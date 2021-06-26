package blockchain

import (
	"fmt"
	"strings"
	"time"

	"github.com/pleed0215/nomadcoin/db"
	"github.com/pleed0215/nomadcoin/utils"
)


type Block struct {
	Data string `json:"data"`
	Hash string `json:"hash"`
	PrevHash string `json:"prevHash,omitempty"`
	Height int `json:"height"`
	Difficulty int `json:"difficulty"`
	Nonce int `json:"nonce"`
	Timestamp int `json:"timestamp"`
}

func (b *Block) mine() {
	target := strings.Repeat("0", b.Difficulty)
	for {
		hash := utils.Hash(b)
		b.Timestamp = int(time.Now().Unix())
		fmt.Printf("Hash:%s\nTarget:%s\nNonce:%d\n\n\n", hash, target, b.Nonce)
		if strings.HasPrefix(hash, target) {
			b.Hash = hash
			break
		} else {
			b.Nonce++
		}
	}
}


func createBlock(data string, prevHash string, height int) *Block {
	block:=&Block{
		Data: data, 
		Hash: "", 
		PrevHash: prevHash, 
		Height: height, 
		Difficulty: BC().difficulty(), 
		Nonce:0,
	}

	block.mine()
	block.persist()

	return block
}



func (b *Block) persist() {
	db.SaveBlock(b.Hash, utils.ToBytes(b))
}

func (b *Block) restore(data []byte) {
	utils.FromBytes(b, data)
}

func FindBlock(hash string) (*Block, error) {
	blockBytes := db.Block(hash)
	if blockBytes == nil {
		return nil, ErrNotFound
	}
	block := Block{}
	block.restore(blockBytes)
	return &block, nil
}