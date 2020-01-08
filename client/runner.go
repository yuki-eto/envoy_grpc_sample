package main

import (
	"context"
	"envoy_grpc_sample/pb"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	polyPb "github.com/toydium/polytank/pb"
	"github.com/toydium/polytank/runner"
	"google.golang.org/grpc"
)

type RunnerPlugins struct {
}

func (p *RunnerPlugins) Run(index, timeout uint32, exMap map[string]string) ([]*polyPb.Result, error) {
	name := "test"
	host := exMap["host"]
	port := exMap["port"]

	target := fmt.Sprintf("%s:%s", host, port)
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewChatClient(conn)

	ctx := context.TODO()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	uid := uuid.NewMD5(uuid.New(), []byte(fmt.Sprint(index))).String()

	var res []*polyPb.Result
	log.Printf("start: %d", index)
	r := p.access("Chat.Join", func() error {
		req := &pb.JoinRequest{
			Uuid: uid,
			Name: name,
		}
		_, err := client.Join(ctx, req)
		return err
	})
	res = append(res, r)

	const sendCount = 10
	for i := 0; i < sendCount; i++ {
		r := p.access("Chat.Speak", func() error {
			req := &pb.SpeakRequest{
				Uuid: uid,
				Msg:  time.Now().String(),
			}
			_, err := client.Speak(ctx, req)
			return err
		})
		res = append(res, r)
	}

	{
		r := p.access("Chat.Leave", func() error {
			req := &pb.LeaveRequest{
				Uuid: uid,
			}
			_, err := client.Leave(ctx, req)
			return err
		})
		res = append(res, r)
	}

	log.Printf("end: %d", index)
	return res, nil
}

func (p *RunnerPlugins) access(processName string, f func() error) *polyPb.Result {
	res := &polyPb.Result{
		ProcessName:        processName,
		IsSuccess:          false,
		StartTimestampUsec: time.Now().UnixNano(),
	}
	defer func() {
		res.ElapsedTimeUsec = time.Now().UnixNano() - res.StartTimestampUsec
	}()

	if err := f(); err != nil {
		return res
	}
	res.IsSuccess = true
	return res
}

func main() {
	runner.Serve(&RunnerPlugins{})
}
