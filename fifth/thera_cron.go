package fifth

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/eddbc/fifth-bot/evescout"
	"github.com/eddbc/fifth-bot/storage"
	"go.etcd.io/bbolt"
	"strconv"
)

func theraCron() {

	newHoles, err := evescout.GetTheraHoles()
	if err != nil {
		return
	}

	var oldHoles []evescout.Wormhole
	storage.DB.View(func(tx *bbolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(storage.TheraHolesKey))

		b.ForEach(func(k, v []byte) error {
			r := evescout.Wormhole{}
			json.Unmarshal(v, &r)
			oldHoles = append(oldHoles, r)
			return nil
		})
		return nil
	})

	for _, nWh := range newHoles {
		isNew := true
		for _, oWh := range oldHoles {
			if nWh.ID == oWh.ID {
				isNew = false
			}
		}
		if isNew {
			go newHole(nWh)
		}
	}

	for _, oWh := range oldHoles {
		dead := true
		for _, nWh := range newHoles {
			if nWh.ID == oWh.ID {
				dead = false
			}
		}
		if dead {
			storage.DB.Update(func(tx *bbolt.Tx) error {
				b := tx.Bucket([]byte(storage.TheraHolesKey))
				b.Delete(storage.Itob(int(oWh.ID)))
				return nil
			})
		}
	}
}

func newHole(wh evescout.Wormhole){
	//log.Printf("New Wormhole - Sig: %v ID: %v", wh.SignatureID, wh.ID)

	storage.DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(storage.TheraHolesKey))
		if b == nil {
			return nil
		}

		buf, _ := json.Marshal(wh)

		b.Put(storage.Itob(int(wh.ID)), buf)
		return nil
	})

	jita := int32(30000142)

	//Jita role
	//610241709753761802

	//Test Lab role
	//370156955005616128

	route, _, err := Eve.RoutesApi.GetRouteOriginDestination(context.Background(), wh.DestinationSolarSystem.ID, jita, nil)
	if err == nil {
		jumps := len(route)-1
		if jumps <= 5 {
			jumpsStr := fmt.Sprintf("%v jumps", strconv.Itoa(jumps))
			m := fmt.Sprintf("%v - %v %v [%.1f] (%v)", wh.SignatureID, jumpsStr, wh.DestinationSolarSystem.Name, wh.DestinationSolarSystem.Security, wh.DestinationSolarSystem.Region.Name)
			SendMsg(fmt.Sprintf("<@&610241709753761802> New Jita Hole - %v", m))
		}
	}
}
