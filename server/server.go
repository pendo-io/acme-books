package server

import (
	"acme-books/models"
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
)

func Init(host string, port string) {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "acme-books")
	defer client.Close()

	bootstrapBooks(client, ctx)

	r := NewRouter(client, ctx)
	r.RunOnAddr(host + ":" + port)
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