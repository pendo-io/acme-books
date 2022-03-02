package server

import (
	"context"

	"cloud.google.com/go/datastore"
)

func Init(host string, port string) {
	ctx := context.Background()
	client, _ := datastore.NewClient(ctx, "acme-books")
	defer client.Close()

	r := NewRouter(ctx, client)
	r.RunOnAddr(host + ":" + port)
}
