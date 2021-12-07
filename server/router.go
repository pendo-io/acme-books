package server

import (
	"github.com/go-martini/martini"

	"acme-books/controllers"
)

func NewRouter(libraryController *controllers.LibraryController) *martini.ClassicMartini {

	router := martini.Classic()

	router.Get("/books", libraryController.ListAll)

	router.Put("/:id/borrow", libraryController.Borrow)

	router.Put("/:id/return", libraryController.Return)

	router.Post("/book", libraryController.Add)
	return router
}
