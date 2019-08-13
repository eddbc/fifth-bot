package fifth

import (
	"github.com/antihax/goesi"
	"github.com/eddbc/fifth-bot/esiStatus"
	"github.com/eddbc/fifth-bot/storage"
	"github.com/gregjones/httpcache"
	bolt "go.etcd.io/bbolt"
	"log"
	"net/http"
	"testing"
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

func TestEsiStatus(t *testing.T) {
	g, y, r, err := esiStatus.GetEsiStatus()
	if err != nil {
		t.Fatalf("error getting status: %v", err)
	}

	t.Logf("ESI Status - green: %v, yellow: %v, red: %v", g, y, r)
}
