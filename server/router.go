package server

import (
	"github.com/go-martini/martini"

	"acme-books/controllers"
)

func NewRouter(library *controllers.LibraryController) *martini.ClassicMartini {
	router := martini.Classic()

	router.Get("/books", library.List)
	router.Get("/books/:id", library.Get)
	router.Put("/:id/borrow", library.BorrowOrReturn(true))
	router.Put("/:id/return", library.BorrowOrReturn(false))

	return router
}
