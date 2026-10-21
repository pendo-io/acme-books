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

// Pendo Track Event names. They must match the track types registered in
// Pendo exactly.
const (
	// eventBookDetailsRetrieved is sent when GET /books/:id returns a book.
	eventBookDetailsRetrieved = "book_details_retrieved"
	// eventBookLookupFailed is sent when GET /books/:id does not return a book.
	eventBookLookupFailed = "book_lookup_failed"
)

type Library struct{}

func (l Library) GetByKey(params martini.Params, r *http.Request, w http.ResponseWriter) {
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		trackBookLookupFailed(r, params["id"], "invalid_id", http.StatusBadRequest)
		return
	}

	book, err := model.BookImplementation{}.GetByKey(id)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		failureReason := "datastore_error"
		if errors.Is(err, datastore.ErrNoSuchEntity) {
			failureReason = "not_found"
		}
		trackBookLookupFailed(r, params["id"], failureReason, http.StatusInternalServerError)
		return
	}

	jsonStr, err := json.MarshalIndent(book, "", "  ")

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		trackBookLookupFailed(r, params["id"], "serialization_error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonStr)

	pendo.Track(r, eventBookDetailsRetrieved, pendo.Properties{
		"book_id":  book.Id,
		"title":    book.Title,
		"author":   book.Author,
		"borrowed": book.Borrowed,
	})
}

func (l Library) ListAll(r *http.Request, w http.ResponseWriter) {
	books := model.BookImplementation{}.ListAll()

	jsonStr, err := json.MarshalIndent(books, "", "  ")

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonStr)
}

// trackBookLookupFailed records a GET /books/:id request that did not return a book.
func trackBookLookupFailed(r *http.Request, requestedID string, failureReason string, statusCode int) {
	pendo.Track(r, eventBookLookupFailed, pendo.Properties{
		"requested_id":   requestedID,
		"failure_reason": failureReason,
		"status_code":    statusCode,
	})
}
