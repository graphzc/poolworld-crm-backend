package main

import (
	"context"
	"log"
	"poolworld-pattaya-backend/internal/config"
	"poolworld-pattaya-backend/internal/database"
	"poolworld-pattaya-backend/internal/repository"
	"poolworld-pattaya-backend/internal/usecase"
)

func main() {
	// Load configuration
	cfg := config.Load()
	if cfg.AdminUsername == "" || cfg.AdminPassword == "" {
		log.Fatal("ADMIN_USERNAME and ADMIN_PASSWORD must be set")
	}

	// Connect to MongoDB
	db, err := database.NewMongoDB(cfg.MongoDBURI, cfg.MongoDBDatabase)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.Database)

	// Initialize usecases
	authUsecase := usecase.NewAuthUsecase(userRepo, cfg.JWTSecret)

	// Create admin user
	ctx := context.Background()
	err = authUsecase.CreateUser(ctx, cfg.AdminUsername, cfg.AdminPassword)
	if err != nil {
		log.Fatal("Failed to create admin user:", err)
	}

	log.Println("Admin user created successfully!")
}
