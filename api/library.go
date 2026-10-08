package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

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
		trackLookupFailed(params["id"], "invalid_id", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book, err := model.BookImplementation{}.GetByKey(id)

	if err != nil {
		fmt.Println(err)
		reason := "datastore_error"
		if errors.Is(err, datastore.ErrNoSuchEntity) {
			reason = "not_found"
		}
		trackLookupFailed(params["id"], reason, http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonStr, err := json.MarshalIndent(book, "", "  ")

	if err != nil {
		fmt.Println(err)
		trackLookupFailed(params["id"], "serialization_error", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pendo.Track(pendo.EventBookDetailsViewed, pendo.Properties{
		"book_id":  strconv.FormatInt(book.Id, 10),
		"title":    book.Title,
		"author":   book.Author,
		"borrowed": book.Borrowed,
	})

	w.WriteHeader(http.StatusOK)
	w.Write(jsonStr)
}

func (l Library) ListAll(r *http.Request, w http.ResponseWriter) {
	books := model.BookImplementation{}.ListAll()

	jsonStr, err := json.MarshalIndent(books, "", "  ")

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	available := 0
	for _, book := range books {
		if !book.Borrowed {
			available++
		}
	}

	// The list is not filtered yet (readme exercise 3), so filter_applied is
	// always false; filter_query records any query the client sent.
	pendo.Track(pendo.EventBookCatalogSearched, pendo.Properties{
		"filter_query":    r.URL.RawQuery,
		"filter_applied":  false,
		"results_count":   len(books),
		"available_count": available,
	})

	w.WriteHeader(http.StatusOK)
	w.Write(jsonStr)
}

// trackLookupFailed sends book_lookup_failed for a GET /books/:id request that
// did not return a book.
func trackLookupFailed(requestedID string, reason string, statusCode int) {
	pendo.Track(pendo.EventBookLookupFailed, pendo.Properties{
		"requested_id":   requestedID,
		"failure_reason": reason,
		"status_code":    statusCode,
	})
}
