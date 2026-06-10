package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-martini/martini"

	"acme-books/model"
	"acme-books/pendo"
)

type Library struct{}

func (l Library) GetByKey(params martini.Params, w http.ResponseWriter) {
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book, err := model.BookImplementation{}.GetByKey(id)

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

func (l Library) ListAll(r *http.Request, w http.ResponseWriter) {
	books := model.BookImplementation{}.ListAll()

	// Check for filter query parameters
	filterField := r.URL.Query().Get("filterField")
	filterValue := r.URL.Query().Get("filterValue")
	sortOrder := r.URL.Query().Get("sortOrder")

	if filterField != "" && filterValue != "" {
		var filtered []model.Book
		lowerValue := strings.ToLower(filterValue)
		for _, b := range books {
			switch strings.ToLower(filterField) {
			case "title":
				if strings.Contains(strings.ToLower(b.Title), lowerValue) {
					filtered = append(filtered, b)
				}
			case "author", "writer":
				if strings.Contains(strings.ToLower(b.Author), lowerValue) {
					filtered = append(filtered, b)
				}
			case "borrowed":
				borrowedStr := strconv.FormatBool(b.Borrowed)
				if borrowedStr == lowerValue {
					filtered = append(filtered, b)
				}
			}
		}
		books = filtered

		// Pendo Track: book_catalog_searched
		go pendo.Track("book_catalog_searched", "system", "system", map[string]interface{}{
			"filterField":  filterField,
			"filterValue":  filterValue,
			"resultsCount": len(books),
			"sortOrder":    sortOrder,
		})
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

func (l Library) BorrowBook(params martini.Params, w http.ResponseWriter) {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book, err := model.BookImplementation{}.GetByKey(id)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if book.Borrowed {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	previousBorrowedStatus := book.Borrowed
	book.Borrowed = true

	err = model.BookImplementation{}.Update(id, book)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Pendo Track: book_borrowed
	go pendo.Track("book_borrowed", "system", "system", map[string]interface{}{
		"bookId":                 strconv.FormatInt(book.Id, 10),
		"bookTitle":              book.Title,
		"bookAuthor":             book.Author,
		"previousBorrowedStatus": strconv.FormatBool(previousBorrowedStatus),
	})

	w.WriteHeader(http.StatusNoContent)
}

func (l Library) ReturnBook(params martini.Params, w http.ResponseWriter) {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book, err := model.BookImplementation{}.GetByKey(id)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !book.Borrowed {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book.Borrowed = false

	err = model.BookImplementation{}.Update(id, book)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Pendo Track: book_returned
	go pendo.Track("book_returned", "system", "system", map[string]interface{}{
		"bookId":     strconv.FormatInt(book.Id, 10),
		"bookTitle":  book.Title,
		"bookAuthor": book.Author,
	})

	w.WriteHeader(http.StatusNoContent)
}

func (l Library) CreateBook(r *http.Request, w http.ResponseWriter) {
	var book model.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	createdBook, err := model.BookImplementation{}.Create(book)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Pendo Track: book_created
	go pendo.Track("book_created", "system", "system", map[string]interface{}{
		"bookId":     strconv.FormatInt(createdBook.Id, 10),
		"bookTitle":  createdBook.Title,
		"bookAuthor": createdBook.Author,
		"borrowed":   strconv.FormatBool(createdBook.Borrowed),
	})

	jsonStr, err := json.MarshalIndent(createdBook, "", "  ")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonStr)
}
