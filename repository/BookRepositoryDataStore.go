package repository

import (
	"acme-books/models"
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

type BookRepositoryDataStore struct {
	client *datastore.Client
}

func NewBookRepositoryDataStore(client *datastore.Client) *BookRepositoryDataStore {
	if client == nil {
		panic("missing client")
	}

	return &BookRepositoryDataStore{client: client}
}

func (brds *BookRepositoryDataStore) GetBook(ctx context.Context, id int) (book models.Book, err error) {

	key := datastore.IDKey("Book", int64(id), nil)
	err = brds.client.Get(ctx, key, &book)

	return book, err
}

func (brds *BookRepositoryDataStore) GetBooks(ctx context.Context, filters map[string]string) []models.Book {

	query := buildQuery(filters)
	it := brds.client.Run(ctx, query)
	return createBooks(it)
}

func (brds *BookRepositoryDataStore) Lending(ctx context.Context, id int, borrow bool) (err error) {
	key := datastore.IDKey("Book", int64(id), nil)
	book := new(models.Book)
	if err = brds.client.Get(ctx, key, book); err != nil {
		return err
	}
	if borrow && book.Borrowed {
		return &BorrowedError{}
	}
	if !borrow && !book.Borrowed {
		return &ReturnedError{}
	}
	book.Borrowed = borrow

	if _, err := brds.client.Put(ctx, key, book); err != nil {
		return err
	}
	return
}

func (brds *BookRepositoryDataStore) AddBook(ctx context.Context, book models.Book) (id int, err error) {
	key := datastore.IncompleteKey("Book", nil)
	//this is dirty, but if book id and DS Id are the same
	key, err = brds.client.Put(ctx, key, &book)
	if err != nil {
		return 0, err
	}
	book.Id = key.ID
	_, err = brds.client.Put(ctx, key, &book)
	if err != nil {
		return 0, err
	}
	return int(book.Id), err

}

func createBooks(it *datastore.Iterator) (books []models.Book) {
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

func buildQuery(filters map[string]string) *datastore.Query {
	query := datastore.NewQuery("Book")
	for filter, value := range filters {
		if strings.EqualFold("title", filter) {
			query = query.Filter("Title=", value)
		}
		if strings.EqualFold("Author", filter) {
			query = query.Filter("Author=", value)
		}
	}
	return query
}
