package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"pancake/maker/gen/api"
	"pancake/maker/handler"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	if err := realMain(os.Args); err != nil {
		log.Fatalf("error %s", err)
	}
}

func realMain(args []string) error {
	port := 50051
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Errorf("failed to listen: %w", err)
	}

	server := grpc.NewServer()
	api.RegisterPancakeBakerServiceServer(
		server,
		handler.NewBakerHandler(),
	)
	reflection.Register(server)

	go func() {
		log.Printf("start gRPC server port: %v", port)
		server.Serve(lis)
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("stopping gRPC server...")
	server.GracefulStop()

	return nil
}
