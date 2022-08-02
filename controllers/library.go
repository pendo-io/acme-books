package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"acme-books/models"
	"cloud.google.com/go/datastore"
	"github.com/go-martini/martini"
)

type LibraryController struct {
	ctx    context.Context
	client *datastore.Client
}

func NewLibraryController() *LibraryController {
	lc := new(LibraryController)
	lc.ctx = context.Background()
	lc.client, _ = datastore.NewClient(lc.ctx, "acme-books")

	books := []models.Book{
		{Id: 1, Author: "George Orwell", Title: "1984", Borrowed: false},
		{Id: 2, Author: "George Orwell", Title: "Animal Farm", Borrowed: false},
		{Id: 3, Author: "Robert Jordan", Title: "Eye of the world", Borrowed: false},
		{Id: 4, Author: "Various", Title: "Collins Dictionary", Borrowed: false},
	}

	var keys []*datastore.Key

	for _, book := range books {
		keys = append(keys, datastore.IDKey("Book", book.Id, nil))
	}

	if _, err := lc.client.PutMulti(lc.ctx, keys, books); err != nil {
		fmt.Println(err)
	}

	return lc
}
func (lc LibraryController) Close() error {
	return lc.client.Close()
}

func (lc LibraryController) GetByKey(params martini.Params, w http.ResponseWriter) {

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	key := datastore.IDKey("Book", int64(id), nil)
	var book models.Book
	err = lc.client.Get(lc.ctx, key, &book)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	writeJson(book, w)
}

func (lc LibraryController) ListAll(r *http.Request, w http.ResponseWriter) {
	filter := r.URL.Query().Get("q")

	var output []models.Book
	query := datastore.NewQuery("Book")

	if filter != "" {
		splitIndices := regexp.MustCompile("(<=|<|>=|>|=)").FindStringIndex(filter)
		if len(splitIndices) < 2 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		sepIndex := splitIndices[1]
		fieldName := filter[:splitIndices[0]]
		filterStr := filter[:sepIndex]

		metaQuery := datastore.NewQuery("__property__")
		type Prop struct {
			Repr []string `datastore:"property_representation"`
		}
		var props []Prop

		keys, err := lc.client.GetAll(lc.ctx, metaQuery, &props)
		//fmt.Println(props)
		//fmt.Println(keys)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var columnKind string
		ok := false
		for i, k := range keys {
			if k.Name == fieldName {
				columnKind = props[i].Repr[0]
				ok = true
				break
			}
		}
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var filterVal interface{}
		switch columnKind {
		case "BOOLEAN":
			filterVal, err = strconv.ParseBool(filter[sepIndex:])
		case "INT64":
			filterVal, err = strconv.Atoi(filter[sepIndex:])
		default:
			filterVal = filter[sepIndex:]
		}
		//fmt.Println(fieldName)
		//fmt.Println(filterStr)
		//fmt.Println(filterVal)

		query = query.Filter(filterStr, filterVal)

	}

	query = query.Order("Id")
	it := lc.client.Run(lc.ctx, query)

	for {
		var b models.Book
		_, err := it.Next(&b)
		if err != nil {
			fmt.Println(err)
			break
		}
		output = append(output, b)
	}

	writeJson(output, w)
}

func (lc LibraryController) Borrow(params martini.Params, w http.ResponseWriter) {
	lc.SetBorrowed(params, w, true)
	return
}

func (lc LibraryController) Return(params martini.Params, w http.ResponseWriter) {
	lc.SetBorrowed(params, w, false)
	return
}
func (lc LibraryController) SetBorrowed(params martini.Params, w http.ResponseWriter, borrowed bool) {
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	key := datastore.IDKey("Book", int64(id), nil)

	var book models.Book
	err = lc.client.Get(lc.ctx, key, &book)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if book.Borrowed == borrowed {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	book.Borrowed = borrowed
	_, err = lc.client.Mutate(lc.ctx, datastore.NewUpdate(key, &book))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	return
}

func writeJson(item interface{}, w http.ResponseWriter) {
	jsonStr, err := json.MarshalIndent(item, "", "  ")

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonStr)
}
