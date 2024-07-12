package utils

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
)

func CreateClient() (context.Context, *datastore.Client) {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, os.Getenv(DATASTORE_PROJECT_ID))

	return ctx, client
}

func HandleHttpError(err error, w http.ResponseWriter, resp int) bool {
	if err != nil {
		w.WriteHeader(resp)
		w.Write([]byte(fmt.Sprintf("ERROR %d: %s", resp, err)))
		return true
	}
	return false
}

func HandleGeneralError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func HandleFatalError(err error, msg string) {
	if err != nil {
		log.Fatal(msg)
	}
}

func CreateHttpResponse(w http.ResponseWriter, jsonStr []byte) {
	w.WriteHeader(http.StatusOK)
	w.Write(jsonStr)
}
