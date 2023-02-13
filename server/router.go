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

	return router
}
