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
	router.Get("/searchByName/:name", libraryController.SearchByName)
	router.Put("/:id/borrow", libraryController.Borrow)
	router.Put("/:id/return", libraryController.Return)
	router.Post("/book", libraryController.CreateBook)

	return router
}
