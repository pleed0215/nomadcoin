package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"
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

func (b *Block) calculateHash()  {
	b.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(b.Data + b.PrevHash)))
}

func createBlock(data string) *Block {
	newBlock := Block{len(GetBlockchain().blocks)+1, data, "", getLastHash(), }
	newBlock.calculateHash()
	return &newBlock
}

func getLastHash() string {
	var length int = len(GetBlockchain().blocks)

	if length > 0 {
		return GetBlockchain().blocks[length-1].Hash;
	}
	return "";
}

func (b *blockchain) AddBlock(data string) {
	newBlock := createBlock(data)
	b.blocks = append(b.blocks, newBlock)
}

func (b *blockchain) ListBlocks() {
	for _, block := range(b.blocks) {
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %s\n", block.Hash)
		fmt.Printf("Previous Hash: %s\n\n", block.PrevHash)
	}
}

func (b *blockchain) AllBlock() []*Block {
	return b.blocks
}



func (b *blockchain) GetBlock(height int) (*Block, error) {
	for _, block := range b.blocks {
		if block.Height == height {
			return block, nil
		}
	}
	return nil,ErrNotFound
}

/*
 for singleton pattern
*/
var b *blockchain
var once sync.Once

func GetBlockchain() *blockchain {
	if b == nil {
		once.Do( func() {
			b = &blockchain{}
			b.AddBlock("Genesis Block")
		})
	}
	return b
}