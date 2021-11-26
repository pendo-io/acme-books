package main

import (
	"log"
	"os"

	"acme-books/models"
	"acme-books/server"
	"acme-books/service"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file! Did you forget to run `gcloud beta emulators datastore env-init > .env`")
	}

	bootstrapBooks()

	host := getEnvWithDefault("HOST", "localhost")
	port := getEnvWithDefault("PORT", "3030")

	server.Init(host, port)
}

func getEnvWithDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func bootstrapBooks() {
	service.AddOrUpdateStore(&models.Book{Id: 1, Author: "George Orwell", Title: "1984", Borrowed: false})
	service.AddOrUpdateStore(&models.Book{Id: 2, Author: "George Orwell", Title: "Animal Farm", Borrowed: false})
	service.AddOrUpdateStore(&models.Book{Id: 3, Author: "Robert Jordan", Title: "Eye of the world", Borrowed: false})
	service.AddOrUpdateStore(&models.Book{Id: 4, Author: "Various", Title: "Collins Dictionary", Borrowed: false})
}
