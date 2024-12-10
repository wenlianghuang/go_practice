package main

import (
	"context"
	"log"
	"net"

	pb "microserviceex/proto/order"

	"google.golang.org/grpc"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
}

func (s *OrderService) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	return &pb.GetOrderResponse{
		Id:     req.Id,
		UserId: req.UserId,
		Amount: 99.99,
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterOrderServiceServer(server, &OrderService{})

	log.Println("Order Service running on port 50052")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
