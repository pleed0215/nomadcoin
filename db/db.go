package db

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/pleed0215/nomadcoin/utils"
)

var db *bolt.DB

const dbName string = "blockchain.db"
const dataBucket string = "data"
const blocksBucket string ="blocks"
const checkpoint string = "checkpoint"

func DB() *bolt.DB {

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
		utils.HandleError(err)
	}
	return db
}

func SaveBlock(hash string, data []byte)  {
	fmt.Printf("Saving Block %s\nData: %b\n", hash, data)
	err := DB().Update(func (t *bolt.Tx) error  {
		bucket := t.Bucket([]byte(blocksBucket))
		err := bucket.Put([]byte(hash), data)
		
		return err
	})
	utils.HandleError(err)
}

func SaveBlockchain(data []byte) {
	err := DB().Update(func(t *bolt.Tx) error {
		b:=t.Bucket([]byte(dataBucket))
		err:=b.Put([]byte(checkpoint), data)
		return err
	})
	utils.HandleError(err)
}

func Blockchain() []byte {
	var data []byte
	err := DB().View( func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		data = bucket.Get([]byte(checkpoint))
		return nil
	})
	utils.HandleError(err)
	return data
}

func Block(hash string) []byte {
	var data []byte
	DB().View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(blocksBucket))
		data = bucket.Get([]byte(hash))
		return nil
	})

	return data
}