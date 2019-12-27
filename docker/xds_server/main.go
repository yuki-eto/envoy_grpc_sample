package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	api "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	listener "github.com/envoyproxy/go-control-plane/envoy/api/v2/listener"
	route "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	accessLog "github.com/envoyproxy/go-control-plane/envoy/config/accesslog/v2"
	filterAccessLog "github.com/envoyproxy/go-control-plane/envoy/config/filter/accesslog/v2"
	manager "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	filterNetworkRbac "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/rbac/v2"
	rbac "github.com/envoyproxy/go-control-plane/envoy/config/rbac/v2"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	xds "github.com/envoyproxy/go-control-plane/pkg/server"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/duration"
	spb "github.com/golang/protobuf/ptypes/struct"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/uuid"
	"github.com/iancoleman/strcase"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/grpc"
)

const (
	ClusterName = "front-envoy"
	NodeID      = "node1"
)

// define blue/green service
type serviceType string

const (
	BlueService  serviceType = "blue"
	GreenService serviceType = "green"
)

type currentService struct {
	mu *sync.RWMutex
	st serviceType
}

func (s *currentService) Get() serviceType {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.st
}
func (s *currentService) GetString() string {
	return string(s.Get())
}
func (s *currentService) Swap() serviceType {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.st == BlueService {
		s.st = GreenService
	} else {
		s.st = BlueService
	}
	return s.st
}

var cs = &currentService{
	mu: &sync.RWMutex{},
	st: BlueService,
}

type hash struct{}

func (h hash) ID(node *core.Node) string {
	if node == nil {
		return "unknown"
	}
	return node.Cluster + "/" + node.Id
}

// generate configuration
func newSnapshot(current serviceType) cache.Snapshot {
	return cache.NewSnapshot(uuid.New().String(), nil, nil, nil, getListeners(current), nil)
}

// listener configuration
func getListeners(current serviceType) []cache.Resource {
	var listeners []cache.Resource

	ipRestrictFilter := &listener.Filter{
		Name: "envoy.filters.network.rbac",
		ConfigType: &listener.Filter_TypedConfig{
			TypedConfig: rbacConfig(),
		},
	}
	ports := map[serviceType]uint32{
		BlueService:  8001,
		GreenService: 8002,
	}
	for st, port := range ports {
		l := &api.Listener{
			Name: "listener_" + string(st),
			Address: &core.Address{
				Address: &core.Address_SocketAddress{
					SocketAddress: &core.SocketAddress{
						Protocol: core.SocketAddress_TCP,
						Address:  "0.0.0.0",
						PortSpecifier: &core.SocketAddress_PortValue{
							PortValue: port,
						},
					},
				},
			},
			FilterChains: []*listener.FilterChain{{
				Name:    "envoy.http_connection_manager",
				Filters: []*listener.Filter{},
			}},
		}
		httpConnectionManager := &listener.Filter{
			Name: "envoy.http_connection_manager",
			ConfigType: &listener.Filter_TypedConfig{
				TypedConfig: httpConnectionManagerConfig(st),
			},
		}
		filters := &l.FilterChains[0].Filters
		if current != st {
			*filters = append(*filters, ipRestrictFilter)
		}
		*filters = append(*filters, httpConnectionManager)
		listeners = append(listeners, l)
	}

	return listeners
}

// ruled base access control
func rbacConfig() *any.Any {
	conf := &filterNetworkRbac.RBAC{
		StatPrefix: "ip_restrictions.",
		Rules: &rbac.RBAC{
			Action: rbac.RBAC_ALLOW,
			Policies: map[string]*rbac.Policy{
				"allow_host": {
					Permissions: []*rbac.Permission{{
						Rule: &rbac.Permission_Any{Any: true},
					}},
					Principals: []*rbac.Principal{{
						Identifier: &rbac.Principal_SourceIp{SourceIp: &core.CidrRange{
							AddressPrefix: "1.1.1.1",
							PrefixLen:     &wrappers.UInt32Value{Value: 16},
						}},
					}},
				},
			},
		},
	}
	rbacConf, err := ptypes.MarshalAny(conf)
	if err != nil {
		panic(err)
	}
	return rbacConf
}

