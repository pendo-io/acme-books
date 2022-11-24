package controllers

import (
	"acme-books/models"
	"encoding/json"
	"github.com/go-martini/martini"
	"log"
	"net/http"
	"strconv"
)

type LibraryController struct{}

func (lc LibraryController) GetByKey(params martini.Params, w http.ResponseWriter) {
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		handleError("Invalid id", http.StatusBadRequest, w)
		return
	}

	book, err := models.GetBookById(int64(id))

	if err != nil {
		handleError("Invalid id", http.StatusBadRequest, w)
		return
	}

	showBookInfo(book, w)
}

func (lc LibraryController) ListAll(r *http.Request, w http.ResponseWriter) {
	books, err := models.GetBooks(getFilters(r), "Id")

	if err != nil {
		handleError("Internal error", http.StatusInternalServerError, w)
		return
	}

	showAllBooksInfo(books, w)
}

func (lc LibraryController) Borrow(params martini.Params, w http.ResponseWriter) {
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		handleError("Invalid id", http.StatusBadRequest, w)
		return
	}

	book, err := models.GetBookById(int64(id))

	if err != nil {
		handleError("Invalid id", http.StatusBadRequest, w)
		return
	}

	if book.Borrowed {
		handleError("Already borrowed", http.StatusBadRequest, w)
		return
	}

	if err := book.Borrow(); err != nil {
		handleError("Internal error", http.StatusInternalServerError, w)
		return
	}

	showBookInfo(book, w)
}

func (lc LibraryController) Return(params martini.Params, w http.ResponseWriter) {
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		handleError("Invalid id", http.StatusBadRequest, w)
		return
	}

	book, err := models.GetBookById(int64(id))

	if err != nil {
		handleError("Invalid id", http.StatusBadRequest, w)
		return
	}

	if !book.Borrowed {
		handleError("Not borrowed", http.StatusBadRequest, w)
		return
	}

	if err := book.Return(); err != nil {
		handleError("Internal error", http.StatusInternalServerError, w)
		return
	}

	showBookInfo(book, w)
}

func (lc LibraryController) Add(r *http.Request, w http.ResponseWriter) {
	r.ParseForm()

	id, err := strconv.ParseInt(r.Form.Get("id"), 10, 64)
	if err != nil {
		handleError("Invalid/missing id", http.StatusBadRequest, w)
		return
	}
	if id == 0 {
		handleError("Missing id", http.StatusBadRequest, w)
		return
	}

	book, err := models.GetBookById(id)

	if err != nil && book.Id != 0 {
		handleError("Internal error", http.StatusInternalServerError, w)
		return
	}
	if book.Id != 0 {
		handleError("Book already exists with id", http.StatusBadRequest, w)
		return
	}

	writer := r.Form.Get("writer")
	if writer == "" {
		handleError("Missing writer", http.StatusBadRequest, w)
		return
	}

	title := r.Form.Get("title")
	if title == "" {
		handleError("Missing title", http.StatusBadRequest, w)
		return
	}

	book.Id = id
	book.Author = writer
	book.Title = title
	book.Borrowed = false

	if err := book.Save(); err != nil {
		handleError("Internal error", http.StatusInternalServerError, w)
		return
	}

	showBookInfo(book, w)
}

func getFilters(r *http.Request) models.Filters {
	filters := models.Filters{}
	query := r.URL.Query()
	title := query.Get("title")
	if title != "" {
		filters.Title = title
	}
	writer := query.Get("writer")
	if writer != "" {
		filters.Writer = writer
	}
	available := query.Get("available")
	if available == "true" {
		filters.Available = true
	}
	return filters
}

func handleError(message string, errorCode int, w http.ResponseWriter) {
	jsonStr, err := json.MarshalIndent(message, "", "  ")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(errorCode)
	w.Write(jsonStr)
	return
}

func showBookInfo(book models.Book, w http.ResponseWriter) {
	jsonStr, err := json.MarshalIndent(book, "", "  ")

	if err != nil {
		handleError("Internal error", http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonStr)
}

func showAllBooksInfo(books []models.Book, w http.ResponseWriter) {
	jsonStr, err := json.MarshalIndent(books, "", "  ")

	if err != nil {
		handleError("Internal error", http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonStr)
}
