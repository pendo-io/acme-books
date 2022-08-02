package models

import (
	"context"
	"fmt"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

type Book struct {
	Id       int64
	Title    string `json:"title"`
	Author   string `json:"writer"`
	Borrowed bool   `json:"borrowed"`
}

type ModelInterface interface {
	Create(book Book) error
	Change(id int, book Book) error
	GetById(id int) (Book, error)
	List(query *datastore.Query) ([]Book, error)
}

type DatastoreQuery struct {
	Ctx    context.Context
	Client *datastore.Client
}

func (dq *DatastoreQuery) Close() error {
	return dq.Client.Close()
}

func (dq *DatastoreQuery) List(query *datastore.Query) ([]interface{}, error) {
	it := dq.Client.Run(dq.Ctx, query)
	var output []interface{}
	for {
		var b Book
		_, err := it.Next(&b)
		if err == iterator.Done {
			fmt.Println(err)
			break
		} else {
			output = append(output, b)
		}
	}
	return output, nil
}

func (dq *DatastoreQuery) Create(book Book) error {
	key := datastore.IncompleteKey("Book", nil)
	_, err := dq.Client.Put(dq.Ctx, key, book)
	return err
}

func (dq *DatastoreQuery) Change(id int, book Book) error {
	key := datastore.IDKey("Book", int64(id), nil)
	_, err := dq.Client.Put(dq.Ctx, key, book)
	return err
}

func (dq *DatastoreQuery) GetById(id int) (Book, error) {
	var book Book
	key := datastore.IDKey("Book", int64(id), nil)
	err := dq.Client.Get(dq.Ctx, key, &book)
	return book, err
}

func CreateConnection(project string, ctx context.Context) (DatastoreQuery, error) {
	client, err := datastore.NewClient(ctx, project)
	if err != nil {
		return DatastoreQuery{nil, nil}, err
	}
	return DatastoreQuery{
		ctx,
		client,
	}, nil
}
