package main

import (
	pb "GRPCADDER/pkg/api/proto"
	"GRPCADDER/pkg/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	s := grpc.NewServer()
	srv := &service.GRPCServer{}
	pb.RegisterCalculatorServer(s, srv)

	l, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}
	if err := s.Serve(l); err != nil {
		log.Fatal(err)
	}
}
