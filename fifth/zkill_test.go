package fifth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

func getKill() (Kill, error) {
	var kill Kill
	r, err := http.Get("https://esi.evetech.net/latest/killmails/72019508/1c768cdb3e5e8eb5a7d7dd12192fe8a7fa83a4d6/")
	if err != nil {
		return kill, err
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	json.Unmarshal(body, &kill)

	//kill = j[0]

	return kill, err
}

func TestSystemRange(t *testing.T) {
	var h8 int32 = 30000974
	var e02 int32 = 30000903

	d, err := distanceBetweenSystems(h8, e02)
	if err != nil {
		t.FailNow()
	}
	if d != 2.073 {
		t.Errorf("Distance 1 was %v, not 2.073", d)
	}
}

func TestInterestingKill(t *testing.T) {

	entitiesOfInterest = []int32{
		1354830081, // goons
		99005338,   // horde
	}

	kill, err := getKill()
	if err != nil {
		t.FailNow()
		return
	}

	isKill, isLoss, err := isEntityRelated(&kill)
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
