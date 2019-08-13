package fifth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
)

func getKill1() (Kill, error) {
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

func getKill2() (Kill, error) {
	var kill Kill
	r, err := http.Get("https://esi.evetech.net/latest/killmails/78390232/02d423efc03d9dd5c1ef2f769b53a940c1213df5/")
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

	kill1, err := getKill1()
	if err != nil {
		t.FailNow()
		return
	}

	kill2, err := getKill2()
	if err != nil {
		t.FailNow()
		return
	}

	isKill1, isLoss1, err := isEntityRelated(&kill1)
	if err != nil {
		t.FailNow()
		return
	}
	isKill2, isLoss2, err := isEntityRelated(&kill2)
	if err != nil {
		t.FailNow()
		return
	}

	if !isKill1 {
		t.Error("Example Kill 1 Not marked as kill when kill")
	}

	if !isLoss1 {
		t.Error("Example Kill 1 Not marked as loss when loss")
	}

	if isKill2 {
		t.Error("Example Kill 2 marked as kill when not kill")
	}

	if isLoss2 {
		t.Error("Example Kill 2 marked as loss when not loss")
	}

}
