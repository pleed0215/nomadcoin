package blockchain

import (
	"errors"
	"sync"

	"github.com/boltdb/bolt"
	"github.com/pleed0215/nomadcoin/db"
)

var ErrNotFound error = errors.New("block not found")

type Block struct {
	Height int `json:"height"`
	Data string `json:"data"`
	Hash string `json:"hash"`
	PrevHash string `json:"prevHash,omitempty"`
}

type blockchain struct {
	blocks []*Block
}



/*
 for singleton pattern
*/
var b *blockchain
var once sync.Once

func (b *blockchain) AddBlock(data string) {
	block := Block
	db.DB().Update(func(t *bolt.Tx) error {})
}

func GetBlockchain() *blockchain {
	if b == nil {
		once.Do( func() {
			b = &blockchain{}
			b.AddBlock("Genesis Block")
		})
	}
	return b
}