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

	// 使用 PostgreSQL + GORM
	fmt.Println("📊 Using PostgreSQL database with GORM")

	gormDB, err := database.NewGormDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v\n\nHint: Make sure PostgreSQL is running:\n  brew services start postgresql", err)
	}

	if err := gormDB.Initialize(); err != nil {
		log.Printf("Warning: Failed to initialize database: %v", err)
	}

	fmt.Println("✅ Connected to PostgreSQL database")

	// 創建 GORM 適配器
	dbAdapter := database.NewGormAdapter(gormDB)

	// 初始化服務
	userService := services.NewUserService(dbAdapter)
	transactionService := services.NewTransactionService(dbAdapter)
	loanService := services.NewLoanService(dbAdapter)

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
