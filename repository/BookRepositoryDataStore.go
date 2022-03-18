package repository

import (
	"acme-books/models"
	"context"
	"fmt"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

type BookRepositoryDataStore struct{
	client *datastore.Client
}

func NewBookRepositoryDataStore(client *datastore.Client) *BookRepositoryDataStore{
	if client == nil {
		panic("missing client")
	}

	return &BookRepositoryDataStore{client: client}
}

func (brds *BookRepositoryDataStore) GetBook(ctx context.Context, id int)(book models.Book, err error){

	key := datastore.IDKey("Book", int64(id), nil)
	err = brds.client.Get(ctx, key, &book)

	return book, err
}

func (brds *BookRepositoryDataStore) GetBooks(ctx context.Context)(books []models.Book, err error){

	it := brds.client.Run(ctx, datastore.NewQuery("Book"))
	for {
		var b models.Book
		_, err := it.Next(&b)
		if err == iterator.Done {
			fmt.Println(err)
			break
		}
		books = append(books, b)
	}
	return books, err
}
