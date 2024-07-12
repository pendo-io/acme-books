package api

import (
	"acme-books/utils"
	"cloud.google.com/go/datastore"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strconv"

	"github.com/go-martini/martini"

	"acme-books/model"
)

type Library struct {
}

func NewLibrary() *Library {
	library := new(Library)

	return library
}

func (l Library) GetByKey(params martini.Params, w http.ResponseWriter, bookInterface model.BookInterface) {
	id, err := strconv.Atoi(params["id"])

	if utils.HandleHttpError(err, w, http.StatusBadRequest) {
		return
	}

	book, err := bookInterface.GetByKey(id)

	if utils.HandleHttpError(err, w, http.StatusNotFound) {
		return
	}

	jsonStr, err := json.MarshalIndent(book, "", "  ")

	if utils.HandleHttpError(err, w, http.StatusInternalServerError) {
		return
	}

	utils.CreateHttpResponse(w, jsonStr)
}

func (l Library) ListAll(r *http.Request, w http.ResponseWriter, bookInterface model.BookInterface) {
	query := datastore.NewQuery("Book")

	t := reflect.TypeOf(model.Book{})

	for i := 0; i < t.NumField(); i++ {
		if v, ok := t.Field(i).Tag.Lookup("json"); ok {
			if filter := r.URL.Query().Get(v); filter != "" {
				switch t.Field(i).Type.Kind() {
				case reflect.Bool:
					boolValue, err := strconv.ParseBool(filter)
					if utils.HandleHttpError(err, w, http.StatusBadRequest) {
						return
					}
					query = query.Filter(t.Field(i).Name+"=", boolValue)
				default:
					query = query.Filter(t.Field(i).Name+"=", filter)
				}
			}
		}
	}

	books := bookInterface.ListAll(query)

	jsonStr, err := json.MarshalIndent(books, "", "  ")

	if utils.HandleHttpError(err, w, http.StatusInternalServerError) {
		return
	}

	utils.CreateHttpResponse(w, jsonStr)
}

func (l Library) Borrow(params martini.Params, w http.ResponseWriter, bookInterface model.BookInterface) {
	id, err := strconv.Atoi(params["id"])

	if utils.HandleHttpError(err, w, http.StatusInternalServerError) {
		return
	}

	err = bookInterface.ChangeBorrowedStatus(id, true)

	if utils.HandleHttpError(err, w, http.StatusInternalServerError) {
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (l Library) Return(params martini.Params, w http.ResponseWriter, bookInterface model.BookInterface) {
	id, err := strconv.Atoi(params["id"])

	if utils.HandleHttpError(err, w, http.StatusBadRequest) {
		return
	}

	err = bookInterface.ChangeBorrowedStatus(id, false)

	if utils.HandleHttpError(err, w, http.StatusInternalServerError) {
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (l Library) AddBook(r *http.Request, w http.ResponseWriter, bookInterface model.BookInterface) {
	var book model.Book

	err := json.NewDecoder(r.Body).Decode(&book)

	if utils.HandleHttpError(err, w, http.StatusInternalServerError) {
		return
	}

	if book.Title == "" {
		utils.HandleHttpError(errors.New("The title field is empty"), w, http.StatusBadRequest)
		return
	}

	if book.Author == "" {
		utils.HandleHttpError(errors.New("The author field is empty"), w, http.StatusBadRequest)
		return
	}

	book, err = bookInterface.AddBook(book)

	if utils.HandleHttpError(err, w, http.StatusInternalServerError) {
		return
	}

	jsonStr, err := json.MarshalIndent(book, "", "  ")

	if utils.HandleHttpError(err, w, http.StatusInternalServerError) {
		return
	}

	utils.CreateHttpResponse(w, jsonStr)
}

func (l Library) DeleteBook(params martini.Params, w http.ResponseWriter, bookInterface model.BookInterface) {
	id, err := strconv.Atoi(params["id"])

	if utils.HandleHttpError(err, w, http.StatusBadRequest) {
		return
	}

	err = bookInterface.Delete(id)

	if utils.HandleHttpError(err, w, http.StatusInternalServerError) {
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
