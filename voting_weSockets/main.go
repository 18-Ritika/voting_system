package main

import (
	"log"
	"net"
	"voting_weSockets/auth"
	"voting_weSockets/voting"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	auth.RegisterAuthServiceServer(grpcServer)
	voting.RegisterVotingServiceServer(grpcServer)

	log.Println("Starting server on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
