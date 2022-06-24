package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonStr)
}

func (lc LibraryController) PostBook(r *http.Request, w http.ResponseWriter) {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "acme-books")

	defer client.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error reading POST JSON body"))
		return
	}

	book := models.Book{}
	err = json.Unmarshal(body, &book)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Error parsing POST JSON into a Book: %s", string(body))))
		return
	}

	// Get the next ID
	count, err := client.Count(ctx, datastore.NewQuery("Book"))

	book.Id = int64(count) + 1
	key := datastore.Key{Kind: "Book", ID: book.Id}
	returnKey, err := client.Put(ctx, &key, &book)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error storing Book to datastore with key: %v", key)))
		return
	}

	if returnKey.ID != key.ID {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Mismatch in stored key: %v != %v", key, returnKey)))
		return
	}

	book.Id = returnKey.ID
	jsonStr, err := json.MarshalIndent(book, "", "  ")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error marshaling book into JSON: %v", book)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonStr)
}
