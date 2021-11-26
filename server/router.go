package server

import (
	"github.com/go-martini/martini"

	"acme-books/controllers"
)

func NewRouter() *martini.ClassicMartini {
	libraryController := new(controllers.LibraryController)

	router := martini.Classic()

	router.Get("/books", libraryController.ListAll)
	router.Put("/books/:id/borrow", libraryController.BorrowBook)
	router.Put("/books/:id/return", libraryController.ReturnBook)

	return router
}
