package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"

	"acme-books/models"
)

type LibraryController struct{}

func (lc LibraryController) ListAll(r *http.Request, w http.ResponseWriter) {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "acme-books")

	defer client.Close()

	output := make([]models.Book, 0)

	it := client.Run(ctx, datastore.NewQuery("Book"))
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
