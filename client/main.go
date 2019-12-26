package main

import (
	"context"
	appGrpc "envoy_grpc_sample/grpc"
	"flag"
	"fmt"
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
		host string
		port uint
	)
	flag.StringVar(&name, "name", "", "name")
	flag.StringVar(&host, "host", "", "host")
	flag.UintVar(&port, "port", 0, "port")
	flag.Parse()

	if name == "" || host == "" || port == 0 {
		flag.Usage()
		return
	}

	target := fmt.Sprintf("%s:%d", host, port)
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("%+v", err)
	}
	defer conn.Close()

	client := appGrpc.NewChatClient(conn)

	ctx := context.TODO()
	uid := uuid.New().String()
	sig := make(chan os.Signal)
	wg := &sync.WaitGroup{}

	{
		req := &appGrpc.JoinRequest{
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
	ticker := time.NewTicker(25 * time.Millisecond)
	go func() {
		for {
			select {
			case _, ok := <-ticker.C:
				if !ok {
					return
				}
				req := &appGrpc.SpeakRequest{
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
		req := &appGrpc.LeaveRequest{
			Uuid: uid,
		}
		if _, err := client.Leave(ctx, req); err != nil {
			log.Fatalf("leave err: %+v", err)
		}
	}
	log.Printf("leaved")
}
