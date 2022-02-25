package models

import (
	"cloud.google.com/go/datastore"
	"context"
	"net/url"
)

type Book struct {
	Id       int64
	Title    string `json:"title"`
	Author   string `json:"writer"`
	Borrowed bool   `json:"borrowed"`
}

func ListBooks(client *datastore.Client, ctx context.Context, orderBy string, params url.Values) (result []Book, err error) {
	query := datastore.NewQuery("Book").Order(orderBy)

	if params.Get("author") != "" {
		query = query.Filter("Author=", params.Get("author"))
	}
	if params.Get("title") != "" {
		query = query.Filter("Title=", params.Get("title"))
	}
	if params.Get("borrowed") == "true" {
		query = query.Filter("Borrowed=", params.Get("borrowed"))
	}

	_, err = client.GetAll(ctx, query, &result)

	return
}