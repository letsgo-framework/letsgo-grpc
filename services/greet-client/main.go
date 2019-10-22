package main

import (
	"context"
	"fmt"
	"log"

	"github.com/letsgo-framework/letsgo-grpc/services/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {

	fmt.Printf("greetpb")
	cred, err := credentials.NewClientTLSFromFile("../../keys/server.crt", "localhost")
	if err != nil {
		log.Fatalf("credentials.NewClientTLSFromFile err: %v", err)
	}

	cc, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(cred))

	if err != nil {
		log.Fatalf("could not connect %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)

	g := &greetpb.Greeting{
		FirstName: "ls",
	}

	req := &greetpb.GreetRequest{
		Greeting: g,
	}

	res, err := c.Greet(context.Background(), req)

	if err != nil {
		log.Fatalf("error while calling watcher %v", err)
	}
	log.Println(res)
}
