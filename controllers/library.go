package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"acme-books/models"
	"cloud.google.com/go/datastore"
	"github.com/go-martini/martini"
)

type LibraryController struct{}

func (lc LibraryController) GetByKey(params martini.Params, w http.ResponseWriter, client *datastore.Client, ctx context.Context) {
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var book models.Book
	key := datastore.IDKey("Book", int64(id), nil)

	err = client.Get(ctx, key, &book)

	if err != nil {
		fmt.Println(err)
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

func (lc LibraryController) ListAll(r *http.Request, w http.ResponseWriter, client *datastore.Client, ctx context.Context) {
	if output, err := models.ListBooks(client, ctx,"Id", r.URL.Query()); err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	} else if jsonStr, err := json.MarshalIndent(output, "", "  "); err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(jsonStr)
	}
}

func (lc LibraryController) BorrowByKey(params martini.Params, w http.ResponseWriter, client *datastore.Client, ctx context.Context) {
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var book models.Book
	key := datastore.IDKey("Book", int64(id), nil)

	err = client.Get(ctx, key, &book)

	if book.Borrowed {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		book.Borrowed = true
	}

	if 	_, err = client.Put(ctx, key, &book); err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func (lc LibraryController) ReturnByKey(params martini.Params, w http.ResponseWriter, client *datastore.Client, ctx context.Context) {
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var book models.Book
	key := datastore.IDKey("Book", int64(id), nil)

	err = client.Get(ctx, key, &book)

	if !book.Borrowed {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		book.Borrowed = false
	}

	if 	_, err = client.Put(ctx, key, &book); err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func (lc LibraryController) CreateBook(w http.ResponseWriter, client *datastore.Client, ctx context.Context, bookBinding models.Book) {
	key := datastore.IncompleteKey("Book", nil)
	newKey, err := client.Put(ctx, key, &bookBinding)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bookBinding.Id = newKey.ID
	if body, err := json.Marshal(bookBinding); err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.Write(body)
		w.WriteHeader(http.StatusOK)
	}
}