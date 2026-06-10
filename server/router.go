package server

import (
	"github.com/go-martini/martini"

	"acme-books/api"
)

func NewRouter() *martini.ClassicMartini {
	library := new(api.Library)

	router := martini.Classic()

	router.Get("/books", library.ListAll)
	router.Get("/books/:id", library.GetByKey)
	router.Put("/books/:id/borrow", library.BorrowBook)
	router.Put("/books/:id/return", library.ReturnBook)
	router.Post("/books", library.CreateBook)

	return router
}
