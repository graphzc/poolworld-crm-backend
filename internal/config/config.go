package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoDBURI      string
	MongoDBDatabase string
	JWTSecret       string
	Port            string
	AdminUsername   string
	AdminPassword   string
}

func Load() *Config {
	// Try to load .env file, but don't fail if it doesn't exist
	// This allows for both .env file and runtime environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using runtime environment variables")
	}

	return &Config{
		MongoDBURI:      getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		MongoDBDatabase: getEnv("MONGODB_DATABASE", "poolworld_pattaya"),
		JWTSecret:       getEnv("JWT_SECRET", "poolworld-secret-key"),
		Port:            getEnv("PORT", "8080"),
		AdminUsername:   getEnv("ADMIN_USERNAME", ""),
		AdminPassword:   getEnv("ADMIN_PASSWORD", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
