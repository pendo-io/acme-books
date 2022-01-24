package main

import (
	"log"
	"os"

	"acme-books/models"
	"acme-books/server"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file! Did you forget to run `gcloud beta emulators datastore env-init > .env`")
	}

	models.Setup()
	models.BootstrapBooks()

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
