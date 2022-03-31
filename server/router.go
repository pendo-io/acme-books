package server

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"

	"acme-books/controllers"
	"acme-books/models"
)

func NewRouter(library *controllers.LibraryController) *martini.ClassicMartini {
	router := martini.Classic()

	router.Get("/books", library.List)
	router.Get("/books/:id", library.Get)
	router.Put("/:id/borrow", library.BorrowOrReturn(true))
	router.Put("/:id/return", library.BorrowOrReturn(false))
	router.Post("/book", binding.Bind(models.Book{}), library.AddBook)

	return router
}
