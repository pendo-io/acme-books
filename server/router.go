package server

import (
	"context"

	"acme-books/controllers"
	"acme-books/models"
	"acme-books/repository"

	"cloud.google.com/go/datastore"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
)

func NewRouter(client *datastore.Client , ctx context.Context) *martini.ClassicMartini {
	libraryController := new(controllers.LibraryController)

	router := martini.Classic()
	router.Map(ctx)

	repository:= repository.NewBookRepositoryDataStore(client)
	router.Map(repository)

	router.Get("/books", libraryController.ListAll)
	router.Get("/books/:id", libraryController.GetByKey)
	router.Delete("/books/:id", libraryController.DeleteBook)
	router.Put("/:id/borrow", libraryController.Borrow)
	router.Put("/:id/return", libraryController.Return)
	router.Post("/book",binding.Bind(models.Book{}), libraryController.AddBook)

	return router
}

