package fifth

import (
	"testing"
)

func TestWhoCommand(t *testing.T) {

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

	if corp != "Avalanche." {
		t.Error("Corp does not match expected value")
	}

	if alli != "Fraternity." {
		t.Error("Alliance does not match expected value")
	}
}
