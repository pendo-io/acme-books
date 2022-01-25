package models

import (
	"context"
	"fmt"

	"cloud.google.com/go/datastore"
)

type Book struct {
	Id       int64
	Title    string `json:"title"`
	Author   string `json:"writer"`
	Borrowed bool   `json:"borrowed"`
}

type BookFilter struct {
	Author string
	Title string
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

func ListBooks(filters BookFilter) ([]Book, error) {
	ctx := context.Background()

	var books []Book

	query := datastore.NewQuery("Book").Order("Id")

	if filters.Author != "" {
		query = query.Filter("Author=", filters.Author)
	}

	if filters.Title != "" {
		query = query.Filter("Title=", filters.Title)
	}

	_, err := client.GetAll(ctx, query, &books)

	if err != nil {
		return books, err
	}

	return books, nil
}
