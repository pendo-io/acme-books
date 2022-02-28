package server

import (
	"github.com/go-martini/martini"

	"acme-books/controllers"
)

func NewRouter(library *controllers.LibraryController) *martini.ClassicMartini {
	router := martini.Classic()

	router.Get("/books", library.ListAll)
	router.Get("/books/:id", library.GetByKey)

	return router
}
