package models

import (
	"context"
	"fmt"
	"sort"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
)

type Book struct {
	Id       int64
	Title    string `json:"title"`
	Author   string `json:"writer"`
	Borrowed bool   `json:"borrowed"`
}

func GetAllBooks(client *datastore.Client, ctx context.Context) []Book {
	var output []Book
	it := client.Run(ctx, datastore.NewQuery("Book"))
	for {
		var b Book
		_, err := it.Next(&b)
		if err == iterator.Done {
			fmt.Println(err)
			break
		}
		output = append(output, b)
	}

	sort.Sort(ById(output))
	return output
}

// Implementing the sort.Interface interface
type ById []Book

func (b ById) Len() int           { return len(b) }
func (b ById) Less(i, j int) bool { return b[i].Id < b[j].Id }
func (b ById) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
