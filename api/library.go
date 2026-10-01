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

func (l Library) GetByKey(params martini.Params, w http.ResponseWriter) {
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		fmt.Println(err)

		// Track invalid book ID request
		pendo.Track("Book Retrieval Failed", "system", "system", map[string]interface{}{
			"raw_id":      params["id"],
			"error_type":  "invalid_id",
			"status_code": http.StatusBadRequest,
		})

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book, err := model.BookImplementation{}.GetByKey(id)

	if err != nil {
		fmt.Println(err)

		// Track datastore lookup failure
		pendo.Track("Book Retrieval Failed", "system", "system", map[string]interface{}{
			"book_id":     id,
			"error_type":  "datastore_error",
			"status_code": http.StatusInternalServerError,
		})

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonStr, err := json.MarshalIndent(book, "", "  ")

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Track successful book retrieval
	pendo.Track("Book Retrieved", "system", "system", map[string]interface{}{
		"book_id":     book.Id,
		"book_title":  book.Title,
		"book_author": book.Author,
		"borrowed":    book.Borrowed,
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

	// Track book list retrieval
	pendo.Track("Books Listed", "system", "system", map[string]interface{}{
		"book_count": len(books),
	})

	w.WriteHeader(http.StatusOK)
	w.Write(jsonStr)
}
