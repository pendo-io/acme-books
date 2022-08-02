package server

import (
	"acme-books/controllers"
	"acme-books/models"
	"context"
	"log"

	"cloud.google.com/go/datastore"
)

var books = []models.Book{
	{Id: 1, Author: "George Orwell", Title: "1984", Borrowed: false},
	{Id: 2, Author: "George Orwell", Title: "Animal Farm", Borrowed: false},
	{Id: 3, Author: "Robert Jordan", Title: "Eye of the world", Borrowed: false},
	{Id: 4, Author: "Various", Title: "Collins Dictionary", Borrowed: false},
}

func Init(host string, port string) {
	ctx := context.Background()
	if client, err := datastore.NewClient(ctx, "acme-books"); err != nil {
		log.Fatalf("Error creating datastore client: %s", err)
	} else {
		defer client.Close()
		if library, err := controllers.NewLibrary(client, ctx, books); err == nil {
			r := NewRouter(library)
			r.RunOnAddr(host + ":" + port)
		}
	}
}
