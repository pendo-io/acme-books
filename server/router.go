package server

import (
	"acme-books/controllers"
	"acme-books/models"
	"log"

	"github.com/go-martini/martini"
)

func NewRouter(datastoreQuery models.DatastoreQuery) *martini.ClassicMartini {
	libraryController := new(controllers.LibraryController)

	router := martini.Classic()
	router.Map(datastoreQuery)
	router.Get("/books", libraryController.ListAll, closeClientConnection)
	router.Put("/books", libraryController.New, closeClientConnection)
	router.Get("/books/:id", libraryController.GetByKey, closeClientConnection)
	router.Put("/books/:id/borrow", libraryController.Borrow, closeClientConnection)
	router.Put("/books/:id/return", libraryController.Return, closeClientConnection)

	return router
}

// uses defer as it wraps the request
func closeClientConnection(c martini.Context, datastoreQuery models.DatastoreQuery, log *log.Logger) {
	defer func() {
		err := datastoreQuery.Close()
		if err != nil {
			log.Println(err)
		}
	}()
}
