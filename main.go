package main

import (
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/datastore"

	"acme-books/databases"
	"acme-books/models"
	"acme-books/server"

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
	client, _ := databases.NewDatabaseClient()

	defer client.Close()

	keys := []*datastore.Key{
		datastore.NameKey("Book", "1984", nil),
		datastore.NameKey("Book", "Animal Farm", nil),
		datastore.NameKey("Book", "Eye of the world", nil),
		datastore.NameKey("Book", "Dictionary", nil),
	}

	books := []interface{}{
		&models.Book{Id: 1, Author: "George Orwell", Title: "1984", Borrowed: false},
		&models.Book{Id: 2, Author: "George Orwell", Title: "Animal Farm", Borrowed: false},
		&models.Book{Id: 3, Author: "Robert Jordan", Title: "Eye of the world", Borrowed: false},
		&models.Book{Id: 4, Author: "Various", Title: "Collins Dictionary", Borrowed: false},
	}

	if _, err := client.PutMulti(keys, books); err != nil {
		fmt.Println(err)
	}
}
