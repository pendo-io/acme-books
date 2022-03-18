package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-martini/martini"

	"acme-books/repository"
	"acme-books/utils"
)

type LibraryController struct{

}


func (lc LibraryController) GetByKey(repository repository.BookRepository, ctx context.Context, params martini.Params, w http.ResponseWriter) {
	id, err := strconv.Atoi(params["id"])

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

func (lc LibraryController) ListAll(repo repository.BookRepository, ctx context.Context, r *http.Request, w http.ResponseWriter) {

	output, _ := repo.GetBooks(ctx)
	jsonStr, err := json.MarshalIndent(output, "", "  ")

	if err != nil {
		utils.ErrorResponse(w,http.StatusInternalServerError, err)
		return
	}

	utils.OKResponse(w, jsonStr)
}
