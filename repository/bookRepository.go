package repository

import (
	"acme-books/models"
	"context"
)
type BookRepository interface {

	GetBook(ctx context.Context, id int)(models.Book, error)
	GetBooks(ctx context.Context)([]models.Book, error)

}
