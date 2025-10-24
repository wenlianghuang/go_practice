package routes

import (
	"net/http"

	"complex_breakpoint_example/handlers"
	"complex_breakpoint_example/middleware"
	"complex_breakpoint_example/services"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

// 設置路由
func SetupRoutes(userService *services.UserService, transactionService *services.TransactionService, loanService *services.LoanService) *chi.Mux {
	router := chi.NewRouter()

	// 內建中間件
	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.Recoverer)
	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.RealIP)
	router.Use(chimiddleware.Timeout(60))

	// 自定義中間件
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.CORSMiddleware)
	router.Use(middleware.ValidationMiddleware)
	router.Use(middleware.ErrorHandlingMiddleware)
	router.Use(middleware.RequestSizeMiddleware(1024 * 1024)) // 1MB 限制

	// 創建處理器
	userHandler := handlers.NewUserHandler(userService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	loanHandler := handlers.NewLoanHandler(loanService)

	// API 路由
	router.Route("/api/v1", func(r chi.Router) {
		// 用戶相關路由
		r.Route("/users", func(r chi.Router) {
			r.Post("/", userHandler.CreateUser)
			r.Get("/", userHandler.GetUsers)
			r.Get("/{id}", userHandler.GetUser)
			r.Get("/{id}/account", userHandler.GetUserAccount)
			r.Get("/{id}/transactions", transactionHandler.GetUserTransactions)
			r.Get("/{id}/loans", loanHandler.GetUserLoanApplications)
			r.Post("/{id}/deposit", transactionHandler.Deposit)
			r.Post("/{id}/withdraw", transactionHandler.Withdraw)
			r.Post("/{id}/apply-loan", loanHandler.ApplyForLoan)
		})

		// 交易相關路由
		r.Post("/transfer", transactionHandler.Transfer)

		// 測試路由
		r.Post("/test/concurrent", transactionHandler.ConcurrentTest)
	})

	// 健康檢查
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return router
}
