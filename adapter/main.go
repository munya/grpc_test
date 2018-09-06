package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	pb "github.com/munya/grpc_test.git/pb"
	"github.com/munya/grpc_test.git/adapter/adapters"
)

const (
	adapterPort          = ":30301"
	serverURL = "server:30302"
)

func main() {
	log.Println("Adapter is running...")
	lis, err := net.Listen("tcp", adapterPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	// client init
	conn, err := grpc.Dial(serverURL, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewServerClient(conn)
	// client init end

	a, err := adapters.NewAdapter(c)
	inLookupAdapter, err := adapters.NewDictInLookupAdapter(a, adapters.MessageMap{
		"marco": "monkey",
		"polo":  "follow",
	})
	inOutLookupAdapter, err := adapters.NewDictOutLookupAdapter(inLookupAdapter, adapters.MessageMap{
		"monkey": "marco",
		"follow": "polo",
	})

	pb.RegisterAdapterServer(s, inOutLookupAdapter)
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
