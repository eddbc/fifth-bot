package fifth

import (
	"testing"
	"github.com/gregjones/httpcache"
	"net/http"
	"github.com/antihax/goesi"
)

func init() {
	transport := httpcache.NewTransport(httpcache.NewMemoryCache())
	transport.Transport = &http.Transport{Proxy: http.ProxyFromEnvironment}
	client := &http.Client{Transport: transport}

	// Get our API Client.
	esiClient := goesi.NewAPIClient(client, useragent)
	Eve = esiClient.ESI
}

func TestWhoCommand(t *testing.T){

	dm, err := getCharacterInfoEmbed("Edd Reynolds")
	if err != nil {
		t.FailNow()
		return
	}

	char := dm.Title

	alli := ""
	corp := ""
	for _, v := range dm.Fields {
		if v.Name == "Alliance" {
			alli = v.Value
		}
		if v.Name == "Corporation" {
			corp = v.Value
		}
	}

	if char != "Edd Reynolds" {
		t.Error("Name does not match expected value")
	}

	if corp != "Girls Lie But Zkill Doesn't" {
		t.Error("Corp does not match expected value")
	}

	if alli != "Pandemic Legion" {
		t.Error("Alliance does not match expected value")
	}
}
