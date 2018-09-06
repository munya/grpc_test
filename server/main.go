package main

import (
	"log"
	"net"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"github.com/munya/grpc_test.git/pb"
)

const (
	port = ":30302"
)

type Server struct {
	Dict map[string]string
}

func (s *Server) Send(ctx context.Context, phrase *pbservice.Params) (*pbservice.Params, error) {
	return &pbservice.Params{Message: s.Dict[phrase.Message]}, nil
}

func main() {
	log.Println("Server is running...")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	m := &Server{
		Dict: map[string]string{
			"monkey": "follow",
			"follow": "monkey",
		},
	}
	pbservice.RegisterServerServer(s, m)
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
