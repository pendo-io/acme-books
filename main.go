package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/datastore"

	"acme-books/models"
	"acme-books/server"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file! Did you forget to run `gcloud beta emulators datastore env-init > .env`")
	}

	host := getEnvWithDefault("HOST", "localhost")
	port := getEnvWithDefault("PORT", "3030")
	client, ctx := initialise()
	bootstrapBooks(client, ctx)
	defer client.Close()

	server.Init(host, port, client, ctx)
}

func getEnvWithDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func bootstrapBooks(client *datastore.Client, ctx context.Context) {
	books := []models.Book{
		{Id: 1, Author: "George Orwell", Title: "1984", Borrowed: false},
		{Id: 2, Author: "George Orwell", Title: "Animal Farm", Borrowed: false},
		{Id: 3, Author: "Robert Jordan", Title: "Eye of the world", Borrowed: false},
		{Id: 4, Author: "Various", Title: "Collins Dictionary", Borrowed: false},
	}

	var keys []*datastore.Key

	for _, book := range books {
		keys = append(keys, datastore.IDKey("Book", book.Id, nil))
	}

	if _, err := client.PutMulti(ctx, keys, books); err != nil {
		fmt.Println(err)
	}
}

func initialise() (*datastore.Client, context.Context) {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "acme-books")
	return client, ctx
}
