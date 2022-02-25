package server

import (
	"acme-books/controllers"
	"acme-books/models"
	"cloud.google.com/go/datastore"
	"context"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
)

func NewRouter(client *datastore.Client, ctx context.Context) *martini.ClassicMartini {
	libraryController := new(controllers.LibraryController)

	router := martini.Classic()
	router.Map(client)
	router.Map(ctx)

	router.Get("/books", libraryController.ListAll)
	router.Get("/books/:id", libraryController.GetByKey)
	router.Put("/books/:id/borrow", libraryController.BorrowByKey)
	router.Put("/books/:id/return", libraryController.ReturnByKey)
	router.Post("/books", binding.Bind(models.Book{}), libraryController.CreateBook)

	return router
}
