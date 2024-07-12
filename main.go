package main

import (
	"acme-books/utils"
	"os"

	"cloud.google.com/go/datastore"

	"acme-books/model"
	"acme-books/server"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	utils.HandleFatalError(err, "Error loading .env file! Did you forget to run `gcloud beta emulators datastore env-init > .env`")

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
	ctx, client := utils.CreateClient()

	defer client.Close()

	books := []model.Book{
		{Id: 1, Author: "George Orwell", Title: "1984", Borrowed: false},
		{Id: 2, Author: "George Orwell", Title: "Animal Farm", Borrowed: false},
		{Id: 3, Author: "Robert Jordan", Title: "Eye of the world", Borrowed: false},
		{Id: 4, Author: "Various", Title: "Collins Dictionary", Borrowed: false},
	}

	var keys []*datastore.Key

	for _, book := range books {
		keys = append(keys, datastore.IDKey("Book", book.Id, nil))
	}

	_, err := client.PutMulti(ctx, keys, books)

	utils.HandleGeneralError(err)
}
