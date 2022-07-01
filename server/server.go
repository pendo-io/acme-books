package server

import (
	"context"

	"cloud.google.com/go/datastore"
)

func Init(host string, port string,client *datastore.Client , ctx context.Context) {
	r := NewRouter(client, ctx)
	r.RunOnAddr(host + ":" + port)
}


