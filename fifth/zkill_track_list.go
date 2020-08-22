package fifth

import (
	"context"
	"errors"
	"github.com/eddbc/fifth-bot/storage"
	"go.etcd.io/bbolt"
	"log"
)

type TrackedEntity struct {
	id   int32
	name string
}

func presetEntries() {
	var presets = []int32{

		// Alliances
		99007126, // NOG8S

		// Corps
		98326834, // NOG8S
		98539693, // NO0RE
		98506452, // SCHNG

		// Characters
		91872672, // Edd Reynolds
	}
	for _, e := range presets {
		err := addTrackedEntityByID(e)
		if err != nil {
			log.Print(err)
		}
	}
}

func getTrackedEntities() (trackedEntities []TrackedEntity) {
	storage.DB.View(func(tx *bbolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(storage.ZKillTrackedKey))
		b.ForEach(func(k, v []byte) error {
			trackedEntities = append(trackedEntities, TrackedEntity{
				id:   int32(storage.Btoi(k)),
				name: string(v),
			})
			return nil
		})
		return nil
	})
	//if Debug {
	//	trackedEntities = append(trackedEntities, getTestingEntities()...)
	//}
	return
}

//func getTestingEntities()(testingEntities []TrackedEntity) {
//	testingEntities = append(testingEntities, TrackedEntity{
//		id:   1354830081,
//		name: "goons",
//	})
//
//	testingEntities = append(testingEntities, TrackedEntity{
//		id:   99005338,
//		name: "horde",
//	})
//
//	testingEntities = append(testingEntities, TrackedEntity{
//		id:   498125261,
//		name: "test",
//	})
//
//	return
//}

func addTrackedEntityByID(id int32) (err error) {
	var i = []int32{id}
	res, _, err := Eve.UniverseApi.PostUniverseNames(context.Background(), i, nil)
	if err == nil {
		addTrackedEntity(id, res[0].Name)
	}

	return
}

func addTrackedEntityByName(name string) (err error) {
	var n = []string{name}
	res, _, err := Eve.UniverseApi.PostUniverseIds(context.Background(), n, nil)
	if err == nil {
		l := len(res.Characters) + len(res.Alliances) + len(res.Corporations)
		if l == 0 {
			err = errors.New("no results found")
		} else if l > 1 {
			err = errors.New("multiple results Found")
		} else {
			if len(res.Characters) == 1 {
				addTrackedEntity(res.Characters[0].Id, res.Characters[0].Name)
			} else if len(res.Corporations) == 1 {
				addTrackedEntity(res.Corporations[0].Id, res.Corporations[0].Name)
			} else if len(res.Alliances) == 1 {
				addTrackedEntity(res.Alliances[0].Id, res.Alliances[0].Name)
			}
		}
	}

	return
}

func addTrackedEntity(id int32, name string) {
	_ = storage.DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(storage.ZKillTrackedKey))
		if b == nil {
			return nil
		}
		_ = b.Put(storage.Itob(int(id)), []byte(name))
		return nil
	})
}
