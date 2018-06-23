package main

import (
	"time"
	"fmt"
	"context"
	"log"
	"errors"
)


//
// Methods
//

/*
Get related data from provided IDs
E.g. get character and corporation names
 */
func (k *Kill) inflate() {
	if k.inflated != true {
		// zkill url
		k.Zkb.url = fmt.Sprintf("https://zkillboard.com/kill/%v/", k.KillmailID)

		ctx := context.Background()

		// get victim ship
		ship, _, err := eve.UniverseApi.GetUniverseTypesTypeId(ctx, int32(k.Victim.ShipTypeID), nil)
		if err == nil {
			k.Victim.ShipTypeName = ship.Name
		} else {log.Printf("inflate victim ship: %v", err)}

		// get victim character
		vic, _, err := eve.CharacterApi.GetCharactersCharacterId(ctx, int32(k.Victim.CharacterID), nil)
		if err == nil {
			k.Victim.CharacterName = vic.Name
		} else {log.Printf("inflate character: %v", err)}

		// get victim corp
		crp, _, err := eve.CorporationApi.GetCorporationsCorporationId(ctx, int32(k.Victim.CorporationID), nil)
		if err == nil {
			k.Victim.CorporationName = crp.Name
		} else {log.Printf("inflate corp: %v", err)}

		// get victim alliance
		ali, _, err := eve.AllianceApi.GetAlliancesAllianceId(ctx, int32(k.Victim.AllianceID), nil)
		if err == nil {
			k.Victim.AllianceName = ali.Name
		} else {log.Printf("inflate alliance: %v", err)}

		k.inflated = true
	}
}

func (k *Kill) getFinalBlow() (Attacker, error) {
	for _, at := range k.Attackers {
		if at.FinalBlow {
			fb, _, err := eve.CharacterApi.GetCharactersCharacterId(context.Background(), int32(at.CharacterID), nil)
			if err == nil {
				at.CharacterName = fb.Name
				return at, nil
			} else {log.Printf("inflate final blow character: %v", err)}
		}
	}

	return Attacker{}, errors.New("error occurred getting final blow character")
}

/*
Check if an entity (character, corp or alliance) was involved in a kill
 */
func (k *Kill) involved(entityId int) (bool) {
	return k.isAttacker(entityId) || k.isVictim(entityId)
}

/*
Check if an entity (character, corp or alliance) was an attacker in a kill
 */
func (k *Kill) isAttacker(entityId int) (bool) {
	for _, a := range k.Attackers {
		if a.CharacterID == entityId {
			return true
		}
		if a.CorporationID == entityId {
			return true
		}
		if a.AllianceID == entityId {
			return true
		}
	}
	return false
}

/*
Check if an entity (character, corp or alliance) was the victim in a kill
 */
func (k *Kill) isVictim(entityId int) (bool) {
	if k.Victim.CharacterID == entityId {
		return true
	}
	if k.Victim.CorporationID == entityId {
		return true
	}
	if k.Victim.AllianceID == entityId {
		return true
	}

	return false
}

//
// Structs
//

type Kill struct {
	Attackers []Attacker 	`json:"attackers"`
	KillmailID    int       `json:"killmail_id"`
	KillmailTime  time.Time `json:"killmail_time"`
	SolarSystemID int       `json:"solar_system_id"`
	Victim 		  Victim 	`json:"victim"`
	Zkb struct {
		LocationID  int     `json:"locationID"`
		Hash        string  `json:"hash"`
		FittedValue float64 `json:"fittedValue"`
		TotalValue  float64 `json:"totalValue"`
		Points      int     `json:"points"`
		Npc         bool    `json:"npc"`
		Solo        bool    `json:"solo"`
		Awox        bool    `json:"awox"`
		url			string
	} `json:"zkb"`
	inflated		bool
}

type Victim struct {
	*Character
	DamageTaken   int `json:"damage_taken"`
	Items         []struct {
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

type Attacker struct {
	*Character
	DamageDone     int     `json:"damage_done"`
	FinalBlow      bool    `json:"final_blow"`
	SecurityStatus float64 `json:"security_status"`
	WeaponTypeID   int     `json:"weapon_type_id"`
}

type Character struct {
	AllianceID     	int     `json:"alliance_id,omitempty"`
	AllianceName   	string
	CharacterID    	int     `json:"character_id,omitempty"`
	CharacterName  	string
	CorporationID  	int     `json:"corporation_id,omitempty"`
	CorporationName	string
	ShipTypeID     	int     `json:"ship_type_id,omitempty"`
	ShipTypeName	string
}