package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-martini/martini"

	"acme-books/models"
)

type LibraryController struct{}

func (lc LibraryController) GetByKey(params martini.Params, w http.ResponseWriter) {
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book, err := models.FindBookById(id)

	if err != nil {
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
	params := r.URL.Query()

	filters := models.BookFilter{
		Author: params.Get("author"),
		Title: params.Get("title"),
	}

	books, err := models.ListBooks(filters)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonStr, err := json.MarshalIndent(books, "", "  ")

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonStr)
}
