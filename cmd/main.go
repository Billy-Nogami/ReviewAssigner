package main

import (
	"log"
)
func main() {
	r := InitApp()
	log.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
