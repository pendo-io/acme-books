package controllers

import (
	"acme-books/helpers"
	"acme-books/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"cloud.google.com/go/datastore"
	"github.com/go-martini/martini"
)

type LibraryController struct{}

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
