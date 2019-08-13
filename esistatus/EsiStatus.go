package esistatus

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

//GetEsiStatus Gets number of green, yellow, and red ESI endpoints
func GetEsiStatus() (green int, yellow int, red int, err error) {

	green = 0
	yellow = 0
	red = 0

	var endpoints []endpoint
	r, err := http.Get("https://esi.evetech.net/status.json")
	if err != nil {
		return
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &endpoints)
	if err != nil {
		return
	}

	for _, endpoint := range endpoints {
		if endpoint.Status == "green" {
			green++
		} else if endpoint.Status == "yellow" {
			yellow++
		} else if endpoint.Status == "red" {
			red++
		}
	}

	return
}

type endpoint struct {
	Endpoint string   `json:"endpoint"`
	Method   string   `json:"method"`
	Route    string   `json:"route"`
	Status   string   `json:"status"`
	Tags     []string `json:"tags"`
}
