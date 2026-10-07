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

func (l Library) GetByKey(params martini.Params, r *http.Request, w http.ResponseWriter) {
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		fmt.Println(err)
		trackLookupFailed(r, params["id"], "invalid_id", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book, err := model.BookImplementation{}.GetByKey(id)

	if err != nil {
		fmt.Println(err)
		// The response is a 500 either way; the reason separates missing books
		// (stale links, deleted books) from datastore failures.
		reason := "server_error"
		if errors.Is(err, datastore.ErrNoSuchEntity) {
			reason = "not_found"
		}
		trackLookupFailed(r, params["id"], reason, http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonStr, err := json.MarshalIndent(book, "", "  ")

	if err != nil {
		fmt.Println(err)
		trackLookupFailed(r, params["id"], "server_error", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pendo.Track(r, "book_details_viewed", pendo.Properties{
		"bookId":   id,
		"title":    book.Title,
		"author":   book.Author,
		"borrowed": book.Borrowed,
	})

	w.WriteHeader(http.StatusOK)
	w.Write(jsonStr)
}

// trackLookupFailed reports a GET /books/:id request that did not return a book.
func trackLookupFailed(r *http.Request, requestedID string, reason string, statusCode int) {
	pendo.Track(r, "book_lookup_failed", pendo.Properties{
		"requestedId":   requestedID,
		"failureReason": reason,
		"statusCode":    statusCode,
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

	borrowedCount := 0
	for _, book := range books {
		if book.Borrowed {
			borrowedCount++
		}
	}

	// The list isn't filtered yet (readme exercise 3), so filters records what
	// callers asked for rather than what was applied.
	pendo.Track(r, "book_catalog_searched", pendo.Properties{
		"filters":        r.URL.RawQuery,
		"hasFilters":     r.URL.RawQuery != "",
		"resultsCount":   len(books),
		"availableCount": len(books) - borrowedCount,
		"borrowedCount":  borrowedCount,
	})

	w.WriteHeader(http.StatusOK)
	w.Write(jsonStr)
}

// TODO(pendo): readme exercises 4 and 5 add PUT /:id/borrow, PUT /:id/return
// and POST /book. Pendo already has track types registered for them, so when
// those handlers are added here, send these events with pendo.Track, using the
// names exactly as written:
//   - book_borrowed, once Borrowed=true is saved (204): bookId, title, author
//   - book_borrow_failed, on each failure response: requestedId, failureReason
//     ("invalid_id", "not_found", "already_borrowed"), statusCode, plus title
//     and author whenever the book was found
//   - book_returned, once Borrowed=false is saved (204): bookId, title, author
//   - book_return_failed, on each failure response: requestedId, failureReason
//     ("invalid_id", "not_found", "not_borrowed"), statusCode
//   - book_added, once the Put succeeds (200): bookId (the generated key),
//     title, author. Not for the demo books that bootstrapBooks in main.go
//     seeds on every start.
