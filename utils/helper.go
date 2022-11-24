package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteJsonResp(w http.ResponseWriter, data interface{}) {
	if jsonStr, err := json.MarshalIndent(data, "", "  "); err != nil {
		HandleError(err, w,  http.StatusInternalServerError)
	} else {
		w.Write(jsonStr)
	}
	return
}

func HandleError(err error, w http.ResponseWriter, errorCode int) {
	fmt.Println(err)
	w.WriteHeader(errorCode)
	return
}
