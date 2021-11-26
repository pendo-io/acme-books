package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"

	"acme-books/databases"
	"acme-books/models"
)

type LibraryController struct{}

func (lc LibraryController) ListAll(r *http.Request, w http.ResponseWriter) {
	writer, _ := r.URL.Query()["writer"]

	client, _ := databases.NewDatabaseClient()

	defer client.Close()

	output := make([]models.Book, 0)

	it := client.Run(datastore.NewQuery("Book").Filter("Author =", writer[0]).Order("Id"))
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

func (lc LibraryController) BorrowBook(r *http.Request, w http.ResponseWriter) {
	client, _ := databases.NewDatabaseClient()
	pathSections := strings.Split(r.URL.Path, "/")
	log.Print((pathSections[2]))

	defer client.Close()

	id, _ := strconv.ParseInt(pathSections[2], 0, 64)
	it := client.Run(datastore.NewQuery("Book").Filter("Id =", id))
	log.Print("The book is")
	var b models.Book
	_, err := it.Next(&b)
	log.Print(b)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else if b.Borrowed == true {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Book is already borrowed"))
	} else {
		b.Borrowed = true
		k := datastore.NameKey("Book", b.Title, nil)
		log.Print(b)
		if _, err := client.Put(k, &b); err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	}
}

func (lc LibraryController) ReturnBook(r *http.Request, w http.ResponseWriter) {
	client, _ := databases.NewDatabaseClient()
	pathSections := strings.Split(r.URL.Path, "/")
	log.Print((pathSections[2]))

	defer client.Close()

	id, _ := strconv.ParseInt(pathSections[2], 0, 64)
	it := client.Run(datastore.NewQuery("Book").Filter("Id =", id))
	log.Print("The book is")
	var b models.Book
	_, err := it.Next(&b)
	log.Print(b)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	} else if b.Borrowed == false {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Book was not borrowed!"))
	} else {
		b.Borrowed = false
		k := datastore.NameKey("Book", b.Title, nil)
		log.Print(b)
		if _, err := client.Put(k, &b); err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	}
}
