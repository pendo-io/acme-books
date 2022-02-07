package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"github.com/go-martini/martini"
	"acme-books/models"
)

type LibraryController struct{}

func (lc LibraryController) GetByKey(params martini.Params, w http.ResponseWriter) {
	if book, err := bookFromId(params["id"], w); err != nil || book == nil {
		return
	} else {
		if json, err := json.MarshalIndent(book, "", "  "); err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write(json)
		}
	}
}

func (lc LibraryController) ListAll(params martini.Params, request *http.Request, writer http.ResponseWriter) {

	if filter, err := getFilter(request); err != nil {
		fmt.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		if output, err := models.GetBooks(filter); err != nil {
			fmt.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			if json, err := json.MarshalIndent(output, "", "  "); err != nil {
				fmt.Println(err)
				writer.WriteHeader(http.StatusInternalServerError)
				return
			} else {
				writer.WriteHeader(http.StatusOK)
				writer.Write(json)
			}
		}
	}
}

func (lc LibraryController) CreateBook(request *http.Request, writer http.ResponseWriter) {
	var book = models.Book{}
	if err := json.NewDecoder(request.Body).Decode(&book); err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := models.AddBook(&book); err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		if json, err := json.MarshalIndent(book, "", "  "); err != nil {
			fmt.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return;
		} else {
			writer.WriteHeader(http.StatusOK)
			writer.Write(json)
		}
	}


}

func (lc LibraryController) BorrowBook(params martini.Params, w http.ResponseWriter) {
	if book, err := bookFromId(params["id"], w); err != nil || book == nil {
		return
	} else {
		if book.Borrowed { // this test should be done in models.BorrowBook
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := book.Borrow(); err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func (lc LibraryController) ReturnBook(params martini.Params, w http.ResponseWriter) {
	if book, err := bookFromId(params["id"], w); err != nil || book == nil {
		return
	} else {
		if !book.Borrowed { // this test should be done in models.ReturnBook
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := book.Return(); err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func bookFromId(bookId string, w http.ResponseWriter) (*models.Book, error) {
	id, err := strconv.ParseInt(bookId, 10, 64)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return nil, err
	}
	var book *models.Book
	book, err = models.GetBook(int64(id))

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return nil, err
	} else if book == nil {
		w.WriteHeader(http.StatusNotFound)
		return nil, nil
	}
	return book, nil
}


// http://localhost:3030/books?field=Title&op=%3d&value=Animal%20Farm
// http://localhost:3030/books?field=Id&op=%3C&value=3
func getFilter(request *http.Request) (*models.Filter, error) {
	var filter *models.Filter
	fieldFilter := getQueryParam(request, "field")
	if fieldFilter != "" {
		var value interface{}
		switch fieldFilter {
		case "Id":
			var err error
			if value, err = strconv.Atoi(getQueryParam(request, "value")); err != nil {
				return nil, err
			}
		default:
			value = getQueryParam(request, "value")
		}
		filter = &models.Filter{
			Operation: fmt.Sprintf("%s %s", fieldFilter, getQueryParam(request, "op")),
			Value : value,
		}
	}
	return filter, nil
}

func getQueryParam(request *http.Request, field string) string {
	return request.URL.Query().Get(field)
}
