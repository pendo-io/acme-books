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

func (brds *BookRepositoryDataStore) GetBooks(ctx context.Context)([]models.Book){

	it := brds.client.Run(ctx, datastore.NewQuery("Book"))
	return createBooks(it)
}

func (brds *BookRepositoryDataStore) GetBooksByTitle(ctx context.Context, title string)([]models.Book){

	query := datastore.NewQuery("Book").Filter("Title =", title)
	it := brds.client.Run(ctx, query)
	return createBooks(it)
}

func createBooks(it *datastore.Iterator) (books []models.Book){
	for {
		var b models.Book
		_, err := it.Next(&b)
		if err == iterator.Done {
			fmt.Println(err)
			break
		}
		books = append(books, b)
	}
	return books
}


