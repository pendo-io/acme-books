package server

import (
	"context"
	"acme-books/controllers"
	"cloud.google.com/go/datastore"
)

func Init(host string, port string) {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "acme-books")
	defer client.Close()

	libraryController := new(controllers.LibraryController)
	libraryController.Init(ctx, client)

	r := NewRouter(libraryController)
	r.RunOnAddr(host + ":" + port)
}

