package controllers

import (
	"context"
	"encoding/json"
	"fmt"
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
	fmt.Println(params)

	if err != nil {
		utils.ErrorResponse(w,http.StatusBadRequest, err)
		return
	}

	book , err := repository.GetBook(ctx, int(id))
	if err != nil {
		utils.ErrorResponse(w,http.StatusInternalServerError, err)
		return
	}

	jsonStr, err := json.MarshalIndent(book, "", "  ")

	if err != nil {
		utils.ErrorResponse(w,http.StatusInternalServerError, err)
		return
	}

	utils.OKResponse(w,jsonStr)
}

func (lc LibraryController) ListAll(repo repository.BookRepository, ctx context.Context, r *http.Request,
	w http.ResponseWriter) {
	qs := r.URL.Query()
	var books []models.Book

	if title:= qs.Get("title"); title!="" {
		fmt.Println(title)
		books = repo.GetBooksByTitle(ctx,title)
	}else{
		books = repo.GetBooks(ctx)
	}

	sort.Slice(books, func(i, j int) bool {return books[i].Id < books[j].Id})

	jsonStr, err := json.MarshalIndent(books, "", "  ")

	if err != nil {
		utils.ErrorResponse(w,http.StatusInternalServerError, err)
		return
	}

	utils.OKResponse(w, jsonStr)
}
