package main

import (
	"github.com/pleed0215/nomadcoin/cli"
	"github.com/pleed0215/nomadcoin/db"
)




func main () {
	defer db.Close()

	cli.Start()
}