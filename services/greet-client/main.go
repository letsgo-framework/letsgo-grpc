package main

import (
	"context"
	"fmt"
	"log"

	"github.com/letsgo-framework/letsgo-grpc/services/greetpb"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

func main() {

	fmt.Printf("greetpb")
	perRPC := oauth.NewOauthAccess(fetchToken())
	cred, err := credentials.NewClientTLSFromFile("../../keys/server.crt", "localhost")
	if err != nil {
		log.Fatalf("credentials.NewClientTLSFromFile err: %v", err)
	}

	opts := []grpc.DialOption{
		// In addition to the following grpc.DialOption, callers may also use
		// the grpc.CallOption grpc.PerRPCCredentials with the RPC invocation
		// itself.
		// See: https://godoc.org/google.golang.org/grpc#PerRPCCredentials
		grpc.WithPerRPCCredentials(perRPC),
		// oauth.NewOauthAccess requires the configuration of transport
		// credentials.
		grpc.WithTransportCredentials(cred),
	}

	cc, err := grpc.Dial("localhost:50051", opts...)

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

func fetchToken() *oauth2.Token {
	return &oauth2.Token{
		AccessToken: "some-secret-token",
	}
}
