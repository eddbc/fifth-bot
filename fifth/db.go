package fifth

import (
	"fmt"
	"go.etcd.io/bbolt"
)

var DB *bbolt.DB

const timersKey = ""

func DBInit() {
	DB.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(timersKey))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		//b.Put([]byte("test key"), []byte("test value"))
		return nil
	})
}

func Get(key string) ([]byte, error) {
	// Start the transaction.
	r := make([]byte, 0)
	tx, err := DB.Begin(false)
	if err != nil {
		return r, err
	}
	defer tx.Rollback()

	b := tx.Bucket([]byte("test"))
	rb := b.Get([]byte(key))

	r = make([]byte, len(rb))
	copy(r, rb)

	// Commit the transaction.
	if err := tx.Commit(); err != nil {
		return r, err
	}

	return r, err
}
