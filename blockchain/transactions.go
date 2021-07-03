package blockchain

import (
	"time"

	"github.com/pleed0215/nomadcoin/utils"
)

const (
	minerReward int = 50
)

type Tx struct {
	Id string `json:"id"`
	Timestamp int `json:"timestampe"`
	TxIns []* TxIn `json:"txIns"`
	TxOuts []* TxOut `json:"txOuts"`
}

type TxIn struct {
	Owner string
	Amount int
}

type TxOut struct {
	Owner string
	Amount int
}

func (t *Tx) getId() {
	t.Id = utils.Hash(t)
}

func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{
			"COINBASE", minerReward,
		},
	}
	txOuts := []*TxOut{
		{
			address, minerReward,
		},
	}
	tx:= Tx{
		Id: "",
		Timestamp: int(time.Now().Unix()),
		TxIns: txIns,
		TxOuts: txOuts,
	}
	tx.getId()
	return &tx
}