package main

import (
	"fmt"
	"log"
	"net/http"

	"complex_breakpoint_example/config"
	"complex_breakpoint_example/database"
	"complex_breakpoint_example/routes"
	"complex_breakpoint_example/services"
)

func main() {
	// 初始化配置
	cfg := config.Load()

	// 初始化數據庫
	db := database.New()
	db.Initialize()

	// 初始化服務
	userService := services.NewUserService(db)
	transactionService := services.NewTransactionService(db)
	loanService := services.NewLoanService(db)

	// 設置路由
	router := routes.SetupRoutes(userService, transactionService, loanService)

	fmt.Printf("🚀 Complex Breakpoint Example Server starting on :%s\n", cfg.Port)
	fmt.Println("📚 Available endpoints:")
	fmt.Println("  POST /api/v1/users - Create user")
	fmt.Println("  GET  /api/v1/users - Get all users")
	fmt.Println("  GET  /api/v1/users/{id} - Get user")
	fmt.Println("  GET  /api/v1/users/{id}/account - Get user account")
	fmt.Println("  GET  /api/v1/users/{id}/transactions - Get user transactions")
	fmt.Println("  GET  /api/v1/users/{id}/loans - Get user loan applications")
	fmt.Println("  POST /api/v1/users/{id}/deposit - Deposit money")
	fmt.Println("  POST /api/v1/users/{id}/withdraw - Withdraw money")
	fmt.Println("  POST /api/v1/transfer - Transfer money")
	fmt.Println("  POST /api/v1/users/{id}/apply-loan - Apply for loan")
	fmt.Println("  POST /api/v1/test/concurrent - Concurrent operations test")
	fmt.Println("  GET  /health - Health check")

	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
