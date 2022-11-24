package models

import (
	"context"
	"log"

	"cloud.google.com/go/datastore"
)

var client *datastore.Client

// Setup initializes the database instance
func Setup() {
	ctx := context.Background()

	var err error
	client, err = datastore.NewClient(ctx, "acme-books")

	if err != nil {
		log.Fatalf("Failed to instantiate new DataStore client: %v", err)
	}
}
