package utils

import (
	"acme-books/repository"
	"fmt"
	"net/http"
	"strconv"
)

func OKResponse(w http.ResponseWriter, jsonStr []byte){
	w.WriteHeader(http.StatusOK)
	w.Write(jsonStr)
}

func ErrorResponse(w http.ResponseWriter, err error){
	if err !=nil{
		fmt.Println(err)
		switch err.(type){
		case *repository.BorrowedError, *strconv.NumError, *repository.ReturnedError:
			w.WriteHeader(http.StatusBadRequest)

		default:
			w.WriteHeader(http.StatusInternalServerError)
		}

	}

}