package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"errors"

	"cloud.google.com/go/datastore"
	"github.com/go-martini/martini"
	"google.golang.org/api/iterator"

	"acme-books/models"
)

type LibraryController struct{}

func (lc LibraryController) Create(r *http.Request, w http.ResponseWriter) {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "acme-books")

	defer client.Close()

	var book models.Book

	err := json.NewDecoder(r.Body).Decode(&book)

	if CheckError(err, http.StatusInternalServerError, w) {
		return
	}

	_, err = client.Put(r.Context(), datastore.IncompleteKey("Book", nil), &book)

	if CheckError(err, http.StatusInternalServerError, w) {
		return
	}

	jsonStr, err := json.MarshalIndent(book, "", "  ")

	if CheckError(err, http.StatusInternalServerError, w) {
		return
	}

	PostResponse(jsonStr, w)
}

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

	it := client.Run(ctx, datastore.NewQuery("Book").Order("Id"))
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

func (lc LibraryController) ListAllSearch(params martini.Params, r *http.Request, w http.ResponseWriter) {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "acme-books")

	defer client.Close()

	var output []models.Book

	it := client.Run(ctx, datastore.NewQuery("Book").Filter("Title <=", params["search"]))
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

func (lc LibraryController) BorrowBook(params martini.Params, w http.ResponseWriter) {
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

	fmt.Println(book.Borrowed)

	if book.Borrowed == false {
		book.Borrowed = true
	} else {
		err = errors.New("400: already borrowed")
	}

	fmt.Println(err)
	if CheckError(err, http.StatusInternalServerError, w) {
		return
	}

	if _, err := client.Put(ctx, key, &book); err != nil {
		fmt.Println(err)
	}

	if CheckError(err, http.StatusInternalServerError, w) {
		return
	}
}

func (lc LibraryController) ReturnBook(params martini.Params, w http.ResponseWriter) {
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

	if book.Borrowed == true {
		book.Borrowed = false
	} else {
		err = errors.New("400: returned")
	}

	fmt.Println(err)
	if CheckError(err, http.StatusInternalServerError, w) {
		return
	}

	if _, err := client.Put(ctx, key, book); err != nil {
		fmt.Println(err)
	}

	if CheckError(err, http.StatusInternalServerError, w) {
		return
	}
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
