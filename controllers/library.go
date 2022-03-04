package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"cloud.google.com/go/datastore"
	"github.com/go-martini/martini"
	"google.golang.org/api/iterator"

	"acme-books/models"
)

type LibraryController struct{}

func (lc LibraryController) GetByKey(params martini.Params, w http.ResponseWriter) {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "acme-books")

	defer client.Close()

	id, err := strconv.Atoi(params["id"])

	if CheckError(err, http.StatusBadRequest, w) {
		return
	}

	var book models.Book
	key := datastore.IDKey("Book", int64(id), nil)

	err = client.Get(ctx, key, &book)

	if CheckError(err, http.StatusInternalServerError, w) {
		return
	}

	jsonStr, err := json.MarshalIndent(book, "", "  ")

	if CheckError(err, http.StatusInternalServerError, w) {
		return
	}

	PostResponse(jsonStr, w)
}

func (lc LibraryController) ListAll(r *http.Request, w http.ResponseWriter) {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "acme-books")

	defer client.Close()

	var output []models.Book

	it := client.Run(ctx, datastore.NewQuery("Book"))
	for {
		var b models.Book
		_, err := it.Next(&b)
		if err == iterator.Done {
			fmt.Println(err)
			break
		}
		output = append(output, b)
	}

	jsonStr, err := json.MarshalIndent(output, "", "  ")

	if CheckError(err, http.StatusInternalServerError, w) {
		return
	}

	PostResponse(jsonStr, w)
}

func CheckError(err error, statuscode int, w http.ResponseWriter) bool {
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(statuscode)
		return true
	}
	return false
}

func PostResponse(jsonStr []byte, w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Write(jsonStr)
	return
}
