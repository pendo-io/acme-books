package server

import (
	"cloud.google.com/go/datastore"
	"github.com/go-martini/martini"

	"acme-books/controllers"
)

func NewRouter(client datastore.Client) *martini.ClassicMartini {
	libraryController := controllers.LibraryController{Client: client}
	router := martini.Classic()
	router.Get("/books", libraryController.ListAll)
	return router
}
