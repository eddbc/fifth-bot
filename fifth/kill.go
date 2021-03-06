package fifth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"time"
)

//
// Methods
//

func (k *Kill) getURL() string {
	return fmt.Sprintf("https://zkillboard.com/kill/%v/", k.KillmailID)
}

/*
Get related data from provided IDs
E.g. get character and corporation names
*/
func (k *Kill) inflate() {
	defer timeTrack(time.Now(), "inflate")
	if k.inflated != true {
		// zkill url
		k.Zkb.url = fmt.Sprintf("https://zkillboard.com/kill/%v/", k.KillmailID)

		ctx := context.Background()

		sys, _, err := Eve.UniverseApi.GetUniverseSystemsSystemId(ctx, int32(k.SolarSystemID), nil)
		if err == nil {
			k.SolarSystemName = sys.Name
			c, _, err := Eve.UniverseApi.GetUniverseConstellationsConstellationId(ctx, int32(sys.ConstellationId), nil)
			if err == nil {
				reg, _, err := Eve.UniverseApi.GetUniverseRegionsRegionId(ctx, int32(c.RegionId), nil)
				if err == nil {
					k.RegionID = reg.RegionId
					k.RegionName = reg.Name
				}
			}
		} else {
			log.Printf("inflate solar system: %v", err)
		}

		// get victim ship
		ship, _, err := Eve.UniverseApi.GetUniverseTypesTypeId(ctx, int32(k.Victim.ShipTypeID), nil)
		if err == nil {
			k.Victim.ShipTypeName = ship.Name
		} else {
			log.Printf("inflate victim ship: %v", err)
		}

		// get victim character
		vic, _, err := Eve.CharacterApi.GetCharactersCharacterId(ctx, int32(k.Victim.CharacterID), nil)
		if err == nil {
			k.Victim.CharacterName = vic.Name
		} else {
			k.Victim.CharacterName = "Someone"
			log.Printf("inflate character: %v", err)
		}

		// get victim corp
		crp, _, err := Eve.CorporationApi.GetCorporationsCorporationId(ctx, int32(k.Victim.CorporationID), nil)
		if err == nil {
			k.Victim.CorporationName = crp.Name
			k.Victim.CorporationTicker = crp.Ticker
		} else {
			log.Printf("inflate corp: %v", err)
		}

		// get victim alliance
		ali, _, err := Eve.AllianceApi.GetAlliancesAllianceId(ctx, int32(k.Victim.AllianceID), nil)
		if err == nil {
			k.Victim.AllianceName = ali.Name
			k.Victim.AllianceTicker = ali.Ticker
		} else {
			log.Printf("inflate alliance: %v", err)
		}

		k.inflated = true
	}
}

func (s byDamage) Len() int {
	return len(s)
}

func (s byDamage) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byDamage) Less(i, j int) bool {
	return s[i].DamageDone < s[j].DamageDone
}

func (k *Kill) getFinalBlow() (Attacker, error) {
	for _, at := range k.Attackers {
		if at.FinalBlow {
			fb, _, err := Eve.CharacterApi.GetCharactersCharacterId(context.Background(), int32(at.CharacterID), nil)
			if err == nil {
				at.CharacterName = fb.Name
				return at, nil
			}
			log.Printf("inflate final blow character: %v", err)
		}
	}

	return Attacker{}, errors.New("error occurred getting final blow character")
}

func (k *Kill) interestingName() (string, error) {

	defer timeTrack(time.Now(), "interestingName")
	name := "Someone"

	fnlBlw, err := k.getFinalBlow()
	if err == nil {
		name = fnlBlw.CharacterName
	}

	if k.InterestingAttackers != nil {
		atkrs := k.InterestingAttackers
		sort.Sort(byDamage(atkrs))
		for _, atk := range atkrs {
			a, _, err := Eve.CharacterApi.GetCharactersCharacterId(context.Background(), int32(atk.CharacterID), nil)
			if err == nil {
				name = a.Name
			}
		}
	}

	return name, err
}

