package main

import (
	bc "github.com/pleed0215/nomadcoin/blockchain"
	"github.com/pleed0215/nomadcoin/explorer"
)



func main() {

	chain := bc.GetBlockchain()
	chain.AddBlock("Genesis block")
	chain.AddBlock("Second block")
	chain.AddBlock("Third block")

	explorer.Start()
}