package main

import (
	"context"
	"flag"
	"log"
	"net"

	api "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	listener "github.com/envoyproxy/go-control-plane/envoy/api/v2/listener"
	route "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	accessLog "github.com/envoyproxy/go-control-plane/envoy/config/accesslog/v2"
	filterAccessLog "github.com/envoyproxy/go-control-plane/envoy/config/filter/accesslog/v2"
	manager "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	envoy_config_filter_network_rbac_v2 "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/rbac/v2"
	envoy_config_rbac_v2 "github.com/envoyproxy/go-control-plane/envoy/config/rbac/v2"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	xds "github.com/envoyproxy/go-control-plane/pkg/server"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc"
)

const (
	ClusterName = "front-envoy"
	NodeID      = "node1"
)

type hash struct {
}

func (h hash) ID(node *core.Node) string {
	if node == nil {
		return "unknown"
	}
	return node.Cluster + "/" + node.Id
}

func httpConnectionManagerConfig(lType string) *any.Any {
	accessLogConf := &accessLog.FileAccessLog{
		Path: "/dev/stdout",
	}
	accessLogTypedConf, err := ptypes.MarshalAny(accessLogConf)
	if err != nil {
		panic(err)
	}
	httpConnMgrConfig := &manager.HttpConnectionManager{
		CodecType:  manager.HttpConnectionManager_AUTO,
		StatPrefix: "ingress_http_" + lType,
		AccessLog: []*filterAccessLog.AccessLog{{
			Name:       "envoy.file_access_log",
			ConfigType: &filterAccessLog.AccessLog_TypedConfig{TypedConfig: accessLogTypedConf},
		}},
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
							Cluster: "grpc_" + lType,
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

func rbacConfig() *any.Any {
	rbac := &envoy_config_filter_network_rbac_v2.RBAC{
		StatPrefix: "ip_restrictions.",
		Rules: &envoy_config_rbac_v2.RBAC{
			Action: envoy_config_rbac_v2.RBAC_ALLOW,
			Policies: map[string]*envoy_config_rbac_v2.Policy{
				"allow_host": {
					Permissions: []*envoy_config_rbac_v2.Permission{{
						Rule: &envoy_config_rbac_v2.Permission_Any{Any: true},
					}},
					Principals: []*envoy_config_rbac_v2.Principal{{
						Identifier: &envoy_config_rbac_v2.Principal_SourceIp{SourceIp: &core.CidrRange{
							AddressPrefix: "192.168.0.1",
							PrefixLen:     &wrappers.UInt32Value{Value: 16},
						}},
					}},
				},
			},
		},
	}
	rbacConf, err := ptypes.MarshalAny(rbac)
	if err != nil {
		panic(err)
	}
	return rbacConf
}

func defaultSnapshot() cache.Snapshot {
	var listeners []cache.Resource

	ipRestrictFilter := &listener.Filter{
		Name: "envoy.filters.network.rbac",
		ConfigType: &listener.Filter_TypedConfig{
			TypedConfig: rbacConfig(),
		},
	}
	_ = ipRestrictFilter
	ports := map[string]uint32{
		"blue":  8001,
		"green": 8002,
	}
	for lType, port := range ports {
		l := &api.Listener{
			Name: "listener_" + lType,
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
				Name: "envoy.http_connection_manager",
				Filters: []*listener.Filter{
					ipRestrictFilter,
					{
						Name: "envoy.http_connection_manager",
						ConfigType: &listener.Filter_TypedConfig{
							TypedConfig: httpConnectionManagerConfig(lType),
						},
					},
				},
			}},
		}
		listeners = append(listeners, l)
	}

	return cache.Snapshot{
		Endpoints: cache.Resources{},
		Clusters:  cache.Resources{},
		Routes:    cache.Resources{},
		Listeners: cache.NewResources("1.0", listeners),
		Secrets:   cache.Resources{},
		Runtimes:  cache.Resources{},
	}
}

func main() {
	var port string
	flag.StringVar(&port, "port", "30000", "listen port")
	flag.Parse()

	c := cache.NewSnapshotCache(false, hash{}, nil)
	if err := c.SetSnapshot(ClusterName+"/"+NodeID, defaultSnapshot()); err != nil {
		panic(err)
	}
	server := xds.NewServer(context.Background(), c, nil)
	grpcServer := grpc.NewServer()
	api.RegisterListenerDiscoveryServiceServer(grpcServer, server)

	host := "0.0.0.0:" + port
	sock, err := net.Listen("tcp", host)
	if err != nil {
		panic(err)
	}

	log.Printf("listen eds sever on %s", host)
	if err := grpcServer.Serve(sock); err != nil {
		panic(err)
	}
}
