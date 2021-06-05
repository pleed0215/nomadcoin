package blockchain

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

type Block struct {
	Data string
	Hash string
	PrevHash string
}

type blockchain struct {
	blocks []*Block
}

func (b *Block) calculateHash()  {
	b.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(b.Data + b.PrevHash)))
}

func createBlock(data string) *Block {
	newBlock := Block{data, "", getLastHash()}
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