package main

import (
	"crypto/sha256"
	"fmt"
)

type block struct {
	data string
	hash string
	prevHash string
}



func main() {
	genesisBlock := block{"Genesis block", "", ""}
	genesisBlock.hash = fmt.Sprintf("%x", sha256.Sum256([]byte(genesisBlock.data+genesisBlock.prevHash)))

	secondBlock := block{"Second block", "", genesisBlock.hash}
	secondBlock.hash = fmt.Sprintf("%x", sha256.Sum256([]byte(secondBlock.data+secondBlock.prevHash)))
}