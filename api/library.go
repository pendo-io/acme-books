package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"cloud.google.com/go/datastore"
	"github.com/go-martini/martini"

	"acme-books/model"
	"acme-books/pendo"
)

type Library struct{}

func (l Library) GetByKey(params martini.Params, w http.ResponseWriter) {
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		fmt.Println(err)
		trackEvent("book_lookup_failed", pendo.Properties{
			"requestedId": params["id"],
			"errorType":   "invalid_id",
		})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book, err := model.BookImplementation{}.GetByKey(id)

	if err != nil {
		fmt.Println(err)
		if errors.Is(err, datastore.ErrNoSuchEntity) {
			trackEvent("book_lookup_failed", pendo.Properties{
				"requestedId": params["id"],
				"errorType":   "not_found",
			})
		}
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
	query := r.URL.Query()
	filter := model.BookFilter{
		Title:  strings.TrimSpace(query.Get("title")),
		Author: strings.TrimSpace(query.Get("author")),
	}

	books := model.BookImplementation{}.ListMatching(filter)

	jsonStr, err := json.MarshalIndent(books, "", "  ")

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !filter.IsEmpty() {
		trackEvent("book_catalog_filtered", pendo.Properties{
			"filters":        filtersProperty(filter),
			"resultsCount":   len(books),
			"availableCount": countAvailable(books),
		})
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonStr)
}

func (l Library) Borrow(params martini.Params, w http.ResponseWriter) {
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		fmt.Println(err)
		trackEvent("book_borrow_failed", pendo.Properties{
			"bookId": params["id"],
			"reason": "invalid_id",
		})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book, err := model.BookImplementation{}.Borrow(id)

	if err != nil {
		fmt.Println(err)

		switch {
		case errors.Is(err, model.ErrAlreadyBorrowed):
			trackEvent("book_borrow_failed", pendo.Properties{
				"bookId": params["id"],
				"reason": "already_borrowed",
				"title":  book.Title,
				"author": book.Author,
			})
			w.WriteHeader(http.StatusBadRequest)
		case errors.Is(err, datastore.ErrNoSuchEntity):
			trackEvent("book_borrow_failed", pendo.Properties{
				"bookId": params["id"],
				"reason": "not_found",
			})
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	trackEvent("book_borrowed", pendo.Properties{
		"bookId": book.Id,
		"title":  book.Title,
		"author": book.Author,
	})
	w.WriteHeader(http.StatusNoContent)
}

func (l Library) Return(params martini.Params, w http.ResponseWriter) {
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		fmt.Println(err)
		trackEvent("book_return_failed", pendo.Properties{
			"bookId": params["id"],
			"reason": "invalid_id",
		})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book, err := model.BookImplementation{}.Return(id)

	if err != nil {
		fmt.Println(err)

		switch {
		case errors.Is(err, model.ErrNotBorrowed):
			trackEvent("book_return_failed", pendo.Properties{
				"bookId": params["id"],
				"reason": "not_borrowed",
				"title":  book.Title,
				"author": book.Author,
			})
			w.WriteHeader(http.StatusBadRequest)
		case errors.Is(err, datastore.ErrNoSuchEntity):
			trackEvent("book_return_failed", pendo.Properties{
				"bookId": params["id"],
				"reason": "not_found",
			})
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	trackEvent("book_returned", pendo.Properties{
		"bookId": book.Id,
		"title":  book.Title,
		"author": book.Author,
	})
	w.WriteHeader(http.StatusNoContent)
}

func (l Library) Add(r *http.Request, w http.ResponseWriter) {
	var book model.Book
	err := json.NewDecoder(r.Body).Decode(&book)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book, err = model.BookImplementation{}.Add(book)

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

	trackEvent("book_added", pendo.Properties{
		"bookId": book.Id,
		"title":  book.Title,
		"author": book.Author,
	})
	w.WriteHeader(http.StatusOK)
	w.Write(jsonStr)
}

// trackEvent sends a Pendo Track Event for a library request. The API has no
// authentication yet, so there is no reader or account id to attribute events
// to, and they are sent as Pendo's "system" visitor and account.
func trackEvent(event string, properties pendo.Properties) {
	pendo.Track(event, pendo.SystemID, pendo.SystemID, properties)
}

// filtersProperty describes the applied filters in query string form,
// e.g. "author=orwell&title=farm".
func filtersProperty(filter model.BookFilter) string {
	values := url.Values{}

	if filter.Title != "" {
		values.Set("title", filter.Title)
	}

	if filter.Author != "" {
		values.Set("author", filter.Author)
	}

	return values.Encode()
}

// countAvailable returns how many of the books are not borrowed.
func countAvailable(books []model.Book) int {
	available := 0

	for _, book := range books {
		if !book.Borrowed {
			available++
		}
	}

	return available
}
