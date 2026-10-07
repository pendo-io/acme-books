package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-martini/martini"

	"acme-books/model"
	"acme-books/pendo"
)

type Library struct{}

func (l Library) GetByKey(params martini.Params, w http.ResponseWriter) {
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book, err := model.BookImplementation{}.GetByKey(id)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonStr, err := json.MarshalIndent(book, "", "  ")

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonStr)
}

// ListAll handles GET /books. The list can be filtered with the title, writer
// and borrowed query parameters.
func (l Library) ListAll(r *http.Request, w http.ResponseWriter) {
	filter, err := newBookFilter(r.URL.Query())

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	books := model.BookImplementation{}.ListAll()
	filtered := filter.count() > 0

	if filtered {
		books = filter.apply(books)
	}

	jsonStr, err := json.MarshalIndent(books, "", "  ")

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if filtered {
		pendo.Track(r, "book_list_filtered", filter.trackProperties(len(books)))
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonStr)
}

// Borrow handles PUT /:id/borrow, responding 204 if the book was borrowed and
// 400 if the id is invalid or the book is already borrowed.
func (l Library) Borrow(r *http.Request, params martini.Params, w http.ResponseWriter) {
	book, err := setBorrowed(params["id"], true)

	if reason := failureReason(err); reason != "" {
		fmt.Println(err)
		pendo.Track(r, "book_borrow_failed", loanProperties(params["id"], book, reason))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pendo.Track(r, "book_borrowed", loanProperties(params["id"], book, ""))
	w.WriteHeader(http.StatusNoContent)
}

// Return handles PUT /:id/return, responding 204 if the book was returned and
// 400 if the id is invalid or the book is not borrowed.
func (l Library) Return(r *http.Request, params martini.Params, w http.ResponseWriter) {
	book, err := setBorrowed(params["id"], false)

	if reason := failureReason(err); reason != "" {
		fmt.Println(err)
		pendo.Track(r, "book_return_failed", loanProperties(params["id"], book, reason))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pendo.Track(r, "book_returned", loanProperties(params["id"], book, ""))
	w.WriteHeader(http.StatusNoContent)
}

// Add handles POST /book, saving the Book in the request body under a newly
// generated key and responding with the saved Book.
func (l Library) Add(r *http.Request, w http.ResponseWriter) {
	var book model.Book
	err := json.NewDecoder(r.Body).Decode(&book)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book, err = model.BookImplementation{}.Add(book)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonStr, err := json.MarshalIndent(book, "", "  ")

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pendo.Track(r, "book_added", map[string]interface{}{
		"book_id": strconv.FormatInt(book.Id, 10),
		"title":   book.Title,
		"author":  book.Author,
	})

	w.WriteHeader(http.StatusOK)
	w.Write(jsonStr)
}

// bookFilter holds the filters applied to the book list. Empty text and a nil
// borrowed are not applied.
type bookFilter struct {
	title    string
	writer   string
	borrowed *bool
}

// newBookFilter reads a bookFilter from the title, writer and borrowed query
// parameters, returning an error if borrowed is not a boolean.
func newBookFilter(query url.Values) (bookFilter, error) {
	filter := bookFilter{
		title:  strings.TrimSpace(query.Get("title")),
		writer: strings.TrimSpace(query.Get("writer")),
	}

	if value := strings.TrimSpace(query.Get("borrowed")); value != "" {
		borrowed, err := strconv.ParseBool(value)

		if err != nil {
			return filter, err
		}

		filter.borrowed = &borrowed
	}

	return filter, nil
}

// count returns the number of filters applied.
func (f bookFilter) count() int {
	count := 0

	if f.title != "" {
		count++
	}

	if f.writer != "" {
		count++
	}

	if f.borrowed != nil {
		count++
	}

	return count
}

// apply returns the books that pass every applied filter. Title and writer
// match any part of the text, ignoring case.
func (f bookFilter) apply(books []model.Book) []model.Book {
	filtered := []model.Book{}

	for _, book := range books {
		if containsFold(book.Title, f.title) &&
			containsFold(book.Author, f.writer) &&
			(f.borrowed == nil || book.Borrowed == *f.borrowed) {
			filtered = append(filtered, book)
		}
	}

	return filtered
}

// trackProperties returns the properties of the book_list_filtered Pendo
// event for a filtered list of results books.
func (f bookFilter) trackProperties(results int) map[string]interface{} {
	properties := map[string]interface{}{
		"filter_count":  f.count(),
		"results_count": results,
	}

	if f.title != "" {
		properties["filter_title"] = f.title
	}

	if f.writer != "" {
		properties["filter_author"] = f.writer
	}

	if f.borrowed != nil {
		properties["filter_borrowed"] = *f.borrowed
	}

	return properties
}

func containsFold(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// errInvalidID is returned for a book id that is not a positive integer.
var errInvalidID = errors.New("invalid book id")

// setBorrowed parses id and marks the book it identifies as borrowed or
// returned.
func setBorrowed(id string, borrowed bool) (model.Book, error) {
	key, err := strconv.Atoi(id)

	if err != nil || key <= 0 {
		return model.Book{}, fmt.Errorf("%w: %q", errInvalidID, id)
	}

	return model.BookImplementation{}.SetBorrowed(key, borrowed)
}

// failureReason returns the failure_reason sent to Pendo when err rejects a
// borrow or return request with 400, or "" for any other error.
func failureReason(err error) string {
	switch {
	case errors.Is(err, errInvalidID), errors.Is(err, model.ErrBookNotFound):
		return "invalid_id"
	case errors.Is(err, model.ErrAlreadyBorrowed):
		return "already_borrowed"
	case errors.Is(err, model.ErrNotBorrowed):
		return "not_borrowed"
	default:
		return ""
	}
}

// loanProperties returns the Pendo event properties of a borrow or return
// request. The book's title and author are included whenever it was found,
// which is every outcome but an invalid id.
func loanProperties(id string, book model.Book, reason string) map[string]interface{} {
	properties := map[string]interface{}{"book_id": id}

	if reason != "" {
		properties["failure_reason"] = reason
	}

	if reason != "invalid_id" {
		properties["title"] = book.Title
		properties["author"] = book.Author
	}

	return properties
}
