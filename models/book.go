package models

import (
	"acme-books/helpers"
	"context"
	"fmt"
	"net/url"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

type Book struct {
	Id       int64
	Title    string `json:"title"`
	Author   string `json:"writer"`
	Borrowed bool   `json:"borrowed"`
}

func GetAllBooks(filters url.Values, client *datastore.Client, ctx context.Context) []Book {
	var output []Book
	it := client.Run(ctx, createQuery(filters))
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

func GetSingleBook(client *datastore.Client, ctx context.Context, id int) (Book, error) {
	var book Book
	key := datastore.IDKey("Book", int64(id), nil)

	err := client.Get(ctx, key, &book)
	return book, err
}

func createQuery(filters url.Values) *datastore.Query {
	if len(filters) > 0 {
		if author := filters.Get("author"); author != "" {
			return datastore.NewQuery("Book").Filter("Author=", author).Order("Id")
		}
		if title := filters.Get("title"); title != "" {
			return datastore.NewQuery("Book").Filter("Title=", title).Order("Id")
		}
		if borrowed := filters.Get("borrowed"); borrowed != "" {
			return datastore.NewQuery("Book").Filter("Borrowed=", (borrowed == "true")).Order("Id")
		}
	}
	return datastore.NewQuery("Book").Order("Id")
}

func CreateBook(client *datastore.Client, ctx context.Context, book Book) (Book, error) {
	key := datastore.IncompleteKey("Book", nil)

	newId, err := client.Put(ctx, key, &book)
	book.Id = newId.ID
	return book, err
}

func BorrowBook(client *datastore.Client, ctx context.Context, id int) (Book, error) {
	var book Book
	key := datastore.IDKey("Book", int64(id), nil)

	err := client.Get(ctx, key, &book)
	if err != nil {
		return book, err
	}
	if book.Borrowed == true {
		return book, &helpers.BoorrowError{"The requested book has already been borrowed"}
	}
	book.Borrowed = true
	_, err = client.Put(ctx, key, &book)
	return book, err
}

func ReturnBook(client *datastore.Client, ctx context.Context, id int) (Book, error) {
	var book Book
	key := datastore.IDKey("Book", int64(id), nil)

	err := client.Get(ctx, key, &book)
	if err != nil {
		return book, err
	}
	if book.Borrowed == false {
		return book, &helpers.ReturnBookError{"The requested book has never been borrowed"}
	}
	book.Borrowed = false
	_, err = client.Put(ctx, key, &book)
	return book, err
}
