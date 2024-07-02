package main

import (
	"github.com/getz-devs/librakeeper-server/internal/server/handlers"
	"github.com/getz-devs/librakeeper-server/internal/server/services/auth"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
)

// WARN: old code
// TODO: move some code from cmd here

func main() {
	r := gin.Default()

	// Configure CORS
	config := cors.Config{
		AllowOrigins:     []string{"https://libra.potat.dev", "http://localhost:3000"}, // Allow specific origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},          // Allow methods
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},          // Allow headers including Authorization
		AllowCredentials: true,                                                         // Allow credentials
	}

	r.Use(cors.New(config))

	// Initialize Firebase
	err := auth.InitializeFirebase("internal/server/.secret.json") // TODO: read config
	if err != nil {
		log.Fatalf("error initializing Firebase: %v", err)
	}

	// Routes
	r.POST("/health", handlers.HealthCheck)

	err = r.Run(":8080") // TODO: read config
	if err != nil {
		log.Fatalf("error starting server: %v", err)
		return
	}
}
