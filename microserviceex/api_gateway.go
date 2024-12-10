package main

import (
	"context"
	"log"
	"net/http"

	orderpb "microserviceex/proto/order"
	userpb "microserviceex/proto/user"

	"google.golang.org/grpc"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 連接 User Service
	userConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to User Service: %v", err)
	}
	defer userConn.Close()
	userClient := userpb.NewUserServiceClient(userConn)

	// 連接 Order Service
	orderConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to Order Service: %v", err)
	}
	defer orderConn.Close()
	orderClient := orderpb.NewOrderServiceClient(orderConn)

	// 定義 API 路由
	r.GET("/user/:id", func(c *gin.Context) {
		id := c.Param("id")
		resp, err := userClient.GetUser(context.Background(), &userpb.GetUserRequest{Id: id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": resp.Id, "name": resp.Name, "age": resp.Age})
	})

	r.GET("/order/:id", func(c *gin.Context) {
		id := c.Param("id")
		userId := c.Query("userId")
		resp, err := orderClient.GetOrder(context.Background(), &orderpb.GetOrderRequest{Id: id, UserId: userId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": resp.Id, "userId": resp.UserId, "amount": resp.Amount})
	})

	// 啟動 API Gateway
	r.Run(":8080")
}
