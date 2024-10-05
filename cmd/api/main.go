package main

import (
	"log"

	"payment-api-learning/internal/api"
	"payment-api-learning/internal/storage"
)

func main() {
	jsonStorage, err := storage.NewJSONStorage("./data/customers.json", "./data/logs.json")
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	router := api.SetupRouter(jsonStorage)
	router.Run(":8080")
}
