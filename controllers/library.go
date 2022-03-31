package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"cloud.google.com/go/datastore"
	"github.com/go-martini/martini"

	"acme-books/models"
)

type LibraryController struct {
	bi *models.BookInt
}

func NewLibrary(client *datastore.Client, ctx context.Context, books []models.Book) (*LibraryController, error) {
	bi := models.NewBookInt(client, ctx)
	lc := LibraryController{bi}
	if err := lc.bootstrapBooks(books); err != nil {
		fmt.Println("Problem bootstrapping library: ", err)
		return nil, err
	}
	return &lc, nil
}

func (lc LibraryController) bootstrapBooks(books []models.Book) error {
	return lc.bi.PutBooks(books)
}

func (lc LibraryController) Get(w http.ResponseWriter, params martini.Params) {
	if id, err := strconv.Atoi(params["id"]); err != nil {
		fmt.Println("Bad id: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if book, err := lc.bi.GetBookByKey(int64(id)); err != nil {
		fmt.Println("Error getting book: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if jsonStr, err := json.MarshalIndent(book, "", "  "); err != nil {
		fmt.Println("Error serializing: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(jsonStr)
	}
}

func (lc LibraryController) List(w http.ResponseWriter, r *http.Request) {
	var output []models.Book
	var err error

	r.ParseForm()

	title := r.Form.Get("title")
	author := r.Form.Get("author")
	borrowed := r.Form.Get("borrowed")

	switch sortBy := r.Form.Get("sort"); sortBy {
	case "author", "title", "id", "":
		output, err = lc.bi.GetBooks(title, author, borrowed, strings.Title(sortBy))
	default:
		fmt.Println("Unknown sorting field: ", sortBy)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err != nil {
		fmt.Println("Error getting books: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if jsonStr, err := json.MarshalIndent(output, "", "  "); err != nil {
		fmt.Println("Error serializing: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(jsonStr)
	}
}

// Return a handler to borrow or return a book.  TODO: sanitize inputs
func (lc LibraryController) BorrowOrReturn(borrow bool) martini.Handler {
	return func(w http.ResponseWriter, params martini.Params) {
		var book models.Book

		state := "return"
		if borrow {
			state = "borrow"
		}
		if id, err := strconv.Atoi(params["id"]); err != nil {
			fmt.Println("Bad id: ", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		} else if book, err = lc.bi.GetBookByKey(int64(id)); err == datastore.ErrNoSuchEntity {
			fmt.Println("Book not found: ", id)
			w.WriteHeader(http.StatusBadRequest)
			return
		} else if err != nil {
			fmt.Println("Error getting book: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else if book.Borrowed == borrow {
			fmt.Printf("Book already %sed: %d\n", state, id)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		book.Borrowed = borrow
		if err := lc.bi.PutBook(book); err != nil {
			fmt.Printf("Error %sing book: %s\n", state, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Printf("Book %sed: %d\n", state, book.Id)
		w.WriteHeader(http.StatusNoContent)
	}
}
