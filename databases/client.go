package databases

import (
	"context"

	"cloud.google.com/go/datastore"
)

type CloudDatastore struct {
	ctx    context.Context
	client *datastore.Client
}

var NewDatabaseClient = NewDatastoreDatabaseClient

func NewDatastoreDatabaseClient() (CloudDatastore, error) {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, "acme-books")

	return CloudDatastore{client: client, ctx: ctx}, err
}

func (cds CloudDatastore) Close() {
	cds.client.Close()
}

func (cds CloudDatastore) PutMulti(keys []*datastore.Key, src interface{}) (ret []*datastore.Key, err error) {
	return cds.client.PutMulti(cds.ctx, keys, src)
}

func (cds CloudDatastore) Put(key *datastore.Key, src interface{}) (ret *datastore.Key, err error) {
	return cds.client.Put(cds.ctx, key, src)
}

func (cds CloudDatastore) Run(q *datastore.Query) *datastore.Iterator {
	return cds.client.Run(cds.ctx, q)
}