// http listener config
func httpConnectionManagerConfig(t serviceType) *any.Any {
	accessLogFields := map[string]string{
		"protocol":                          "",
		"duration":                          "",
		"start_time":                        "",
		"bytes_received":                    "",
		"response_code":                     "",
		"bytes_sent":                        "",
		"response_flags":                    "",
		"route_name":                        "",
		"upstream_host":                     "",
		"upstream_transport_failure_reason": "",
		"downstream_remote_address":         "",
		"method":                            "%REQ(:METHOD)%",
		"path":                              "%REQ(X-ENVOY-ORIGINAL-PATH?:PATH)%",
		"x-forwarded-for":                   "%REQ(X-FORWARDED-FOR)%",
		"user-agent":                        "%REQ(USER-AGENT)%",
	}
	jsonFormat := &spb.Struct{
		Fields: map[string]*spb.Value{},
	}
	for k, v := range accessLogFields {
		format := v
		if format == "" {
			format = "%" + strcase.ToScreamingSnake(k) + "%"
		}
		jsonFormat.Fields[k] = &spb.Value{
			Kind: &spb.Value_StringValue{StringValue: format},
		}
	}
	accessLogConf := &accessLog.FileAccessLog{
		Path:            "/dev/stdout",
		AccessLogFormat: &accessLog.FileAccessLog_JsonFormat{JsonFormat: jsonFormat},
	}
	accessLogTypedConf, err := ptypes.MarshalAny(accessLogConf)
	if err != nil {
		panic(err)
	}
	httpConnMgrConfig := &manager.HttpConnectionManager{
		CodecType:  manager.HttpConnectionManager_AUTO,
		StatPrefix: "ingress_http_" + string(t),
		AccessLog: []*filterAccessLog.AccessLog{{
			Name:       "envoy.file_access_log",
			ConfigType: &filterAccessLog.AccessLog_TypedConfig{TypedConfig: accessLogTypedConf},
		}},
		UseRemoteAddress: &wrappers.BoolValue{Value: true},
		HttpFilters: []*manager.HttpFilter{{
			Name: "envoy.router",
		}},
		RouteSpecifier: &manager.HttpConnectionManager_RouteConfig{RouteConfig: &api.RouteConfiguration{
			Name: "route_to_grpc",
			VirtualHosts: []*route.VirtualHost{{
				Name:    "backend",
				Domains: []string{"*"},
				Routes: []*route.Route{{
					Match: &route.RouteMatch{
						PathSpecifier: &route.RouteMatch_Prefix{
							Prefix: "/",
						},
					},
					Action: &route.Route_Route{Route: &route.RouteAction{
						ClusterSpecifier: &route.RouteAction_Cluster{
							Cluster: "grpc_" + string(t),
						},
						Timeout: &duration.Duration{Seconds: 0, Nanos: 0},
					}},
				}},
			}},
		}},
	}
	httpConnMgrTypedConf, err := ptypes.MarshalAny(httpConnMgrConfig)
	if err != nil {
		panic(err)
	}
	return httpConnMgrTypedConf
}

// logger interface for SnapshotCache
type stdLogger struct {
}

func (s *stdLogger) Infof(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (s *stdLogger) Errorf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

// callbacks interface for xds.Server
type callbacks struct {
}

func (c *callbacks) OnStreamOpen(ctx context.Context, id int64, url string) error {
	log.Printf("id: %d, url: %s", id, url)
	return nil
}

func (c *callbacks) OnStreamClosed(int64) {
}

func (c *callbacks) OnStreamRequest(id int64, req *api.DiscoveryRequest) error {
	log.Printf("id: %d, request: %s", id, req.String())
	return nil
}

func (c *callbacks) OnStreamResponse(int64, *api.DiscoveryRequest, *api.DiscoveryResponse) {
}

func (c *callbacks) OnFetchRequest(ctx context.Context, req *api.DiscoveryRequest) error {
	log.Printf("request: %s", req.String())
	return nil
}

func (c *callbacks) OnFetchResponse(*api.DiscoveryRequest, *api.DiscoveryResponse) {
}

func main() {
	var (
		xdsPort   string
		adminPort string
	)
	flag.StringVar(&xdsPort, "port", "30000", "listen xds server port")
	flag.StringVar(&adminPort, "admin-port", "40000", "listen admin api port")
	flag.Parse()

	ctx := context.Background()

	c := cache.NewSnapshotCache(false, hash{}, &stdLogger{})
	nodeName := ClusterName + "/" + NodeID
	if err := c.SetSnapshot(nodeName, newSnapshot(cs.Get())); err != nil {
		panic(err)
	}

	server := xds.NewServer(ctx, c, &callbacks{})
	grpcServer := grpc.NewServer()
	api.RegisterListenerDiscoveryServiceServer(grpcServer, server)

	router := httprouter.New()
	router.GET("/service", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.WriteHeader(200)
		if _, err := w.Write([]byte(cs.GetString())); err != nil {
			log.Printf("err: %+v", err)
		}
	})
	router.POST("/service/swap", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		current := cs.Swap()
		if err := c.SetSnapshot(nodeName, newSnapshot(current)); err != nil {
			w.WriteHeader(500)
			if _, err := w.Write([]byte("cannot set snapshot")); err != nil {
				log.Printf("err: %+v", err)
			}
			return
		}
		w.WriteHeader(200)
		if _, err := w.Write([]byte(current)); err != nil {
			log.Printf("err: %+v", err)
		}
	})

	xdsHost := "0.0.0.0:" + xdsPort
	xdsSock, err := net.Listen("tcp", xdsHost)
	if err != nil {
		panic(err)
	}
	go func() {
		log.Printf("listen xds sever on %s", xdsHost)
		if err := grpcServer.Serve(xdsSock); err != nil {
			log.Printf("closed: %+v", err)
		}
	}()

	adminHost := "0.0.0.0:" + adminPort
	adminSock, err := net.Listen("tcp", adminHost)
	if err != nil {
		panic(err)
	}
	defer adminSock.Close()
	go func() {
		log.Printf("listen admin api server on %s", adminHost)
		if err := http.Serve(adminSock, router); err != nil {
			log.Printf("closed: %+v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGABRT)
	<-sig

	grpcServer.GracefulStop()
}
