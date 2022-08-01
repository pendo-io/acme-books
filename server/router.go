package server

import (
	"github.com/go-martini/martini"

	"acme-books/controllers"
)

func NewRouter() *martini.ClassicMartini {
	libraryController := controllers.NewLibraryController()

	router := martini.Classic()

	router.Get("/books", libraryController.ListAll)
	router.Get("/books/:id", libraryController.GetByKey)

	return router
}
