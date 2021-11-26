package service

import (
	"context"
	"fmt"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"

	"acme-books/models"
)

func client() *datastore.Client {
	c, _ := datastore.NewClient(context.Background(), "acme-books")
	return c
}

func AddOrUpdateStore(book *models.Book) {
	c := client()

	defer c.Close()

	if _, err := c.Put(context.Background(), datastore.NameKey("Book", book.Title, nil), book); err != nil {
		fmt.Println(err)
	}
}

func GetFromStore(query *datastore.Query) []models.Book {
	c := client()

	defer c.Close()

	it := c.Run(context.Background(), query)
	output := make([]models.Book, 0)
	for {
		var b models.Book
		_, err := it.Next(&b)
		if err == iterator.Done {
			break
		}
		output = append(output, b)
	}
	return output
}
