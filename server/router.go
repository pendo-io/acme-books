package server

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/go-martini/martini"

	"acme-books/controllers"
	"acme-books/repository"
)

func NewRouter(client *datastore.Client , ctx context.Context) *martini.ClassicMartini {
	libraryController := new(controllers.LibraryController)

	router := martini.Classic()
	router.Map(ctx)
	repository:= repository.NewBookRepositoryDataStore(client)
	router.Map(repository)


	router.Get("/books", libraryController.ListAll)
	router.Get("/books/:id", libraryController.GetByKey)
	router.Put("/:id/borrow", libraryController.Borrow)
	router.Put("/:id/return", libraryController.Return)

	return router
}

