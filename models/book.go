package models

import (
	"context"
	"strconv"

	"cloud.google.com/go/datastore"
)

type Book struct {
	Id       int64
	Title    string `json:"title"`
	Author   string `json:"writer"`
	Borrowed bool   `json:"borrowed"`
}

type BookInt struct {
	client *datastore.Client
	ctx    context.Context
}

func NewBookInt(client *datastore.Client, ctx context.Context) *BookInt {
	return &BookInt{client, ctx}
}

func (bi BookInt) GetBookByKey(id int64) (book Book, err error) {
	key := datastore.IDKey("Book", id, nil)
	err = bi.client.Get(bi.ctx, key, &book)
	return book, err
}

func (bi BookInt) GetBooks(title, author, borrowed, sort string) ([]Book, error) {
	var books []Book
	query := datastore.NewQuery("Book")

	if title != "" {
		query = query.Filter("Title =", title)
	}

	if author != "" {
		query = query.Filter("Author =", author)
	}

	if borrowed != "" {
		if b, err := strconv.ParseBool(borrowed); err != nil {
			return books, err

		} else {
			query = query.Filter("Borrowed =", b)
		}
	}

	if sort == "" {
		sort = "Id"
	}

	query = query.Order(sort)

	_, err := bi.client.GetAll(bi.ctx, query, &books)
	return books, err
}

func (bi BookInt) PutBooks(books []Book) error {
	var keys []*datastore.Key

	for _, book := range books {
		keys = append(keys, datastore.IDKey("Book", book.Id, nil))
	}

	_, err := bi.client.PutMulti(bi.ctx, keys, books)

	return err
}

func (bi BookInt) PutBook(book Book) error {
	return bi.PutBooks([]Book{book})
}
