package models

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

func BootstrapBooks() {
	ctx := context.Background()

	books := []Book{
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

func FindBookById(id int) (Book, error) {
	ctx := context.Background()

	var book Book
	key := datastore.IDKey("Book", int64(id), nil)

	if err := client.Get(ctx, key, &book); err != nil {
		fmt.Println(err)
		return Book{}, err
	}

	return book, nil
}

func ListBooks() ([]Book, error) {
	ctx := context.Background()

	var books []Book

	it := client.Run(ctx, datastore.NewQuery("Book"))
	for {
		var b Book
		_, err := it.Next(&b)
		if err == iterator.Done {
			fmt.Println(err)
			break
		}
		books = append(books, b)
	}

	return books, nil
}
