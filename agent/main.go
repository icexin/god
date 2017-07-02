package main

import (
	"flag"
	"log"
	"net"

	"github.com/icexin/god/pb"

	"google.golang.org/grpc"
)

var (
	listAddr = flag.String("listen", ":8080", "listen address")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", *listAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterGodAgentServer(grpcServer, &GodAgentService{})
	grpcServer.Serve(lis)
}
