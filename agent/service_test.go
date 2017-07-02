package main

import (
	"context"
	"io"
	"testing"

	"github.com/icexin/god/pb"

	"google.golang.org/grpc"
)

func TestGodAgent(t *testing.T) {
	conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewGodAgentClient(conn)
	req := &pb.RunJobRequest{
		Id: 1,
		Desc: &pb.JobDesc{
			Cmd: "ls sdf",
		},
	}
	stream, err := client.RunJob(context.TODO(), req)
	if err != nil {
		t.Error(err)
	}
	var resp *pb.RunJobResponse
	for {
		resp, err = stream.Recv()
		if err == io.EOF {
			break
		}
		t.Log(resp.Body)
		if resp.Stoped {
			t.Log("code", resp.Code)
		}
	}
}
