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
	"google.golang.org/api/iterator"

	"acme-books/models"
)

type LibraryController struct{}

func (lc LibraryController) GetByKey(params martini.Params, w http.ResponseWriter) {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "acme-books")

	defer client.Close()

	id, err := strconv.Atoi(params["id"])

	HandleError(err, w)

	var book models.Book
	key := datastore.IDKey("Book", int64(id), nil)

	err = client.Get(ctx, key, &book)

	HandleError(err, w)

	jsonStr, err := json.MarshalIndent(book, "", "  ")

	HandleError(err, w)


	w.WriteHeader(http.StatusOK)
	w.Write(jsonStr)
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

	HandleError(err, w)

	w.WriteHeader(http.StatusOK)
	w.Write(jsonStr)
}

func HandleError (err error, w http.ResponseWriter) {
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (lc LibraryController) SearchByName(params martini.Params, w http.ResponseWriter) {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "acme-books")

	defer client.Close()

	name := strings.ReplaceAll(params["name"], "%20"," ")


	fmt.Println(name)


	//HandleError(err, w)

	var output []models.Book

	q := datastore.NewQuery("Book").Filter("Title=", name)

	it := client.Run(ctx,q)

	for {
		var b models.Book
		_, err := it.Next(&b)
		fmt.Println(b)
		if err == iterator.Done {
			fmt.Println(err)
			break
		}
		output = append(output, b)
	}

	jsonStr, err := json.MarshalIndent(output, "", "  ")

	HandleError(err, w)

	w.WriteHeader(http.StatusOK)
	w.Write(jsonStr)
}

func (lc LibraryController) Borrow (params martini.Params, w http.ResponseWriter) {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "acme-books")

	defer client.Close()

	id, err := strconv.Atoi(params["id"])

	HandleError(err, w)

	var book models.Book
	key := datastore.IDKey("Book", int64(id), nil)

	err = client.Get(ctx, key, &book)

	switch {
	case err != nil:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid id"))
	case book.Borrowed:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("already borrowed"))
	default:
		book.Borrowed = true
		_, err = client.Put(ctx, key, &book)
		HandleError(err, w)

		jsonStr, err := json.MarshalIndent(book, "", "  ")

		HandleError(err, w)

		w.WriteHeader(http.StatusOK)
		w.Write(jsonStr)
	}


}

func (lc LibraryController) Return (params martini.Params, w http.ResponseWriter) {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "acme-books")

	defer client.Close()

	id, err := strconv.Atoi(params["id"])

	HandleError(err, w)

	var book models.Book
	key := datastore.IDKey("Book", int64(id), nil)

	err = client.Get(ctx, key, &book)

	switch {
	case err != nil:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid id"))
		return
	case !book.Borrowed:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("already returned"))
	default:
		book.Borrowed = false
		_, err = client.Put(ctx, key, &book)
		HandleError(err, w)

		jsonStr, err := json.MarshalIndent(book, "", "  ")

		HandleError(err, w)

		w.WriteHeader(http.StatusOK)
		w.Write(jsonStr)
	}
}

func (lc LibraryController) CreateBook (r *http.Request, w http.ResponseWriter) {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "acme-books")

	defer client.Close()

	var book models.Book


	err := json.NewDecoder(r.Body).Decode(&book)
	HandleError(err, w)

	key := datastore.IncompleteKey("Book", nil)
	fmt.Println(key)

	newKey, err := client.Put(ctx, key, &book)
	book.Id = newKey.ID
	fmt.Println(book)
	HandleError(err, w)


	jsonStr, err := json.MarshalIndent(book, "", "  ")

	HandleError(err, w)

	w.WriteHeader(http.StatusOK)
	w.Write(jsonStr)

}
