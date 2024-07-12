package server

import (
	"acme-books/api"
	"acme-books/model"
	"acme-books/utils"
	"github.com/go-martini/martini"
)

func NewRouter() *martini.ClassicMartini {
	library := api.NewLibrary()
	context, client := utils.CreateClient()

	bookImplementation := model.BookImplementation{
		Ctx:    context,
		Client: client,
	}

	router := martini.Classic()
	router.Map(bookImplementation)
	router.Get("/books", library.ListAll, closeContextClient)
	router.Get("/books/:id", library.GetByKey, closeContextClient)
	router.Put("/books/:id/borrow", library.Borrow, closeContextClient)
	router.Put("/books/:id/return", library.Return, closeContextClient)
	router.Post("/books", library.AddBook, closeContextClient)
	router.Delete("/books/:id", library.DeleteBook, closeContextClient)

	return router
}

func closeContextClient(bi model.BookImplementation) {
	defer func() {
		err := bi.Client.Close()
		utils.HandleGeneralError(err)
	}()
}
