package controllers

import (
	"encoding/json"
	"net/http"

	"cloud.google.com/go/datastore"

	"acme-books/service"
)

type LibraryController struct{}

func (lc LibraryController) ListAll(r *http.Request, w http.ResponseWriter) {
	qs := r.URL.Query()
	author, title := qs.Get("author"), qs.Get("title")

	query := datastore.NewQuery("Book").
				Order("Id")

	if author != "" {
		query = query.Filter("Author=", author)
	}
	if title != "" {
		query = query.Filter("Title=", title)
	}

	output := service.GetFromStore(query)

	j, _ := json.MarshalIndent(output, "", "  ")
	w.Write(j)
}
