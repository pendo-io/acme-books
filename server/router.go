package server

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/go-martini/martini"

	"acme-books/controllers"
)

func NewRouter(ctx context.Context, client *datastore.Client) *martini.ClassicMartini {
	libraryController := new(controllers.LibraryController)

	router := martini.Classic()
	router.Map(ctx)
	router.Map(client)

	router.Get("/books", libraryController.ListAll)
	router.Get("/books/:id", libraryController.GetByKey)
	router.Put("/books/:id/borrow", libraryController.BorrowBook)
	router.Put("/books/:id/return", libraryController.ReturnBook)
	router.Post("/book", libraryController.Create)
	return router
}
