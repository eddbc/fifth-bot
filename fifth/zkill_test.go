package fifth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

func getKill() (Kill, error) {
	var kill Kill
	r, err := http.Get("https://zkillboard.com/api/killID/72019508/")
	if err != nil {
		return kill, err
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	var j []Kill
	json.Unmarshal(body, &j)

	kill = j[0]

	return kill, err
}

func TestInterestingKill(t *testing.T) {

	entitiesOfInterest = []int32{
		1354830081, // goons
		99005338,   // horde
		98481691,   // nogrl
	}

	kill, err := getKill()
	if err != nil {
		t.FailNow()
		return
	}

	isKill, isLoss, err := isEntityRelated(kill)
	if err != nil {
		t.FailNow()
		return
	}

	if !isKill {
		t.Error("Not marked as kill when kill")
	}

	if !isLoss {
		t.Error("Not marked as loss when loss")
	}

}
