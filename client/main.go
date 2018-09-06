package main

import (
	pb "github.com/munya/grpc_test.git/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"os"
	"time"
)

const (
	adapterURL= "adapter:30301"
)

func main() {
	log.Println("Client is running...")
	if len(os.Args) == 1 {
		log.Fatalln("Parameter in command line is missing.")
	}
	message := os.Args[1]

	conn, err := grpc.Dial(adapterURL, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Can't connect to adapter: %v", err)
	}
	defer conn.Close()
	adapterClient := pb.NewAdapterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sendMessage := &pb.Params{Message: message}
	receivedMessage, err := adapterClient.Send(ctx, sendMessage)
	if err != nil {
		log.Fatalf("%v", err)
	}

	log.Printf("%s", sendMessage)
	log.Printf("%s", receivedMessage)
}