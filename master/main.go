package main

import (
	"flag"
	"log"
	"net"
	"net/http"

	"github.com/icexin/god/pb"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"

	_ "github.com/mattn/go-sqlite3"
)

var (
	listAddr = flag.String("listen", ":9090", "listen address")
	httpAddr = flag.String("http", ":9091", "dashboard address")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", *listAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	db, err := sqlx.Connect("sqlite3", "god.db")
	if err != nil {
		log.Fatal(err)
	}

	gdb = db

	go func() {
		log.Fatal(http.ListenAndServe(*httpAddr, nil))
	}()

	grpcServer := grpc.NewServer()
	pb.RegisterGodMasterServer(grpcServer, &GodMasterService{db: db})
	grpcServer.Serve(lis)
}