/*
Check if an entity (character, corp or alliance) was involved in a kill
*/
func (k *Kill) involved(entityID int32) bool {
	return k.isAttacker(entityID) || k.isVictim(entityID)
}

/*
Check if an entity (character, corp or alliance) was an attacker in a kill
*/
func (k *Kill) isAttacker(entityID int32) bool {

	k.InterestingAttackers = nil

	for _, a := range k.Attackers {
		r := false

		if a.CharacterID == entityID {
			r = true
		}
		if a.CorporationID == entityID {
			r = true
		}
		if a.AllianceID == entityID {
			r = true
		}

		if r {
			if k.InterestingAttackers == nil {
				k.InterestingAttackers = []Attacker{a}
			} else {
				k.InterestingAttackers = append(k.InterestingAttackers, a)
			}
		}
	}

	return k.InterestingAttackers != nil
}

/*
Check if an entity (character, corp or alliance) was the victim in a kill
*/
func (k *Kill) isVictim(entityID int32) bool {
	loss := false
	if k.Victim.CharacterID == entityID {
		loss = true
	}
	if k.Victim.CorporationID == entityID {
		loss = true
	}
	if k.Victim.AllianceID == entityID {
		loss = true
	}

	return loss
}

func (c *Character) ticker() string {
	if c.AllianceTicker == "" {
		return c.CorporationTicker
	}
	return c.AllianceTicker
}

//
// Structs
//

type byDamage []Attacker

//Kill Killmail
type Kill struct {
	Attackers            []Attacker `json:"attackers"`
	InterestingAttackers []Attacker
	KillmailID           int32     `json:"killmail_id"`
	KillmailTime         time.Time `json:"killmail_time"`
	SolarSystemID        int32     `json:"solar_system_id"`
	SolarSystemName      string
	RegionID             int32
	RegionName           string
	Victim               Victim `json:"victim"`
	Zkb                  struct {
		LocationID  int     `json:"locationID"`
		Hash        string  `json:"hash"`
		FittedValue float64 `json:"fittedValue"`
		TotalValue  float64 `json:"totalValue"`
		Points      int     `json:"points"`
		Npc         bool    `json:"npc"`
		Solo        bool    `json:"solo"`
		Awox        bool    `json:"awox"`
		url         string
	} `json:"zkb"`
	inflated bool
}

//Victim Thing destroyed in a killmaill
type Victim struct {
	Character
	DamageTaken int `json:"damage_taken"`
	Items       []struct {
		Flag              int `json:"flag"`
		ItemTypeID        int `json:"item_type_id"`
		QuantityDropped   int `json:"quantity_dropped,omitempty"`
		Singleton         int `json:"singleton"`
		QuantityDestroyed int `json:"quantity_destroyed,omitempty"`
		Items             []struct {
			Flag              int `json:"flag"`
			ItemTypeID        int `json:"item_type_id"`
			QuantityDestroyed int `json:"quantity_destroyed"`
			Singleton         int `json:"singleton"`
		} `json:"items,omitempty"`
	} `json:"items"`
	Position struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
		Z float64 `json:"z"`
	} `json:"position"`
}

//Attacker Aggressor on a killmail
type Attacker struct {
	Character
	DamageDone     int     `json:"damage_done"`
	FinalBlow      bool    `json:"final_blow"`
	SecurityStatus float64 `json:"security_status"`
	WeaponTypeID   int     `json:"weapon_type_id"`
}

//Character EVE Online character on a killmail
type Character struct {
	AllianceID     int32 `json:"alliance_id,omitempty"`
	AllianceName   string
	AllianceTicker string

	CorporationID     int32 `json:"corporation_id,omitempty"`
	CorporationName   string
	CorporationTicker string

	CharacterID   int32 `json:"character_id,omitempty"`
	CharacterName string

	ShipTypeID   int32 `json:"ship_type_id,omitempty"`
	ShipTypeName string
}
