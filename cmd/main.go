// cmd/main.go
package main

import (
	"log"
	"os"

	"wallet-api/handlers"
	"wallet-api/models"
	"wallet-api/repositories"
	"wallet-api/services"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Database connection
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=wallet_db port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&models.User{}, &models.Wallet{}, &models.Transaction{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Repositories
	userRepo := repositories.NewUserRepository(db)
	walletRepo := repositories.NewWalletRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)

	// Services
	userService := services.NewUserService(userRepo)
	walletService := services.NewWalletService(walletRepo, userRepo)
	transferService := services.NewTransferService(transactionRepo, walletRepo, db)

	// Handlers
	userHandler := handlers.NewUserHandler(userService)
	walletHandler := handlers.NewWalletHandler(walletService)
	transferHandler := handlers.NewTransferHandler(transferService)

	// Router
	router := gin.Default()

	// Routes
	v1 := router.Group("/api/v1")
	{
		// User routes
		v1.POST("/users", userHandler.Create)
		v1.GET("/users/:id", userHandler.GetByID)

		// Wallet routes
		v1.POST("/wallets", walletHandler.Create)
		v1.GET("/wallets/:id", walletHandler.GetByID)
		v1.GET("/users/:id/wallets", walletHandler.GetByUserID)

		// Transfer routes
		v1.POST("/transfers", transferHandler.Transfer)
		v1.POST("/deposits", transferHandler.Deposit)
		v1.GET("/wallets/:id/transactions", transferHandler.GetTransactions)
	}

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
