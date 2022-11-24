package models

import (
	"context"

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
	Query(title string) ([]Book, error)
	PutMulti([]Book) error
	Put(Book, bool) (int64, error)
	CloseClientConnection()
}

type BookImplementation struct {
	client *datastore.Client
	context context.Context
}

func NewBookImplementation() BookInterface {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "acme-books")

	return BookImplementation{client: client, context: ctx}
}

func (bookImpl BookImplementation) CloseClientConnection () {
	defer bookImpl.client.Close()
}

func (bookImpl BookImplementation) GetByKey(id int) (Book, error) {
	var book Book
	key := datastore.IDKey("Book", int64(id), nil)

	err := bookImpl.client.Get(bookImpl.context, key, &book)

	if err != nil {
		return Book{}, err
	}
    
	return book, nil
}

func (bookImpl BookImplementation) Query(title string) ([]Book, error) {
	var books []Book

	query := datastore.NewQuery("Book").Order("Id")

	if len(title) > 0 {
		query = query.Filter("Title =", title)
	}

	it := bookImpl.client.Run(bookImpl.context, query)

	for {
		var b Book
		_, err := it.Next(&b)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return []Book{}, err
		}
		books = append(books, b)
	}
	return books, nil
}

func (bookImpl BookImplementation) Put(book Book, insert bool) (int64, error) {
	var key *datastore.Key
	if insert {
		key = datastore.IncompleteKey("Book", nil)
	} else {
		key = datastore.IDKey("Book", book.Id, nil)
	}

	if newKey, err := bookImpl.client.Put(bookImpl.context, key, &book); err != nil {
		return 0, err
	} else {
		return newKey.ID, nil
	}
}

func (bookImpl BookImplementation) PutMulti(books []Book) error {
	var keys []*datastore.Key

	for _, book := range books {
		keys = append(keys, datastore.IDKey("Book", book.Id, nil))
	}

	if _, err := bookImpl.client.PutMulti(bookImpl.context, keys, books); err != nil {
		return err
	}

	return nil
}
