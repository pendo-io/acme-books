package model

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

var (
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

// BookFilter narrows down a book list. Each field matches anywhere in the
// book's title or author, ignoring case. Empty fields match every book.
type BookFilter struct {
	Title  string
	Author string
}

// IsEmpty reports whether the filter matches every book.
func (f BookFilter) IsEmpty() bool {
	return f.Title == "" && f.Author == ""
}

// Matches reports whether the book matches every field of the filter.
func (f BookFilter) Matches(book Book) bool {
	return containsFold(book.Title, f.Title) && containsFold(book.Author, f.Author)
}

type BookInterface interface {
	GetByKey(id int) (Book, error)
	ListAll() []Book
	ListMatching(filter BookFilter) []Book
	Borrow(id int) (Book, error)
	Return(id int) (Book, error)
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

// ListMatching returns the books that match the filter.
func (bi BookImplementation) ListMatching(filter BookFilter) []Book {
	books := bi.ListAll()

	if filter.IsEmpty() {
		return books
	}

	matches := []Book{}

	for _, book := range books {
		if filter.Matches(book) {
			matches = append(matches, book)
		}
	}

	return matches
}

// Borrow marks the book as borrowed. It returns ErrAlreadyBorrowed if the book
// is already borrowed, and datastore.ErrNoSuchEntity if there is no such book.
func (bi BookImplementation) Borrow(id int) (Book, error) {
	return bi.setBorrowed(id, true, ErrAlreadyBorrowed)
}

// Return marks the book as no longer borrowed. It returns ErrNotBorrowed if the
// book is not borrowed, and datastore.ErrNoSuchEntity if there is no such book.
func (bi BookImplementation) Return(id int) (Book, error) {
	return bi.setBorrowed(id, false, ErrNotBorrowed)
}

// Add saves the book under a newly generated key and returns it with its Id set.
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

// setBorrowed sets the book's Borrowed flag in a transaction, so that concurrent
// requests cannot both borrow, or both return, the same book. If the flag
// already has the requested value, it returns errUnchanged along with the book.
func (bi BookImplementation) setBorrowed(id int, borrowed bool, errUnchanged error) (Book, error) {
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

		if err := tx.Get(key, &book); err != nil {
			return err
		}

		if book.Borrowed == borrowed {
			return errUnchanged
		}

		book.Borrowed = borrowed

		_, err := tx.Put(key, &book)
		return err
	})

	return book, err
}

// containsFold reports whether substr is within s, ignoring case.
func containsFold(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
