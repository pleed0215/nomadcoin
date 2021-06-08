package main

import (
	"github.com/pleed0215/nomadcoin/explorer"
	"github.com/pleed0215/nomadcoin/rest"
)

func main () {
	go rest.Start(3001)
	explorer.Start(3000)
}