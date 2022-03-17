package server

import (
	"github.com/go-martini/martini"

	"acme-books/controllers"
	"acme-books/models"
)

func NewRouter() *martini.ClassicMartini {
	libraryController := new(controllers.LibraryController)

	router := martini.Classic()

	router.Get("/books", initBookInterface, libraryController.ListAll, closeClientConnection)
	router.Get("/books/:id", initBookInterface, libraryController.GetByKey, closeClientConnection)
	router.Put("/books/:id/borrow", initBookInterface, libraryController.Borrow, closeClientConnection)
	router.Put("/books/:id/return", initBookInterface, libraryController.Return, closeClientConnection)
	router.Post("/book", initBookInterface, libraryController.Create, closeClientConnection)

	return router
}

func initBookInterface(c martini.Context) {
	bookInt := models.NewBookImplementation()
	c.Map(bookInt)
}

func closeClientConnection(c martini.Context, bookInt models.BookInterface) {
	defer bookInt.CloseClientConnection()
}
