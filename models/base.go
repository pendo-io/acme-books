package models

import (
	"cloud.google.com/go/datastore"
	"context"
	"google.golang.org/api/iterator"
)

type Filter struct {
	Operation string
	Value interface{}
}

func (filter Filter) applyTo(query *datastore.Query) *datastore.Query {
	return query.Filter(filter.Operation, filter.Value)
}

type modelListener struct {
	itemFound func([]datastore.Property)
}

func getAll(itemType string, filter *Filter, listener modelListener) error {
	client, ctx := getClient()
	defer client.Close()
	query := datastore.NewQuery(itemType).Order("Id")
	if filter != nil {
		query = filter.applyTo(query)
	}
    it := client.Run(ctx, query)
    for {
		var itemProperties datastore.PropertyList
        if _, err := it.Next(&itemProperties); err != nil {
			if err == iterator.Done {
				err = nil
				break
			} else {
				return err
			}
		}
        listener.itemFound(itemProperties)
    }
    return nil
}

func update(itemType string, itemId int64, entity interface{}) error {
	client, ctx := getClient()
	defer client.Close()
	key := datastore.IDKey(itemType, itemId, nil)
	if _, err := client.Put(ctx, key, entity); err != nil {
		return err
	}
	return nil
}

func create(itemType string, entity interface{}) error {
	client, ctx := getClient()
	defer client.Close()
	key := datastore.IncompleteKey(itemType, nil)
	if _, err := client.Put(ctx, key, entity); err != nil {
		return err
	}
	return nil
}


func getClient() (*datastore.Client, context.Context) {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "acme-books")
	return client, ctx
}


