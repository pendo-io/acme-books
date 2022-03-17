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
	"acme-books/utils"
)

type LibraryController struct{}


func (lc LibraryController) GetByKey(client *datastore.Client , ctx context.Context, params martini.Params, w http.ResponseWriter) {
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		utils.ErrorResponse(w,http.StatusBadRequest, err)
		return
	}

	var book models.Book
	key := datastore.IDKey("Book", int64(id), nil)

	err = client.Get(ctx, key, &book)

	if err != nil {
		utils.ErrorResponse(w,http.StatusInternalServerError, err)
		return
	}

	jsonStr, err := json.MarshalIndent(book, "", "  ")

	if err != nil {
		utils.ErrorResponse(w,http.StatusInternalServerError, err)
		return
	}

	utils.OKResponse(w,jsonStr)
}

func (lc LibraryController) ListAll(client *datastore.Client , ctx context.Context, r *http.Request, w http.ResponseWriter) {

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

	if err != nil {
		utils.ErrorResponse(w,http.StatusInternalServerError, err)
		return
	}

	utils.OKResponse(w, jsonStr)
}
