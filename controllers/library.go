package controllers

import (
	"acme-books/helpers"
	"acme-books/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"cloud.google.com/go/datastore"
	"github.com/go-martini/martini"
)

type LibraryController struct{}

func (lc LibraryController) Create(r *http.Request, w http.ResponseWriter, ctx context.Context, client *datastore.Client) {

	author := r.FormValue("author")
	if author == "" {
		err := errors.New("author parameter is missing")
		writeError(err, http.StatusBadRequest, w)
		return
	}

	title := r.FormValue("title")
	if title == "" {
		err := errors.New("title parameter is missing")
		writeError(err, http.StatusBadRequest, w)
		return
	}

	var book models.Book
	book.Author = author
	book.Title = title
	book.Borrowed = false
	book, err := models.CreateBook(client, ctx, book)
	if err != nil {
		writeError(err, http.StatusBadRequest, w)
		return
	}

	jsonStr, err := json.MarshalIndent(book, "", "  ")

	if err != nil {
		writeError(err, http.StatusInternalServerError, w)
		return
	}

	writeResponse(jsonStr, w)
}

func (lc LibraryController) GetByKey(params martini.Params, w http.ResponseWriter, ctx context.Context, client *datastore.Client) {

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		writeError(err, http.StatusBadRequest, w)
		return
	}

	book, err := models.GetSingleBook(client, ctx, id)

	if err != nil {
		writeError(err, http.StatusInternalServerError, w)
		return
	}

	jsonStr, err := json.MarshalIndent(book, "", "  ")

	if err != nil {
		writeError(err, http.StatusInternalServerError, w)
		return
	}

	writeResponse(jsonStr, w)
}

func (lc LibraryController) BorrowBook(params martini.Params, w http.ResponseWriter, ctx context.Context, client *datastore.Client) {

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		writeError(err, http.StatusBadRequest, w)
		return
	}

	book, err := models.BorrowBook(client, ctx, id)

	if err, isBorrow := err.(helpers.IBorrowError); isBorrow && err.IsBorrowed() {
		writeError(err, http.StatusBadRequest, w)
		return
	}

	if err != nil {
		writeError(err, http.StatusBadRequest, w)
		return
	}

	jsonStr, err := json.MarshalIndent(book, "", "  ")

	if err != nil {
		writeError(err, http.StatusInternalServerError, w)
		return
	}

	writeResponse(jsonStr, w)
}

func (lc LibraryController) ReturnBook(params martini.Params, w http.ResponseWriter, ctx context.Context, client *datastore.Client) {

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		writeError(err, http.StatusBadRequest, w)
		return
	}

	book, err := models.ReturnBook(client, ctx, id)

	if err, isReturn := err.(helpers.IReturneBookError); isReturn && !err.IsBorrowed() {
		writeError(err, http.StatusBadRequest, w)
		return
	}

	if err != nil {
		writeError(err, http.StatusBadRequest, w)
		return
	}

	jsonStr, err := json.MarshalIndent(book, "", "  ")

	if err != nil {
		writeError(err, http.StatusInternalServerError, w)
		return
	}

	writeResponse(jsonStr, w)
}

func (lc LibraryController) ListAll(r *http.Request, w http.ResponseWriter, ctx context.Context, client *datastore.Client) {
	filters := r.URL.Query()
	output := models.GetAllBooks(filters, client, ctx)

	jsonStr, err := json.MarshalIndent(output, "", "  ")

	if err != nil {
		writeError(err, http.StatusInternalServerError, w)
		return
	}
	writeResponse(jsonStr, w)
}

func writeError(err error, status int, w http.ResponseWriter) {
	fmt.Println(err)
	w.WriteHeader(status)
}

func writeResponse(resp []byte, w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
