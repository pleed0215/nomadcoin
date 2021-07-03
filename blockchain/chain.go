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
	CurrentDifficulty int `json:"currentDifficulty"`
}
const (
	difaultDifficulty int = 2
	difficultyInterval int = 5
	blockInterval int = 2
	allowedRange int = 2
)


/*
 for singleton pattern
*/
var b *blockchain
var once sync.Once

func (b *blockchain) persist() {
	db.SaveCheckpoint(utils.ToBytes(b))
}

func (b *blockchain) AddBlock() {
	block := createBlock(b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = BC().difficulty()
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

func (b *blockchain) recalculateDifficulty() int {
	allBlocks := b.AllBlocks()
	newestBlock := allBlocks[0]
	lastRecalculatedBlock := allBlocks[difficultyInterval-1]
	actualTime := (newestBlock.Timestamp-lastRecalculatedBlock.Timestamp)/60
	expectedTime := difficultyInterval * blockInterval
	if actualTime <= (expectedTime-allowedRange) {
		return b.CurrentDifficulty + 1
	} else if actualTime > (expectedTime+allowedRange) {
		return b.CurrentDifficulty - 1
	}
	return b.CurrentDifficulty
}

func (b *blockchain) difficulty() int {
	if b.Height == 0 {
		return difaultDifficulty
	} else if b.Height % difficultyInterval == 0 {
		// recalc the difficulty
		return b.recalculateDifficulty()
	} else {
		return b.CurrentDifficulty
	}
}

func (b *blockchain) txOuts() []*TxOut {
	var txOuts []*TxOut
	blocks := b.AllBlocks()
	for _, block := range blocks {
		for _, tx := range block.Transactions {
			txOuts = append(txOuts, tx.TxOuts...)
		}
	}
	return txOuts
}

func (b *blockchain) TxOutsByAddress(address string) []*TxOut {
	var ownedTxOuts []*TxOut
	txOuts := b.txOuts()
	for _, txOut := range txOuts {
		if txOut.Owner == address {
			ownedTxOuts = append(ownedTxOuts, txOut)
		}
	}
	return ownedTxOuts
}

func (b *blockchain) BalanceByAddress(address string) int {
	txOuts := b.TxOutsByAddress(address)
	var amount int = 0
	for _, txOut := range txOuts {
		amount += txOut.Amount
	}
	return amount
}

func (b *blockchain) txIns() []*TxIn {
	var txIns []*TxIn
	blocks := b.AllBlocks()
	for _, block := range blocks {
		for _, tx := range block.Transactions {
			txIns = append(txIns, tx.TxIns...)
		}
	}
	return txIns
}

func BC() *blockchain {
	if b == nil {
		once.Do( func() {
			b = &blockchain{Height: 0}
			checkpoint := db.Blockchain()
			// search for checkpoint on the db
			// restore b from bytes
			fmt.Printf("NewestHash: %s\nHeight: %d\n", b.NewestHash, b.Height)
			if(checkpoint == nil) {
				b.AddBlock()
			} else {
				fmt.Println("Restoring...")
				b.restore(checkpoint)
			}
			fmt.Printf("NewestHash: %s\nHeight: %d\n", b.NewestHash, b.Height)
		})
	}
	return b
}