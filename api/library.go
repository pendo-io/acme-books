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

// getVisitorID extracts the visitor identifier from the request.
// Falls back to "anonymous" when no X-Visitor-ID header is present.
func getVisitorID(r *http.Request) string {
	if id := r.Header.Get("X-Visitor-ID"); id != "" {
		return id
	}
	return "anonymous"
}

// getAccountID extracts the account identifier from the request.
// Falls back to "acme-books" when no X-Account-ID header is present.
func getAccountID(r *http.Request) string {
	if id := r.Header.Get("X-Account-ID"); id != "" {
		return id
	}
	return "acme-books"
}

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

	// Filter support via query parameters (author, title, borrowed).
	authorFilter := r.URL.Query().Get("author")
	titleFilter := r.URL.Query().Get("title")
	borrowedFilter := r.URL.Query().Get("borrowed")

	totalCount := len(books)
	hasFilters := authorFilter != "" || titleFilter != "" || borrowedFilter != ""

	if hasFilters {
		var filtered []model.Book
		for _, book := range books {
			if authorFilter != "" && !strings.Contains(strings.ToLower(book.Author), strings.ToLower(authorFilter)) {
				continue
			}
			if titleFilter != "" && !strings.Contains(strings.ToLower(book.Title), strings.ToLower(titleFilter)) {
				continue
			}
			if borrowedFilter != "" {
				isBorrowed := borrowedFilter == "true"
				if book.Borrowed != isBorrowed {
					continue
				}
			}
			filtered = append(filtered, book)
		}
		books = filtered

		var filterFields []string
		var filterValues []string
		if authorFilter != "" {
			filterFields = append(filterFields, "author")
			filterValues = append(filterValues, authorFilter)
		}
		if titleFilter != "" {
			filterFields = append(filterFields, "title")
			filterValues = append(filterValues, titleFilter)
		}
		if borrowedFilter != "" {
			filterFields = append(filterFields, "borrowed")
			filterValues = append(filterValues, borrowedFilter)
		}

		// Pendo track: catalog filtered
		pendo.Track("book_catalog_filtered", getVisitorID(r), getAccountID(r), map[string]interface{}{
			"filter_field":      strings.Join(filterFields, ","),
			"filter_value":      strings.Join(filterValues, ","),
			"results_count":     len(books),
			"total_books_count": totalCount,
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

func (l Library) Borrow(params martini.Params, r *http.Request, w http.ResponseWriter) {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book, err := model.BookImplementation{}.Borrow(id)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Pendo track: book borrowed successfully
	pendo.Track("book_borrowed", getVisitorID(r), getAccountID(r), map[string]interface{}{
		"book_id":                strconv.FormatInt(book.Id, 10),
		"book_title":             book.Title,
		"book_author":            book.Author,
		"borrowed_status_before": false,
	})

	w.WriteHeader(http.StatusNoContent)
}

func (l Library) Return(params martini.Params, r *http.Request, w http.ResponseWriter) {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book, err := model.BookImplementation{}.Return(id)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Pendo track: book returned successfully
	pendo.Track("book_returned", getVisitorID(r), getAccountID(r), map[string]interface{}{
		"book_id":                strconv.FormatInt(book.Id, 10),
		"book_title":             book.Title,
		"book_author":            book.Author,
		"borrowed_status_before": true,
	})

	w.WriteHeader(http.StatusNoContent)
}

func (l Library) Create(r *http.Request, w http.ResponseWriter) {
	var book model.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
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

	// Pendo track: book added successfully
	pendo.Track("book_added", getVisitorID(r), getAccountID(r), map[string]interface{}{
		"book_id":     strconv.FormatInt(createdBook.Id, 10),
		"book_title":  createdBook.Title,
		"book_author": createdBook.Author,
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
