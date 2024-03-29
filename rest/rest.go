package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pleed0215/nomadcoin/blockchain"
	bc "github.com/pleed0215/nomadcoin/blockchain"
	"github.com/pleed0215/nomadcoin/utils"
)

var port string
const BASE_URL string = "http://localhost"

type url string
func (u url) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("%s%s%s", BASE_URL, port, u)), nil
}

type errorResponse struct {
	ErrorMessage string `json:"errorMessage"`
}

type balanceResponse struct {
	Address string `json:"address"`
	Balance int `json:"balance"`
}

type urlDescription struct {
	URL url `json:"url"`
	Method string `json:"method"`
	Description string `json:"description"`
	Payload string `json:"payload,omitempty"`
}

type addBlockBody struct {
	Message string
}

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []urlDescription{
		{
			URL: url("/"),
			Method: "GET",
			Description: "See Documentation",
		},
		{
			URL: url("/blocks"),
			Method: "POST",
			Description: "Create a block",
			Payload:"data:string",
		},
		{
			URL: url("/status"),
			Method: "GET",
			Description: "See Status",
		},
		{
			URL: url("/blocks"),
			Method: "GET",
			Description: "See all blocks",
		},
		{
			URL: url("/blocks/{height}"),
			Method: "GET",
			Description: "See a block",
			Payload: "height:int",
		},
		{
			URL: url("/balance/{address}"),
			Method: "GET",
			Description: "Get Transaction Outputs for an Address",
			Payload: "address:string",
		},
	}
	
	json.NewEncoder(rw).Encode(data)
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":

		json.NewEncoder(rw).Encode(bc.BC().AllBlocks())

	case "POST":
		var addBlockBody addBlockBody
		utils.HandleError(json.NewDecoder(r.Body).Decode(&addBlockBody))
		bc.BC().AddBlock()
		rw.WriteHeader(http.StatusCreated)
	default:
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func block(rw http.ResponseWriter, r *http.Request) {
	if(r.Method == "GET") {
		vars := mux.Vars(r)
		hash := vars["hash"]
		block, err := bc.FindBlock(hash)
		encoder := json.NewEncoder(rw)
		if err == nil {
			encoder.Encode(block);
			rw.WriteHeader(http.StatusOK)
		} else {
			errorResponse := errorResponse{fmt.Sprintf("Block height: %s is not found", hash)}
			encoder.Encode(errorResponse)
			rw.WriteHeader(http.StatusNotFound)
		}
	}
}

func status(rw http.ResponseWriter, r *http.Request) {
	json.NewEncoder(rw).Encode(blockchain.BC())
}

func balance(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	total := r.URL.Query().Get("total");
	switch total {
	case "true":
		balance := blockchain.BC().BalanceByAddress(address)
		utils.HandleError(
			json.NewEncoder(rw).
			Encode(
				balanceResponse{
					Address: address, 
					Balance: balance,
				},
			),
		)
	default:
		utils.HandleError(json.NewEncoder(rw).Encode(blockchain.BC().TxOutsByAddress(address))) 
	}
	
}

func Start(aPort int) {
	newServer := mux.NewRouter()
	port=fmt.Sprintf(":%d", aPort)
	newServer.Use(jsonMiddleware)
	newServer.HandleFunc("/", documentation).Methods("GET")
	newServer.HandleFunc("/status", status).Methods("GET")
	newServer.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	newServer.HandleFunc("/block/{hash:[a-f0-9]+}", block).Methods("GET")
	newServer.HandleFunc("/balance/{address}",balance)
	fmt.Println("listening on http://localhost", port)
	log.Fatal(http.ListenAndServe(port, newServer))
}