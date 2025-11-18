package main

import (
	"log"
)

// @title           Review Assigner API
// @version         1.0
// @description     API for assigning code reviewers to pull requests

// @host            localhost:8080
// @BasePath        /

// @securityDefinitions.apikey BearerAuth
// @in                         header
// @name                       Authorization
// @description                Type "Bearer" followed by a space and JWT token.

func main() {
	r := InitApp()
	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
