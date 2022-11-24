package controllers

import (
	"errors"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"cloud.google.com/go/datastore"
	"github.com/go-martini/martini"

	"acme-books/models"
	"acme-books/utils"
)

type LibraryController struct{}

func (lc LibraryController) GetByKey(params martini.Params, w http.ResponseWriter, bookInt models.BookInterface) {
	id, err := strconv.Atoi(params["id"])

	if err != nil { 
		utils.HandleError(err, w, http.StatusBadRequest)
		return
	}

	book, err := bookInt.GetByKey(id)

	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			utils.HandleError(err, w, http.StatusNotFound)
		} else {
			utils.HandleError(err, w, http.StatusInternalServerError)
		}
		return
	}

	utils.WriteJsonResp(w, book)
}

func (lc LibraryController) ListAll(r *http.Request, w http.ResponseWriter, bookInt models.BookInterface) {
	title := r.URL.Query().Get("title"); 

	books, err := bookInt.Query(title)
	
	if err != nil {
		utils.HandleError(err, w, http.StatusInternalServerError)
		return
	}

	utils.WriteJsonResp(w, books)
}

func (lc LibraryController) Borrow(params martini.Params, w http.ResponseWriter, bookInt models.BookInterface) {
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		utils.HandleError(err, w, http.StatusBadRequest)
		return
	}

	getBookAndUpdateStatus(id, "borrow", w, bookInt)
}

func (lc LibraryController) Return(params martini.Params, w http.ResponseWriter, bookInt models.BookInterface) {
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		utils.HandleError(err, w, http.StatusBadRequest)
		return
	}

	getBookAndUpdateStatus(id, "return", w, bookInt)
}

func getBookAndUpdateStatus(id int, status string, w http.ResponseWriter, bookInt models.BookInterface) {
	book, err := bookInt.GetByKey(id)

	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			utils.HandleError(err, w, http.StatusBadRequest)
		} else {
			utils.HandleError(err, w, http.StatusInternalServerError)
		}
		return
	}

	if status == "borrow" && book.Borrowed {
		err := errors.New("Already Borrowed")
		utils.HandleError(err, w, http.StatusBadRequest)
		return
	} else if status == "return" && !book.Borrowed {
		err := errors.New("Not Borrowed") 
		utils.HandleError(err, w, http.StatusBadRequest)
		return
	}

	book.Borrowed = !book.Borrowed

	_, error := bookInt.Put(book, false)

	if error != nil {
		utils.HandleError(err, w, http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
}

func (lc LibraryController) Create(r *http.Request, w http.ResponseWriter, bookInt models.BookInterface) {
	if body, err := ioutil.ReadAll(r.Body); err != nil { 
		utils.HandleError(err, w, http.StatusInternalServerError)
	} else {
		var book models.Book
		
		if err = json.Unmarshal(body, &book); err != nil { 
			utils.HandleError(err, w, http.StatusInternalServerError)
		} else {
			id, err := bookInt.Put(book, true)
			
			if err != nil {
				utils.HandleError(err, w, http.StatusInternalServerError)
			}
			book.Id = id
			utils.WriteJsonResp(w, book)
		}
	}
}
