package model

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

var (
	// ErrBookNotFound is returned when there is no book with the given id.
	ErrBookNotFound = errors.New("book not found")
	// ErrAlreadyBorrowed is returned when borrowing a book that is already borrowed.
	ErrAlreadyBorrowed = errors.New("book is already borrowed")
	// ErrNotBorrowed is returned when returning a book that is not borrowed.
	ErrNotBorrowed = errors.New("book is not borrowed")
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
	SetBorrowed(id int, borrowed bool) (Book, error)
	Add(book Book) (Book, error)
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

// SetBorrowed marks the book with the given id as borrowed or returned. The
// book is loaded, checked and saved in one transaction, so concurrent requests
// cannot borrow or return the same book twice. When the book is already in
// the requested state, it is returned as loaded alongside ErrAlreadyBorrowed
// or ErrNotBorrowed.
func (bi BookImplementation) SetBorrowed(id int, borrowed bool) (Book, error) {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, "acme-books")

	if err != nil {
		return Book{}, err
	}

	defer client.Close()

	var book Book
	key := datastore.IDKey("Book", int64(id), nil)

	_, err = client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		book = Book{}
		err := tx.Get(key, &book)

		if errors.Is(err, datastore.ErrNoSuchEntity) {
			return ErrBookNotFound
		}

		if err != nil {
			return err
		}

		if book.Borrowed == borrowed {
			if borrowed {
				return ErrAlreadyBorrowed
			}

			return ErrNotBorrowed
		}

		book.Borrowed = borrowed
		_, err = tx.Put(key, &book)

		return err
	})

	return book, err
}

// Add saves book under a newly generated key and returns it with its Id set
// to that key's ID.
func (bi BookImplementation) Add(book Book) (Book, error) {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, "acme-books")

	if err != nil {
		return Book{}, err
	}

	defer client.Close()

	keys, err := client.AllocateIDs(ctx, []*datastore.Key{datastore.IncompleteKey("Book", nil)})

	if err != nil {
		return Book{}, err
	}

	book.Id = keys[0].ID

	if _, err := client.Put(ctx, keys[0], &book); err != nil {
		return Book{}, err
	}

	return book, nil
}
