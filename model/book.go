package model

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

var (
	// ErrAlreadyBorrowed is returned when attempting to borrow a book that is already borrowed.
	ErrAlreadyBorrowed = errors.New("book is already borrowed")
	// ErrNotBorrowed is returned when attempting to return a book that is not currently borrowed.
	ErrNotBorrowed = errors.New("book is not currently borrowed")
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
	Borrow(id int) (Book, error)
	Return(id int) (Book, error)
	Create(book Book) (Book, error)
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

func (bi BookImplementation) Borrow(id int) (Book, error) {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "acme-books")

	defer client.Close()

	var book Book
	key := datastore.IDKey("Book", int64(id), nil)

	err := client.Get(ctx, key, &book)
	if err != nil {
		return book, err
	}

	if book.Borrowed {
		return book, ErrAlreadyBorrowed
	}

	book.Borrowed = true
	_, err = client.Put(ctx, key, &book)

	return book, err
}

func (bi BookImplementation) Return(id int) (Book, error) {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "acme-books")

	defer client.Close()

	var book Book
	key := datastore.IDKey("Book", int64(id), nil)

	err := client.Get(ctx, key, &book)
	if err != nil {
		return book, err
	}

	if !book.Borrowed {
		return book, ErrNotBorrowed
	}

	book.Borrowed = false
	_, err = client.Put(ctx, key, &book)

	return book, err
}

func (bi BookImplementation) Create(book Book) (Book, error) {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "acme-books")

	defer client.Close()

	key := datastore.IncompleteKey("Book", nil)
	completeKey, err := client.Put(ctx, key, &book)
	if err != nil {
		return book, err
	}

	book.Id = completeKey.ID

	return book, nil
}
