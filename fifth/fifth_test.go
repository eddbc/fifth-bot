package fifth

import (
	"github.com/antihax/goesi"
	"github.com/eddbc/fifth-bot/storage"
	"github.com/gregjones/httpcache"
	bolt "go.etcd.io/bbolt"
	"log"
	"net/http"
)

func init() {
	transport := httpcache.NewTransport(httpcache.NewMemoryCache())
	transport.Transport = &http.Transport{Proxy: http.ProxyFromEnvironment}
	client := &http.Client{Transport: transport}

	// Get our API Client.
	esiClient := goesi.NewAPIClient(client, useragent)
	Eve = esiClient.ESI

	db, err := bolt.Open("../fifth.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	storage.DB = db
}
