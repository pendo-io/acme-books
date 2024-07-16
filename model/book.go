package model

import (
	"context"
	"fmt"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

type Book struct {
	Id       int64
	Title    string `json:"title"`
	Author   string `json:"writer"`
	Borrowed bool   `json:"borrowed"`
}

type BookInterface interface {
	GetByKey(id int) (Book, error)
	ListAll() []Book
}

type BookImplementation struct {
}

func (bi BookImplementation) GetByKey(id int) (Book, error) {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "acme-books")

	defer client.Close()

	var book Book
	key := datastore.IDKey("Book", int64(id), nil)

	err := client.Get(ctx, key, &book)

	return book, err
}

func (bi BookImplementation) ListAll() []Book {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "acme-books")

	defer client.Close()

	var output []Book

	it := client.Run(ctx, datastore.NewQuery("Book"))
	for {
		var b Book
		_, err := it.Next(&b)
		if err == iterator.Done {
			fmt.Println(err)
			break
		}
		output = append(output, b)
	}

	return output
}
