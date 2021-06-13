package blockchain

import (
	"errors"
	"fmt"
	"sync"

	"github.com/pleed0215/nomadcoin/db"
	"github.com/pleed0215/nomadcoin/utils"
)

var ErrNotFound error = errors.New("block not found")



type blockchain struct {
	NewestHash string `json:"newsetHash"`
	Height	   int	  `json:"height"`
}



/*
 for singleton pattern
*/
var b *blockchain
var once sync.Once

func (b *blockchain) persist() {
	db.SaveCheckpoint(utils.ToBytes(b))
}

func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.persist()
}

func (b *blockchain) AllBlocks() []*Block {
	var blocks []*Block

	hashCursor := b.NewestHash

	for {
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		} else {
			break
		}
	}

	return blocks
}

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(b, data)
}

func BC() *blockchain {
	if b == nil {
		once.Do( func() {
			b = &blockchain{"", 0}
			checkpoint := db.Blockchain()
			// search for checkpoint on the db
			// restore b from bytes
			fmt.Printf("NewestHash: %s\nHeight: %d\n", b.NewestHash, b.Height)
			if(checkpoint == nil) {
				b.AddBlock("Genesis Block")
			} else {
				fmt.Println("Restoring...")
				b.restore(checkpoint)
			}
			fmt.Printf("NewestHash: %s\nHeight: %d\n", b.NewestHash, b.Height)
		})
	}
	return b
}