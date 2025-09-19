package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/go-martini/martini"

	"acme-books/model"
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

	writeJsonResponse(w, book)
}

func (l Library) ListAll(r *http.Request, w http.ResponseWriter, bookInterface model.BookInterface) {
	books := bookInterface.ListAll()

	books, err := queryBooks(books, r)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	slices.SortFunc(books, func(a, b model.Book) int {
		return strings.Compare(strings.ToLower(a.Author), strings.ToLower(b.Author))

	})

	writeJsonResponse(w, books)
}

func (l Library) NewBook(r *http.Request, w http.ResponseWriter) {
	var book model.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book, err = model.BookImplementation{}.Put(book)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJsonResponse(w, book)

}

func (l Library) Borrow(params martini.Params, w http.ResponseWriter, bookInterface model.BookInterface) {
	changeBorrowStatus(params, w, bookInterface, true)
}

func (l Library) Return(params martini.Params, w http.ResponseWriter, bookInterface model.BookInterface) {
	changeBorrowStatus(params, w, bookInterface, false)
}

func writeJsonResponse(w http.ResponseWriter, value any) {
	jsonStr, err := json.MarshalIndent(value, "", "  ")

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonStr)

}

func changeBorrowStatus(params martini.Params, w http.ResponseWriter, bookInterface model.BookInterface, borrow bool) {
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book, err := bookInterface.GetByKey(id)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	changedBook := book
	changedBook.Borrowed = borrow

	if book == changedBook {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = bookInterface.Put(changedBook)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func queryBooks(books []model.Book, r *http.Request) ([]model.Book, error) {
	title := strings.ToLower(r.URL.Query().Get("title"))
	writer := strings.ToLower(r.URL.Query().Get("writer"))
	borrowed := strings.ToLower(r.URL.Query().Get("borrowed"))

	if title == "" && writer == "" && borrowed == "" {
		return books, nil
	}

	var queriedBooks []model.Book

	b, err := strconv.ParseBool(borrowed)

	if err != nil && borrowed != "" {
		return nil, err
	}

	for _, book := range books {
		if strings.ToLower(book.Title) == title {
			queriedBooks = append(queriedBooks, book)
		}
		if strings.ToLower(book.Author) == writer {
			queriedBooks = append(queriedBooks, book)
		}
		if borrowed != "" && book.Borrowed == b {
			queriedBooks = append(queriedBooks, book)
		}
	}
	return queriedBooks, nil
}
