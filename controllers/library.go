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

type LibraryController struct {
	client *datastore.Client
	ctx    context.Context
}

func NewLibrary(client *datastore.Client, ctx context.Context, books []models.Book) (*LibraryController, error) {
	lc := LibraryController{client, ctx}
	if err := lc.bootstrapBooks(books); err != nil {
		fmt.Println("Problem bootstrapping library: %s", err)
		return nil, err
	}
	return &lc, nil
}

func (lc LibraryController) bootstrapBooks(books []models.Book) error {
	var keys []*datastore.Key

	for _, book := range books {
		keys = append(keys, datastore.IDKey("Book", book.Id, nil))
	}

	_, err := lc.client.PutMulti(lc.ctx, keys, books)

	return err
}

func (lc LibraryController) GetByKey(w http.ResponseWriter, params martini.Params) {
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var book models.Book
	key := datastore.IDKey("Book", int64(id), nil)

	err = lc.client.Get(lc.ctx, key, &book)

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

func (lc LibraryController) ListAll(w http.ResponseWriter, r *http.Request) {
	var (
		output []models.Book
		query  *datastore.Query
	)

	r.ParseForm()
	switch sortBy := r.Form.Get("sort"); sortBy {
	case "author", "title":
		query = datastore.NewQuery("Book").Order(strings.Title(sortBy))
	case "", "id":
		query = datastore.NewQuery("Book").Order("Id")
	default:
		fmt.Printf("Unknown sorting field: %s\n", sortBy)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	it := lc.client.Run(lc.ctx, query)
	for {
		var b models.Book
		if _, err := it.Next(&b); err == iterator.Done {
			fmt.Println(err)
			break
		} else if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			output = append(output, b)
		}
	}

	if jsonStr, err := json.MarshalIndent(output, "", "  "); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(jsonStr)
	}
}
