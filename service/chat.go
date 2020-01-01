package service

import (
	"context"
	"envoy_grpc_sample/pb"
	"log"
	"os"

	"github.com/cornelk/hashmap"
	"github.com/go-redis/redis/v7"
	"github.com/golang/protobuf/proto"
)

const psChannel = "broadcast"
const mapSize = 100000

type ChatImpl struct {
	redis    redis.UniversalClient
	streamCh chan *pb.ChatStream
	users    *userMap
}

type userMap struct {
	m *hashmap.HashMap
}

func NewUserMap() *userMap {
	return &userMap{
		m: hashmap.New(mapSize),
	}
}
func (m *userMap) Set(u *User) {
	m.m.Set(u.UUID, u)
}
func (m *userMap) Get(uuid string) *User {
	u, ok := m.m.Get(uuid)
	if !ok {
		return nil
	}
	return u.(*User)
}
func (m *userMap) Del(uuid string) bool {
	u := m.Get(uuid)
	if u == nil {
		return false
	}
	m.m.Del(uuid)
	u.deletedCh <- true
	return true
}
func (m *userMap) Broadcast(msg *pb.ChatStream) {
	iter := m.m.Iter()
	for {
		kv, ok := <-iter
		if !ok {
			break
		}
		key, val := kv.Key, kv.Value
		v := val.(*User)
		if err := v.server.Send(msg); err != nil {
			log.Printf("send err: %+v", err)
			m.m.Del(key)
		}
	}
}

func NewRedisCluster() redis.UniversalClient {
	role := os.Getenv("ROLE")
	if role == "prod" {
		return redis.NewClusterClient(&redis.ClusterOptions{
			Addrs: []string{
				"grpc-app-0001-001.02pcoc.0001.apne1.cache.amazonaws.com:6379",
				"grpc-app-0002-001.02pcoc.0001.apne1.cache.amazonaws.com:6379",
				"grpc-app-0003-001.02pcoc.0001.apne1.cache.amazonaws.com:6379",
			},
		})
	}
	return redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{
			"redis-cluster:7000",
			"redis-cluster:7001",
			"redis-cluster:7002",
		},
	})
}

func NewChat() *ChatImpl {
	c := &ChatImpl{
		redis: NewRedisCluster(),
		users: NewUserMap(),
	}
	go c.Subscribe()
	return c
}

func (s *ChatImpl) Subscribe() {
	ps := s.redis.Subscribe(psChannel)
	defer ps.Close()

	if err := ps.Ping(); err != nil {
		log.Printf("[%s] cannot start subscription", hostname)
		return
	}
	if _, err := ps.Receive(); err != nil {
		log.Printf("[%s] cannot start subscription", hostname)
		return
	}

	log.Printf("[%s] start subscription", hostname)

	for msg := range ps.Channel() {
		stream := &pb.ChatStream{}
		if err := proto.Unmarshal([]byte(msg.Payload), stream); err != nil {
			log.Printf("unmarshal err: %+v", err)
			continue
		}

		switch stream.Type {
		case pb.ChatStream_JOINED, pb.ChatStream_LEAVED, pb.ChatStream_SPOKE:
			s.users.Broadcast(stream)
		case pb.ChatStream_LEAVE, pb.ChatStream_SPEAK:
			next := s.next(stream)
			if next != nil {
				if err := s.publish(next); err != nil {
					log.Printf("publish err: %+v", err)
				}
			}
		}
	}

	log.Printf("[%s] end subscription", hostname)
}

func (s *ChatImpl) next(stream *pb.ChatStream) *pb.ChatStream {
	switch stream.Type {
	case pb.ChatStream_LEAVE:
		return s.leave(stream.Uuid)
	case pb.ChatStream_SPEAK:
		return s.speak(stream.Uuid, stream.Msg)
	}
	return nil
}

func (s *ChatImpl) leave(uuid string) *pb.ChatStream {
	u := s.users.Get(uuid)
	if u == nil {
		return nil
	}
	u.deletedCh <- true
	s.users.Del(uuid)
	return &pb.ChatStream{
		Type: pb.ChatStream_LEAVED,
		Uuid: uuid,
		Name: u.Name,
	}
}

func (s *ChatImpl) speak(uuid, msg string) *pb.ChatStream {
	u := s.users.Get(uuid)
	if u == nil {
		return nil
	}
	return &pb.ChatStream{
		Type: pb.ChatStream_SPOKE,
		Uuid: uuid,
		Name: u.Name,
		Msg:  msg,
	}
}

func (s *ChatImpl) publish(stream *pb.ChatStream) error {
	b, err := proto.Marshal(stream)
	if err != nil {
		return err
	}
	if err := s.redis.Publish(psChannel, b).Err(); err != nil {
		return err
	}
	log.Printf("published: %s", stream.String())
	return nil
}

type User struct {
	UUID      string
	Name      string
	server    pb.Chat_JoinServer
	deletedCh chan bool
}

func (s *ChatImpl) Join(req *pb.JoinRequest, server pb.Chat_JoinServer) error {
	uuid := req.Uuid
	name := req.Name
	user := &User{
		UUID:      uuid,
		Name:      name,
		server:    server,
		deletedCh: make(chan bool, 1),
	}
	s.users.Set(user)

	if err := s.publish(&pb.ChatStream{
		Type: pb.ChatStream_JOINED,
		Uuid: user.UUID,
		Name: user.Name,
	}); err != nil {
		return err
	}

	<-user.deletedCh
	return nil
}

func (s *ChatImpl) Leave(ctx context.Context, req *pb.LeaveRequest) (*pb.CommonResponse, error) {
	stream := &pb.ChatStream{
		Type: pb.ChatStream_LEAVE,
		Uuid: req.Uuid,
	}
	if err := s.publish(stream); err != nil {
		return nil, err
	}
	return &pb.CommonResponse{Result: true}, nil
}

func (s *ChatImpl) Speak(ctx context.Context, req *pb.SpeakRequest) (*pb.CommonResponse, error) {
	stream := &pb.ChatStream{
		Type: pb.ChatStream_SPEAK,
		Uuid: req.Uuid,
		Msg:  req.Msg,
	}
	if err := s.publish(stream); err != nil {
		return nil, err
	}
	return &pb.CommonResponse{Result: true}, nil
}

var hostname string

func init() {
	h, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	hostname = h
}
