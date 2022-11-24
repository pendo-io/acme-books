package server

import (
	"github.com/go-martini/martini"

	"acme-books/controllers"
)

func NewRouter() *martini.ClassicMartini {
	libraryController := new(controllers.LibraryController)

	router := martini.Classic()

	router.Get("/books", libraryController.ListAll)
	router.Get("/books/:id", libraryController.GetByKey)
	router.Put("/books/:id/borrow", libraryController.Borrow)
	router.Put("/books/:id/return", libraryController.Return)
	router.Post("/books", libraryController.Add)

	return router
}
