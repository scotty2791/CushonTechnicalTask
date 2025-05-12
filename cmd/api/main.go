package main

import (
	"log"
	"os"

	"cushon/internal/adapters/primary/http"
	"cushon/internal/adapters/secondary/persistence/mysql"
	"cushon/internal/core/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize MySQL connection
	dbConfig := mysql.Config{
		Host:     "localhost",
		Port:     3306,
		User:     "root",
		Password: os.Getenv("DB_PASSWORD"),
		Database: "cushon",
	}

	db, err := mysql.NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	directUserRepo := mysql.NewDirectUserRepository(db)
	transactionRepo := mysql.NewTransactionRepository(db)

	// Initialize services
	directUserService := services.NewDirectUserService(directUserRepo)
	transactionService := services.NewTransactionService(transactionRepo)

	// Initialize handlers
	directUserHandler := http.NewDirectUserHandler(directUserService)
	transactionHandler := http.NewTransactionHandler(transactionService)

	// Initialize router
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12 hours
	}))

	// Register routes
	directUserHandler.RegisterRoutes(router)
	transactionHandler.RegisterRoutes(router)

	// Add health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
} 