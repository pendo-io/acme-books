package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"cloud.google.com/go/datastore"
	"github.com/go-martini/martini"
	"google.golang.org/api/iterator"

	"acme-books/models"
)

type LibraryController struct {
	Client datastore.Client
}

func (lc LibraryController) Borrow(r *http.Request, w http.ResponseWriter, params martini.Params) {

	// 204 if ok
	// 400 if already borrowed
	// 400 if not borrowed when returned
	paramId := params["id"]
	id, err := strconv.Atoi(paramId)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("invalid id:%v", params["id"])))
		w.WriteHeader(400)
		return
	}
	b := &models.Book{}
	if err := lc.Client.Get(r.Context(), datastore.IDKey("Book", int64(id), nil), b); err != nil {
		if err == datastore.ErrNoSuchEntity {
			w.WriteHeader(404)
			return
		}
	}
	if b.Borrowed {
		w.WriteHeader(400)
		return
	}
	log.Println("loaded book:", b)
	b.Borrowed = true
	_, err = lc.Client.Put(r.Context(), datastore.IDKey("Book", int64(b.Id), nil), b)
	if err != nil {
		log.Fatalln(err)
	}
	w.WriteHeader(204)
}

func (lc LibraryController) ListAll(r *http.Request, w http.ResponseWriter) {

	output := make([]models.Book, 0)
	query := datastore.NewQuery("Book").Order("Id")

	if af := r.URL.Query().Get("writer"); af != "" {
		query = query.Filter("Author=", af)
	}
	it := lc.Client.Run(r.Context(), query)

	for {
		var b models.Book
		_, err := it.Next(&b)
		if err == iterator.Done {
			break
		}
		output = append(output, b)
	}

	j, _ := json.MarshalIndent(output, "", "  ")
	w.Write(j)
}
