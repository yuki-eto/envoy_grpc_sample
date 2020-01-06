package main

import (
	"context"
	"envoy_grpc_sample/pb"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

func Run(index, timeout uint32, ch chan []string, exMap map[string]string) error {
	name := "test"
	host := exMap["host"]
	port := exMap["port"]

	target := fmt.Sprintf("%s:%s", host, port)
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewChatClient(conn)

	ctx := context.TODO()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	uid := uuid.NewMD5(uuid.New(), []byte(fmt.Sprint(index))).String()

	req := &pb.JoinRequest{
		Uuid: uid,
		Name: name,
	}
	log.Printf("start: %d", index)
	start := time.Now()
	_, err = client.Join(ctx, req)
	isSuccess := true
	if err != nil {
		isSuccess = false
	}
	end := time.Now()
	ch <- []string{
		"Chat.Join",
		fmt.Sprint(isSuccess),
		fmt.Sprint(start.UnixNano()),
		fmt.Sprint(end.UnixNano() - start.UnixNano()),
	}

	const sendCount = 10
	for i := 0; i < sendCount; i++ {
		req := &pb.SpeakRequest{
			Uuid: uid,
			Msg:  time.Now().String(),
		}
		start := time.Now()
		isSuccess := true
		if _, err := client.Speak(ctx, req); err != nil {
			isSuccess = false
		}
		end := time.Now()
		ch <- []string{
			"Chat.Speak",
			fmt.Sprint(isSuccess),
			fmt.Sprint(start.UnixNano()),
			fmt.Sprint(end.UnixNano() - start.UnixNano()),
		}
	}

	{
		req := &pb.LeaveRequest{
			Uuid: uid,
		}
		start := time.Now()
		isSuccess := true
		if _, err := client.Leave(ctx, req); err != nil {
			isSuccess = false
		}
		end := time.Now()
		ch <- []string{
			"Chat.Leave",
			fmt.Sprint(isSuccess),
			fmt.Sprint(start.UnixNano()),
			fmt.Sprint(end.UnixNano() - start.UnixNano()),
		}
	}

	log.Printf("end: %d", index)
	return nil
}
