package main

import (
	"envoy_grpc_sample/pb"
	"envoy_grpc_sample/service"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

func main() {
	port := os.Getenv("PORT_NUMBER")
	if port == "" {
		port = "29999"
	}
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("%+v", err)
	}

	server := grpc.NewServer()
	s := service.NewChat()
	pb.RegisterChatServer(server, s)

	log.Printf("running server on port: %s", port)
	go func() {
		if err := server.Serve(l); err != nil {
			log.Printf("end Serve: %+v", err)
		}
	}()

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGABRT)

	log.Print("wait for interrupt...")
	<-sig

	server.GracefulStop()
}
