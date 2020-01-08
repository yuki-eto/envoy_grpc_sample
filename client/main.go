package main

import (
	"context"
	"envoy_grpc_sample/pb"
	"flag"
	"io"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

func main() {
	var (
		name string
		addr string
	)
	flag.StringVar(&name, "name", "", "name")
	flag.StringVar(&addr, "addr", "", "addr")
	flag.Parse()

	if name == "" || addr == "" {
		flag.Usage()
		return
	}

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("%+v", err)
	}
	defer conn.Close()

	client := pb.NewChatClient(conn)

	ctx := context.TODO()
	uid := uuid.New().String()
	sig := make(chan os.Signal)
	wg := &sync.WaitGroup{}

	{
		req := &pb.JoinRequest{
			Uuid: uid,
			Name: name,
		}
		res, err := client.Join(ctx, req)
		if err != nil {
			log.Fatalf("%+v", err)
		}

		wg.Add(1)
		go func() {
			wg.Done()
			for {
				recv, err := res.Recv()
				if err == io.EOF {
					log.Print("receive: EOF")
					return
				}
				if err != nil {
					log.Printf("recv err: %+v", err)
					sig <- syscall.SIGABRT
					return
				}
				log.Printf("receive: %s", recv.String())
			}
		}()
	}

	log.Print("wait for speak")
	wg.Wait()
	ticker := time.NewTicker(100 * time.Millisecond)
	go func() {
		for {
			select {
			case _, ok := <-ticker.C:
				if !ok {
					return
				}
				req := &pb.SpeakRequest{
					Uuid: uid,
					Msg:  time.Now().String(),
				}
				if _, err := client.Speak(ctx, req); err != nil {
					log.Fatalf("speak err: %+v", err)
				}
				log.Printf("speaked")
			}
		}
	}()

	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGABRT)

	<-sig
	{
		req := &pb.LeaveRequest{
			Uuid: uid,
		}
		if _, err := client.Leave(ctx, req); err != nil {
			log.Fatalf("leave err: %+v", err)
		}
	}
	log.Printf("leaved")
}
