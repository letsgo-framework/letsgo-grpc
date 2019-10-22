package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/letsgo-framework/letsgo-grpc/services/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	firstName := req.GetGreeting().GetFirstName()
	result := "Greetings " + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}
func main() {
	creds, err := credentials.NewServerTLSFromFile("./keys/server.crt", "./keys/server.key")
	if err != nil {
		log.Fatalf("Failed to generate credentials %v", err)
	}
	grpcServer := grpc.NewServer(grpc.Creds(creds))
	greetpb.RegisterGreetServiceServer(grpcServer, &server{})

	// grpc
	listen, err := net.Listen("tcp", ":50051")

	if err != nil {
		grpclog.Fatalf("failed starting grpc server: %v", err)
	}

	go func() {
		fmt.Println("grpc server running on port 50051")
		if err := grpcServer.Serve(listen); err != nil {
			grpclog.Fatalf("failed starting grpc server: %v", err)
		}
	}()

	// grpc Web
	wrappedServer := grpcweb.WrapServer(grpcServer)
	handler := func(resp http.ResponseWriter, req *http.Request) {
		allowCors(resp, req)
		wrappedServer.ServeHTTP(resp, req)
	}
	httpServer := http.Server{
		Addr:    ":8000",
		Handler: http.HandlerFunc(handler),
	}
	go func() {
		fmt.Println("http server running on port 8000")
		if err := httpServer.ListenAndServe(); err != nil {
			grpclog.Fatalf("failed starting http server: %v", err)
		}
	}()

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}

func allowCors(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Expose-Headers", "grpc-status, grpc-message")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, XMLHttpRequest, x-user-agent, x-grpc-web, grpc-status, grpc-message")
}
