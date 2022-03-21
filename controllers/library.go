package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"strconv"

	"github.com/go-martini/martini"

	"acme-books/models"
	"acme-books/repository"
	"acme-books/utils"
)

type LibraryController struct{

}


func (lc LibraryController) GetByKey(repository repository.BookRepository, ctx context.Context, params martini.Params, w http.ResponseWriter) {
	id, err := strconv.Atoi(params["id"])
	if err !=nil{
		utils.ErrorResponse(w,err)
		return
	}

	book , err := repository.GetBook(ctx, int(id))
	if err !=nil{
		utils.ErrorResponse(w,err)
		return
   }

	jsonStr, err := json.MarshalIndent(book, "", "  ")
	if err !=nil{
		utils.ErrorResponse(w,err)
		return
   }

   utils.OKResponse(w,jsonStr)
}

func (lc LibraryController) ListAll(repo repository.BookRepository, ctx context.Context, r *http.Request,
	w http.ResponseWriter) {
	qs := r.URL.Query()
	var books []models.Book
	var filters map[string]string = make(map[string]string)

	if qs.Has("title"){
		filters["Title"]=qs.Get("title")
	}
	if qs.Has("author"){
		filters["Author"]=qs.Get("author")
	}

	books = repo.GetBooks(ctx,filters)

	sort.Slice(books, func(i, j int) bool {return books[i].Id < books[j].Id})

	jsonStr, err := json.MarshalIndent(books, "", "  ")

	if err !=nil{
		utils.ErrorResponse(w,err)
		return
   }

	utils.OKResponse(w, jsonStr)
}

func (lc LibraryController) Borrow(repo repository.BookRepository, ctx context.Context, r *http.Request,
	w http.ResponseWriter, params martini.Params) {

	id, err := strconv.Atoi(params["id"])

	if err !=nil{
		utils.ErrorResponse(w,err)
		return
   }
	err = repo.Lending(ctx, id, true)
	if err !=nil{
		 utils.ErrorResponse(w,err)
		 return
	}
	w.WriteHeader(http.StatusNoContent)
}


func (lc LibraryController) Return(repo repository.BookRepository, ctx context.Context, r *http.Request,
		w http.ResponseWriter, params martini.Params) {
	id, err := strconv.Atoi(params["id"])
	if err !=nil{
		utils.ErrorResponse(w,err)
		return
   }
	err = repo.Lending(ctx, id, false)
	if err !=nil{
		utils.ErrorResponse(w,err)
		return
   }
	w.WriteHeader(http.StatusNoContent)
}
