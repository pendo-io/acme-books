package repository

import (
	"acme-books/models"
	"context"
)
type BookRepository interface {

	GetBook(ctx context.Context, id int)(models.Book, error)
	GetBooks(ctx context.Context, filters map[string]string)([]models.Book)

}
