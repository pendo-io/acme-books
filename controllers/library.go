package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"cloud.google.com/go/datastore"
	"github.com/go-martini/martini"

	"acme-books/models"
	"acme-books/service"
)

type LibraryController struct{}

func (lc LibraryController) ListAll(r *http.Request, w http.ResponseWriter) {
	qs := r.URL.Query()
	author, title := qs.Get("author"), qs.Get("title")

	query := datastore.NewQuery("Book").
				Order("Id")

	if author != "" {
		query = query.Filter("Author=", author)
	}
	if title != "" {
		query = query.Filter("Title=", title)
	}

	output := service.GetFromStore(query)

	j, _ := json.MarshalIndent(output, "", "  ")
	w.Write(j)
}

func (lc LibraryController) Borrow(params martini.Params, w http.ResponseWriter) {
	book, err := getBookById(params)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	if book.Borrowed {
		w.WriteHeader(400)
		w.Write([]byte("Already borrowed"))
		return
	}

	book.Borrowed = true
	service.AddOrUpdateStore(book)

	w.WriteHeader(204)
}

func (lc LibraryController) Return(params martini.Params, w http.ResponseWriter) {
	book, err := getBookById(params)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	if !book.Borrowed {
		w.WriteHeader(400)
		w.Write([]byte("Not yet borrowed"))
		return
	}

	book.Borrowed = false
	service.AddOrUpdateStore(book)

	w.WriteHeader(204)
}

func getBookById(params martini.Params) (*models.Book, error) {
	id := params["id"]

	if id == "" {
		return nil, errors.New("No id provided")
	}

	i, err := strconv.Atoi(id)
	if err != nil {
		return nil, errors.New("Non numeric id")
	}

	query := datastore.NewQuery("Book").
		Filter("Id=", i)

	book := service.GetFromStore(query)

	if book == nil || len(book) != 1 {
		return nil, errors.New("No single book found")
	}

	return &book[0], nil
}
