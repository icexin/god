package main

import (
	"context"
	"testing"

	"github.com/icexin/god/pb"
	"google.golang.org/grpc"
)

func TestGodMasterService(t *testing.T) {
	conn, err := grpc.Dial("127.0.0.1:9090", grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewGodMasterClient(conn)
	req := &pb.SubmitJobRequest{
		Desc: &pb.JobDesc{
			Concurrent: 2,
			Interval:   2,
			Cmd:        "sleep 1",
			Agent:      []string{"127.0.0.1:8080", "127.0.0.1:8080", "127.0.0.1:8080"},
		},
	}
	_, err = client.SubmitJob(context.TODO(), req)
	if err != nil {
		t.Error(err)
	}
}
