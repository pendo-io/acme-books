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

type LibraryController struct {
	ctx    context.Context
	client *datastore.Client
}

func NewLibraryController() *LibraryController {
	lc := new(LibraryController)
	lc.ctx = context.Background()
	lc.client, _ = datastore.NewClient(lc.ctx, "acme-books")

	books := []models.Book{
		{Id: 1, Author: "George Orwell", Title: "1984", Borrowed: false},
		{Id: 2, Author: "George Orwell", Title: "Animal Farm", Borrowed: false},
		{Id: 3, Author: "Robert Jordan", Title: "Eye of the world", Borrowed: false},
		{Id: 4, Author: "Various", Title: "Collins Dictionary", Borrowed: false},
	}

	var keys []*datastore.Key

	for _, book := range books {
		keys = append(keys, datastore.IDKey("Book", book.Id, nil))
	}

	if _, err := lc.client.PutMulti(lc.ctx, keys, books); err != nil {
		fmt.Println(err)
	}

	return lc
}
func (lc LibraryController) Close() error {
	return lc.client.Close()
}

func (lc LibraryController) GetByKey(params martini.Params, w http.ResponseWriter) {

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	key := datastore.IDKey("Book", int64(id), nil)
	var book models.Book
	err = lc.client.Get(lc.ctx, key, &book)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJson(book, w)
}

func (lc LibraryController) ListAll(r *http.Request, w http.ResponseWriter) {

	var output []models.Book
	it := lc.client.Run(lc.ctx, datastore.NewQuery("Book"))
	for {
		var b models.Book
		_, err := it.Next(&b)
		if err == iterator.Done {
			fmt.Println(err)
			break
		}
		output = append(output, b)
	}

	writeJson(output, w)
}

func writeJson(item interface{}, w http.ResponseWriter) {
	jsonStr, err := json.MarshalIndent(item, "", "  ")

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonStr)
}
