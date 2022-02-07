package models

import (
	"cloud.google.com/go/datastore"
)

type Book struct {
	Id       int64
	Title    string `json:"title"`
	Author   string `json:"writer"`
	Borrowed bool   `json:"borrowed"`
}

func (book *Book) Borrow() error {
	return setBorrowed(book, true)
}

func (book *Book) Return() error {
	return setBorrowed(book, false)
}

func AddBook(book *Book) error {
	return create("Book", book)
}

func GetBook(id int64) (*Book, error) {
	filter := Filter{
		Operation: "Id =",
		Value : id,
	}
	if books, err := GetBooks(&filter); err != nil {
		return nil, err
	} else if len(books) != 1 {
		return nil, nil
	} else {
		return &books[0], nil
	}
}

func GetBooks(filter *Filter) ([]Book, error) {
	var books []Book
	listener := modelListener {
		itemFound: func(properties []datastore.Property) {
			if book, err := createBookFromProperties(properties); err != nil {
				return
			} else {
				books = append(books, book)
			}
		},
	}
	error := getAll("Book", filter, listener)
	return books, error
}

func setBorrowed(book *Book, borrowed bool) error {
	book.Borrowed = borrowed
	// book probably shouldn't do this update - would have a library struct to do this maybe?
	return update("Book", book.Id, book)
}

func createBookFromProperties(properties []datastore.Property) (Book, error) {
	book := Book{}
	error := datastore.LoadStruct(&book, properties)
	return book, error
}



