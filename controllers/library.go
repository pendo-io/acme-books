package controllers

import (
	"encoding/json"
	"net/http"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"

	"acme-books/models"
)

type LibraryController struct {
	Client datastore.Client
}

func (lc LibraryController) ListAll(r *http.Request, w http.ResponseWriter) {

	output := make([]models.Book, 0)
	it := lc.Client.Run(r.Context(), datastore.NewQuery("Book"))
	for {
		var b models.Book
		_, err := it.Next(&b)
		if err == iterator.Done {
			break
		}
		output = append(output, b)
	}

	j, _ := json.MarshalIndent(output, "", "  ")
	w.Write(j)
}
