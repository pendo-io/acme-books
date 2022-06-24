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
	router.Post("/book", libraryController.PostBook)

	return router
}
