package models

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"google.golang.org/api/iterator"
)

type Book struct {
	Id       int64  `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"writer"`
	Borrowed bool   `json:"borrowed"`
}

func GetBookById(id int64) (Book, error) {
	ctx, client := initDatastore()
	defer client.Close()

	var book Book
	key := datastore.IDKey("Book", id, nil)
	err := client.Get(ctx, key, &book)

	return book, err
}

func GetBooks(filters Filters, order string) ([]Book, error) {
	ctx, client := initDatastore()
	defer client.Close()

	var books []Book

	query := datastore.NewQuery("Book")
	if filters.Title != "" {
		query = query.Filter("Title =", filters.Title)
	}
	if filters.Writer != "" {
		query = query.Filter("Author =", filters.Writer)
	}
	if filters.Available {
		query = query.Filter("Borrowed =", false)
	}
	if order != "" {
		query = query.Order(order)
	}

	it := client.Run(ctx, query)
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

func (book *Book) Borrow() error {
	if book.Borrowed {
		return nil
	}

	book.Borrowed = true
	return book.Save()
}

func (book *Book) Return() error {
	if !book.Borrowed {
		return nil
	}

	book.Borrowed = false
	return book.Save()
}

func (book *Book) Save() error {
	ctx, client := initDatastore()
	defer client.Close()

	_, err := client.Put(ctx, datastore.IDKey("Book", book.Id, nil), book)
	return err
}

func initDatastore() (context.Context, *datastore.Client) {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "acme-books")
	return ctx, client
}
