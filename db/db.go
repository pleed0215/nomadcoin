package db

import (
	"github.com/boltdb/bolt"
	"github.com/pleed0215/nomadcoin/utils"
)

var db *bolt.DB

const dbName string = "blockchain.db"
const dataBucket string = "data"
const blocksBucket string ="blocks"

func DB() *bolt.DB {
	defer db.Close()
	if db == nil {

		dbPointer, err := bolt.Open(dbName, 0600, nil)
		utils.HandleError(err)
		db = dbPointer

		err = db.Update(func (t *bolt.Tx) error {
			_, err := t.CreateBucketIfNotExists([]byte(dataBucket))
			utils.HandleError(err)
			_, err = t.CreateBucketIfNotExists([]byte(blocksBucket))
			utils.HandleError(err)

			return nil
		})
	}
	return db
}