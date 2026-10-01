package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-martini/martini"

	"acme-books/model"
	"acme-books/pendo"
)

type Library struct{}

func (l Library) GetByKey(r *http.Request, params martini.Params, w http.ResponseWriter) {
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		fmt.Println(err)

		// Track invalid book ID request
		pendo.Track("Book Lookup Failed", "anonymous", "system", map[string]interface{}{
			"raw_id":     params["id"],
			"error_type": "invalid_id",
		}, pendo.ContextFromRequest(r))

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book, err := model.BookImplementation{}.GetByKey(id)

	if err != nil {
		fmt.Println(err)

		// Track book not found / datastore error
		pendo.Track("Book Lookup Failed", "anonymous", "system", map[string]interface{}{
			"book_id":    id,
			"error_type": "not_found",
		}, pendo.ContextFromRequest(r))

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonStr, err := json.MarshalIndent(book, "", "  ")

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Track successful book detail view
	pendo.Track("Book Details Viewed", "anonymous", "system", map[string]interface{}{
		"book_id":     book.Id,
		"book_title":  book.Title,
		"book_author": book.Author,
		"is_borrowed": book.Borrowed,
	}, pendo.ContextFromRequest(r))

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

	// Track book list view
	pendo.Track("Book List Viewed", "anonymous", "system", map[string]interface{}{
		"book_count": len(books),
	}, pendo.ContextFromRequest(r))

	w.WriteHeader(http.StatusOK)
	w.Write(jsonStr)
}
