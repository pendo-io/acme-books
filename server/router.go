package server

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/go-martini/martini"

	"acme-books/controllers"
)

func NewRouter(client *datastore.Client , ctx context.Context) *martini.ClassicMartini {
	libraryController := new(controllers.LibraryController)

	router := martini.Classic()
	router.Map(client)
	router.Map(ctx)

	router.Get("/books", libraryController.ListAll)
	router.Get("/books/:id", libraryController.GetByKey)

	return router
}

