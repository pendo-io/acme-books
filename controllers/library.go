package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"cloud.google.com/go/datastore"
	"github.com/go-martini/martini"

	"acme-books/models"
)

type LibraryController struct{}

func internalErrorRendor(err error, w http.ResponseWriter) {
	fmt.Println(err)
	w.WriteHeader(http.StatusInternalServerError)
}

func renderData(data func() (interface{}, error), w http.ResponseWriter) {
	foo, err := data()
	render(foo, err, w)
}

func render(data interface{}, err error, w http.ResponseWriter) {
	if err != nil {
		internalErrorRendor(err, w)
		return
	}
	jsonStr, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		internalErrorRendor(err, w)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonStr)
}

func get_id(params martini.Params, w http.ResponseWriter) (int, error) {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return 0, err
	}

	return id, nil
}
func (lc LibraryController) New(req *http.Request, w http.ResponseWriter, datastoreQuery models.ModelInterface) {
	decoder := json.NewDecoder(req.Body)

	var book models.Book
	err := decoder.Decode(&book)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = datastoreQuery.Create(book)

	render(book, err, w)
}
func (lc LibraryController) Borrow(params martini.Params, w http.ResponseWriter, datastoreQuery models.ModelInterface) {
	id, err := get_id(params, w)
	if err != nil {
		return
	}
	book, err := datastoreQuery.GetById(id)
	if err != nil {
		internalErrorRendor(err, w)
		return
	}
	if book.Borrowed {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	book.Borrowed = !book.Borrowed
	if err = datastoreQuery.Change(id, book); err != nil {
		internalErrorRendor(err, w)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}
func (lc LibraryController) Return(params martini.Params, w http.ResponseWriter, datastoreQuery models.ModelInterface) {
	id, err := get_id(params, w)
	if err != nil {
		internalErrorRendor(err, w)
		return
	}
	book, err := datastoreQuery.GetById(id)
	if err != nil {
		internalErrorRendor(err, w)
		return
	}
	if !book.Borrowed {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	book.Borrowed = !book.Borrowed
	if err = datastoreQuery.Change(id, book); err != nil {
		internalErrorRendor(err, w)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func (lc LibraryController) GetByKey(params martini.Params, w http.ResponseWriter, datastoreQuery models.ModelInterface) {
	id, err := get_id(params, w)
	if err != nil {
		return
	}
	book, err := datastoreQuery.GetById(id)
	render(book, err, w)
}

func (lc LibraryController) ListAll(params martini.Params, w http.ResponseWriter, datastoreQuery models.ModelInterface) {
	query := datastore.NewQuery("Book").Order("Id")

	// Using the reflections api go through all the feilds of the struct looking up there
	// json tag and then seeing if we have any parameters that match that json tag. If so filter on
	// that value.
	// Yes this method does incur a higher run time cost, but it does mean that filtering will always
	// work for that model. The alternative would be to write some templating code but that would be
	// a mess.
	t := reflect.TypeOf(models.Book{})
	for i := 0; i < t.NumField(); i++ {
		if value, ok := t.Field(i).Tag.Lookup("json"); ok {
			if params[value] != "" {
				query = query.Filter(t.Field(i).Name+" =", params[value])
			}
		}
	}

	// inlining this does make it more ergonomic in anyay... lol
	// renderData(func() (interface{}, error) { return datastoreQuery.List(query) }, w)

	output, err := datastoreQuery.List(query)
	render(output, err, w)
}
