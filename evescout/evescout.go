package evescout

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var httpClient *http.Client

func init() {
	httpClient = &http.Client{}
}

func GetTheraHoles()(holes []Wormhole, err error){
	r, err := httpClient.Get("https://www.eve-scout.com/api/wormholes")
	if err != nil {
		return
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &holes)

	return
}