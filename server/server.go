package server

import "cloud.google.com/go/datastore"

func Init(host string, port string, client datastore.Client) {
	r := NewRouter(client)
	r.RunOnAddr(host + ":" + port)
}
