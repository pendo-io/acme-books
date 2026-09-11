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

func (l Library) GetByKey(params martini.Params, r *http.Request, w http.ResponseWriter) {
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

	pendo.Track("Book Viewed", "system", "system", map[string]interface{}{
		"book_id":     id,
		"book_title":  book.Title,
		"book_author": book.Author,
		"borrowed":    book.Borrowed,
	}, map[string]string{
		"userAgent": r.UserAgent(),
		"ip":        r.RemoteAddr,
		"url":       r.URL.String(),
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

	pendo.Track("Books Listed", "system", "system", map[string]interface{}{
		"book_count": len(books),
	}, map[string]string{
		"userAgent": r.UserAgent(),
		"ip":        r.RemoteAddr,
		"url":       r.URL.String(),
	})
}
