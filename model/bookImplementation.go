package model

import (
	"cloud.google.com/go/datastore"
	"context"
	"errors"
	"fmt"
	"google.golang.org/api/iterator"
	"math/rand/v2"
)

type BookInterface interface {
	GetByKey(id int) (Book, error)
	ListAll(query *datastore.Query) []Book
	ChangeBorrowedStatus(id int, borrowed bool) error
	AddBook(book Book) (Book, error)
	Delete(id int) error
}

type BookImplementation struct {
	Ctx    context.Context
	Client *datastore.Client
}

func (bi BookImplementation) GetByKey(id int) (Book, error) {
	var book Book
	key := datastore.IDKey("Book", int64(id), nil)

	err := bi.Client.Get(bi.Ctx, key, &book)
	return book, err
}

func (bi BookImplementation) ListAll(query *datastore.Query) []Book {
	var output []Book

	it := bi.Client.Run(bi.Ctx, query)
	for {
		var b Book
		_, err := it.Next(&b)
		if errors.Is(err, iterator.Done) {
			fmt.Println(err)
			break
		}

		output = append(output, b)
	}

	return output
}

func (bi BookImplementation) ChangeBorrowedStatus(id int, borrowed bool) error {
	var book Book
	key := datastore.IDKey("Book", int64(id), nil)

	err := bi.Client.Get(bi.Ctx, key, &book)

	if err != nil {
		return err
	}

	if book.Borrowed == borrowed {
		var errorMessage string

		if book.Borrowed {
			errorMessage = "The Book has already been borrowed"
		} else {
			errorMessage = "The has already been returned"
		}

		return errors.New(errorMessage)
	}

	book.Borrowed = borrowed
	_, err = bi.Client.Put(bi.Ctx, key, &book)
	return err
}

func (bi BookImplementation) AddBook(book Book) (Book, error) {
	if book.Id == 0 {
		book.Id = rand.Int64()
	}

	key := datastore.IDKey("Book", book.Id, nil)
	_, err := bi.Client.Put(bi.Ctx, key, &book)
	return book, err
}

func (bi BookImplementation) Delete(id int) error {
	key := datastore.IDKey("Book", int64(id), nil)
	return bi.Client.Delete(bi.Ctx, key)
}
