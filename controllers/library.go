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

	id, err := strconv.Atoi(params["id"])

	writeError(err, http.StatusBadRequest, w)

	var book models.Book
	key := datastore.IDKey("Book", int64(id), nil)

	err = client.Get(ctx, key, &book)

	writeError(err, http.StatusInternalServerError, w)

	jsonStr, err := json.MarshalIndent(book, "", "  ")

	writeError(err, http.StatusInternalServerError, w)

	writeResponse(jsonStr, w)
}

func (lc LibraryController) ListAll(r *http.Request, w http.ResponseWriter, ctx context.Context, client *datastore.Client) {

	output := models.GetAllBooks(client, ctx)

	jsonStr, err := json.MarshalIndent(output, "", "  ")

	writeError(err, http.StatusInternalServerError, w)

	writeResponse(jsonStr, w)
}

func writeError(err error, status int, w http.ResponseWriter) {
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(status)
		return
	}
}

func writeResponse(resp []byte, w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
