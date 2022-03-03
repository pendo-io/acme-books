package server

import (
	"acme-books/models"
	"context"
	"fmt"

	"cloud.google.com/go/datastore"
)

func Init(host string, port string) {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "acme-books")
	defer client.Close()

	bootstrapBooks(ctx, client)
	r := NewRouter(ctx, client)
	r.RunOnAddr(host + ":" + port)
}

func bootstrapBooks(ctx context.Context, client *datastore.Client) {
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
