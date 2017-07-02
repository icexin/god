package main

import (
	"flag"
	"log"
	"net"

	"github.com/icexin/god/pb"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"

	_ "github.com/cznic/ql/driver"
)

var (
	listAddr = flag.String("listen", ":9090", "listen address")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", *listAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	db, err := sqlx.Connect("ql", "god.db")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterGodMasterServer(grpcServer, &GodMasterService{db: db})
	grpcServer.Serve(lis)
}
