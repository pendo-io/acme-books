package pendo

// Names of the Track Events sent by the API. Pendo matches events to their
// track type by exact, case-sensitive name, so these values must not change.
const (
	EventBookCatalogSearched = "book_catalog_searched" // GET /books returned the book list
	EventBookDetailsViewed   = "book_details_viewed"   // GET /books/:id returned a book
	EventBookLookupFailed    = "book_lookup_failed"    // GET /books/:id did not return a book
)

// Names of Track Events registered in Pendo for endpoints that are not built
// yet (readme exercises 4 and 5). Send them from those handlers once they exist.
const (
	EventBookBorrowed     = "book_borrowed"      // PUT /:id/borrow marked the book as borrowed
	EventBookBorrowFailed = "book_borrow_failed" // PUT /:id/borrow rejected the request
	EventBookReturned     = "book_returned"      // PUT /:id/return marked the book as returned
	EventBookReturnFailed = "book_return_failed" // PUT /:id/return rejected the request
	EventBookAdded        = "book_added"         // POST /book stored a new book
)
