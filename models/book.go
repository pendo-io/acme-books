package models

import (
	"context"
	"fmt"

	"cloud.google.com/go/datastore"
)

type Book struct {
	Id       int64
	Title    string `json:"title"`
	Author   string `json:"author"`
	Borrowed bool   `json:"borrowed"`
}

type BookFilter struct {
	Author string
	Title  string
}

func BootstrapBooks() {
	ctx := context.Background()
	books, _ := ListBooks(BookFilter{})

	for _, book := range books {
		client.Delete(ctx, datastore.IDKey("Book", book.Id, nil))
	}

	AddBook(Book{Author: "George Orwell", Title: "1984"})
	AddBook(Book{Author: "George Orwell", Title: "Animal Farm", Borrowed: false})
	AddBook(Book{Author: "Robert Jordan", Title: "Eye of the world", Borrowed: false})
	AddBook(Book{Author: "Various", Title: "Collins Dictionary", Borrowed: false})
}

func FindBookById(id int) (Book, error) {
	ctx := context.Background()

	var book Book
	key := datastore.IDKey("Book", int64(id), nil)

	if err := client.Get(ctx, key, &book); err != nil {
		fmt.Println(err)

		if book.Id == 0 {
			return Book{}, datastore.ErrNoSuchEntity
		}
		return Book{}, err
	}
	book.Id = key.ID

	return book, nil
}

func AddBook(book Book) (Book, error) {
	ctx := context.Background()

	key := datastore.IncompleteKey("Book", nil)

	key, err := client.Put(ctx, key, &book)

	if err != nil {
		fmt.Println(err)
		return book, err
	}

	book.Id = key.ID

	return book, nil
}

func ListBooks(filters BookFilter) ([]*Book, error) {
	ctx := context.Background()

	var books []*Book

	query := buildBookQuery(filters)

	keys, err := client.GetAll(ctx, query, &books)

	if err != nil {
		return nil, err
	}

	for i, key := range keys {
		books[i].Id = key.ID
	}

	return books, nil
}

func BorrowBook(book Book) (Book, error) {
	book.Borrowed = true

	book, _ = updateBook(book)

	return book, nil
}

func ReturnBook(book Book) (Book, error) {
	book.Borrowed = false

	book, _ = updateBook(book)

	return book, nil
}

func updateBook(book Book) (Book, error) {
	ctx := context.Background()

	key := datastore.IDKey("Book", book.Id, nil)
	_, err := client.Put(ctx, key, &book)

	if err != nil {
		fmt.Println(err)
		return book, err
	}

	return book, err
}

func buildBookQuery(filters BookFilter) *datastore.Query {
	query := datastore.NewQuery("Book").Order("Id")

	if filters.Author != "" {
		query = query.Filter("Author=", filters.Author)
	}

	if filters.Title != "" {
		query = query.Filter("Title=", filters.Title)
	}
	return query
}
