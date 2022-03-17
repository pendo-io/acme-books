package utils

import (
	"fmt"
	"net/http"
)

func ErrorResponse(w http.ResponseWriter,status int, err error){
	fmt.Println(err)
	w.WriteHeader(status)
}

func OKResponse(w http.ResponseWriter, jsonStr []byte){
	w.WriteHeader(http.StatusOK)
	w.Write(jsonStr)
}
