package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"fmt"
	"reflect"
	"strconv"
	"errors"
	"math/rand"

	"cloud.google.com/go/datastore"
	"github.com/go-martini/martini"

	"acme-books/models"
)

type LibraryController struct{}

var (
	ctx context.Context
	client *datastore.Client
)

func (lc LibraryController) Init(ct context.Context, cl *datastore.Client) {
	ctx = ct
	client = cl

	deleteAllBooks()

	addBook("George Orwell", "1984", false)
	addBook("Robert Jordan", "Eye of the World", false)
	addBook("George Orwell", "Animal Farm", false)
	addBook("Various", "Collins Dictionary", false)
}

func addBook(author string, title string, borrowed bool) (*models.Book, *datastore.Key) {
	id := rand.Int63()
	book := &models.Book{
		Id: id,
		Author: author,
		Title: title,
		Borrowed: borrowed,
	}
	key := datastore.IDKey("Book", id, nil)
	key, _ = client.Put(ctx, key, book)
	return book, key
}

func deleteBook(bookID int64) error {
	return client.Delete(ctx, datastore.IDKey("Book", bookID, nil))
}

func deleteAllBooks() () {

	var books []*models.Book
	// Create a query to fetch all Task entities, ordered by "created".
	query := datastore.NewQuery("Book")
	keys, err := client.GetAll(ctx, query, &books)
	if err != nil {
		fmt.Println(err)
	}

	// Set the id field on each Task from the corresponding key.
	for _, key := range keys {
		err := deleteBook(key.ID)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (lc LibraryController) ListAll(r *http.Request, w http.ResponseWriter) {
	var books []*models.Book

	query := datastore.NewQuery("Book").Order("Id")
	keys, err := client.GetAll(ctx, query, &books)
	if err != nil {
		fmt.Println(err)
	}

	output := make([]models.Book, 0)
	for i, _ := range keys {
		match := true
		for qKey, qVal := range r.URL.Query() {
			if (getField(books[i], qKey) != qVal[0]) {
				match = false
			}
		}
		if (match) {
			output = append(output, *books[i])
		}
	}
	j, _ := json.MarshalIndent(output, "", "  ")
	w.Write(j)
}

func getField(b *models.Book, field string) string {
    r := reflect.ValueOf(b)
    f := reflect.Indirect(r).FieldByName(field)
    return fmt.Sprintf("%v", f)
}

func (lc LibraryController) Borrow(r *http.Request, w http.ResponseWriter, params martini.Params) (int, string) {
	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		return 400, "invalid id"
	}
	key := datastore.IDKey("Book", id, nil)

	var borrowedError error = errors.New("already borrowed")

	_, err = client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		var book models.Book
		if err := tx.Get(key, &book); err != nil {
			return err
		}
		if book.Borrowed {
			return borrowedError
		}
		book.Borrowed = true
		_, err := tx.Put(key, &book)
		return err
	})
	if errors.Is(err, borrowedError) {
		return 400, "already borrowed"
	}
	if err != nil  {
		return 400, fmt.Sprintf("%v", err)
	}
	return 200, "ok"
}

func (lc LibraryController) Return(r *http.Request, w http.ResponseWriter, params martini.Params) (int, string) {
	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		return 400, "invalid id"
	}
	fmt.Println(id)
	// Create a key using the given integer ID.
	key := datastore.IDKey("Book", id, nil)

	var borrowedError error = errors.New("already not borrowed")

	// In a transaction load each task, set done to true and store.
	_, err = client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		var book models.Book
		if err := tx.Get(key, &book); err != nil {
			return err
		}
		if !book.Borrowed {
			return borrowedError
		}
		book.Borrowed = false
		_, err := tx.Put(key, &book)
		return err
	})
	if errors.Is(err, borrowedError) {
		return 400, "already not borrowed"
	}
	if err != nil  {
		return 400, fmt.Sprintf("%v", err)
	}
	return 200, "ok"
}

func (lc LibraryController) Add(r *http.Request, w http.ResponseWriter) {
	book := new(models.Book)
	json.NewDecoder(r.Body).Decode(&book)
	book, key := addBook(book.Title, book.Author, book.Borrowed)

	type Output struct {
		Key		*datastore.Key
		Book	*models.Book
	}
	output := &Output{
		Key: key,
		Book: book,
	}
	j, _ := json.MarshalIndent(output, "", "  ")
	w.Write(j)
}
