package controllers

import (
	"encoding/json"
	"net/http"

	"cloud.google.com/go/datastore"

	"acme-books/service"
)

type LibraryController struct{}

func (lc LibraryController) ListAll(r *http.Request, w http.ResponseWriter) {
	query := datastore.NewQuery("Book")

	output := service.GetFromStore(query)

	j, _ := json.MarshalIndent(output, "", "  ")
	w.Write(j)
}
