package server

import (
	"acme-books/models"
	"context"
)

func Init(host string, port string) error {
	datastoreQuery, err := models.CreateConnection("acme-books", context.Background())
	if err != nil {
		return err
	}
	r := NewRouter(datastoreQuery)
	r.RunOnAddr(host + ":" + port)
	return nil
}
