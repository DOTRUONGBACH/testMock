package main

import (
	"log"
	"mock-project/database"
	"mock-project/services/customer-service/middlewares"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	database.NewClient()
	lis, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer(
		grpc.UnaryInterceptor(middlewares.JWTUnaryInterceptor),
		grpc.StreamInterceptor(middlewares.JWTStreamInterceptor),
	)
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
