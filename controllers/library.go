package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"cloud.google.com/go/datastore"
	"github.com/go-martini/martini"

	"acme-books/models"
)

type LibraryController struct{}

func (lc LibraryController) GetByKey(params martini.Params, w http.ResponseWriter, ctx context.Context, client *datastore.Client) {

	fmt.Println(params)
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
