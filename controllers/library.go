package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"cloud.google.com/go/datastore"
	"github.com/go-martini/martini"

	"acme-books/models"
)

type LibraryController struct{}

func (lc LibraryController) GetByKey(params martini.Params, w http.ResponseWriter) {
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book, err := models.FindBookById(id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonStr, err := json.MarshalIndent(book, "", "  ")

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonStr)
}

func (lc LibraryController) ListAll(r *http.Request, w http.ResponseWriter) {
	params := r.URL.Query()

	filters := models.BookFilter{
		Author: params.Get("author"),
		Title: params.Get("title"),
	}

	books, err := models.ListBooks(filters)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonStr, err := json.MarshalIndent(books, "", "  ")

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonStr)
}

func (lc LibraryController) CreateBook(book models.Book, r *http.Request, w http.ResponseWriter) {
	book, err := models.AddBook(book)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonStr, err := json.MarshalIndent(book, "", "  ")

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonStr)
}

func (lc LibraryController) Borrow(params martini.Params, w http.ResponseWriter) {
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid id"))
		return
	}

	book, err := models.FindBookById(id)

	if book.Borrowed == true {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Book already borrowed"))
		return
	}

	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

	_, err = models.BorrowBook(book)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (lc LibraryController) Return(params martini.Params, w http.ResponseWriter) {
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid id"))
		return
	}

	book, err := models.FindBookById(id)

	if book.Borrowed == false {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Cannot return a book that is not borrowed"))
		return
	}

	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

	_, err = models.ReturnBook(book)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
