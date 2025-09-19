package server

import (
	"acme-books/model"

	"github.com/go-martini/martini"

	"acme-books/api"
)

func NewRouter() *martini.ClassicMartini {
	library := new(api.Library)

	router := martini.Classic()

	bookImplementation := model.BookImplementation{}

	router.Map(bookImplementation)
	router.Get("/books", library.ListAll)
	router.Get("/books/:id", library.GetByKey)
	router.Put("/:id/borrow", library.Borrow)
	router.Put("/:id/return", library.Return)
	router.Post("/book", library.NewBook)

	return router

}
