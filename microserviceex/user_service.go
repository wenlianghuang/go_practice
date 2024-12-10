package main

import (
	pb "D:/MattCode/go_practice/microserviceex/matttest/microservice/proto/user"
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
}

func (s *UserService) GetUser(ctx context.Context, in *pb.UserRequest) (*pb.UserResponse, error) {
	return &pb.UserResponse{Message: "Hello " + in.Name}, nil
}

func main() {
	listner, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, &UserService{})
	log.Println("User Service running on port 50051")
	if err := server.Serve(listner); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
