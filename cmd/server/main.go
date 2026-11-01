package main

import (
	"context"
	"log"
	"poolworld-pattaya-backend/internal/config"
	"poolworld-pattaya-backend/internal/database"
	"poolworld-pattaya-backend/internal/handler"
	customMiddleware "poolworld-pattaya-backend/internal/middleware"
	"poolworld-pattaya-backend/internal/repository"
	"poolworld-pattaya-backend/internal/usecase"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to MongoDB
	db, err := database.NewMongoDB(cfg.MongoDBURI, cfg.MongoDBDatabase)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.Database)
	customerRepo := repository.NewCustomerRepository(db.Database)
	poolCareRepo := repository.NewPoolCareRepository(db.Database)
	poolCareReportRepo := repository.NewPoolCareReportRepository(db.Database)

	if err := poolCareReportRepo.EnsureIndexes(context.Background()); err != nil {
		log.Println("Warning: Failed to create pool care report indexes (duplicate reports may exist):", err)
	}

	// Initialize usecases
	authUsecase := usecase.NewAuthUsecase(userRepo, cfg.JWTSecret)
	customerUsecase := usecase.NewCustomerUsecase(customerRepo, poolCareRepo, poolCareReportRepo)
	poolCareUsecase := usecase.NewPoolCareUsecase(poolCareRepo, customerRepo, poolCareReportRepo)
	poolCareReportUsecase := usecase.NewPoolCareReportUsecase(poolCareRepo, poolCareReportRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authUsecase)
	customerHandler := handler.NewCustomerHandler(customerUsecase)
	poolCareHandler := handler.NewPoolCareHandler(poolCareUsecase, customerUsecase)
	poolCareReportHandler := handler.NewPoolCareReportHandler(poolCareReportUsecase)

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(customMiddleware.CORS())

	// Routes
	api := e.Group("/api")

	// Auth routes
	api.POST("/login", authHandler.Login)

	// Protected routes
	protected := api.Group("")
	protected.Use(customMiddleware.JWTMiddleware(cfg.JWTSecret))

	// Customer routes
	protected.GET("/customers", customerHandler.GetAllCustomers)
	protected.GET("/customers/:id", customerHandler.GetCustomerByID)
	protected.POST("/customers", customerHandler.CreateCustomer)
	protected.PUT("/customers/:id", customerHandler.UpdateCustomer)

	// Pool care routes

	// Pool care report routes
	poolCareReportHandler.RegisterRoutes(protected)

	poolCareHandler.RegisterRoutes(protected)
	protected.DELETE("/customers/:id", customerHandler.DeleteCustomer)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := e.Start(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
