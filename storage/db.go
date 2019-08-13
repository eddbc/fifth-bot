package storage

import (
	"encoding/binary"
	"fmt"
	"go.etcd.io/bbolt"
)

var DB *bbolt.DB

const TimersKey = "timers"
const TheraHolesKey = "theraHoles"

func DBInit() {
	DB.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(TimersKey))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	DB.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(TheraHolesKey))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}

func Put(bucket string, key string, data []byte) {
	DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return nil
		}
		b.Put([]byte(key), data)
		return nil
	})
}

func Get(bucket string, key string) ([]byte, error) {
	// Start the transaction.
	r := make([]byte, 0)
	tx, err := DB.Begin(false)
	if err != nil {
		return r, err
	}
	defer tx.Rollback()

	b := tx.Bucket([]byte(bucket))
	rb := b.Get([]byte(key))

	r = make([]byte, len(rb))
	copy(r, rb)

	// Commit the transaction.
	if err := tx.Commit(); err != nil {
		return r, err
	}

	return r, err
}

// Itob returns an 8-byte big endian representation of v.
func Itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
