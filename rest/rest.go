package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	bc "github.com/pleed0215/nomadcoin/blockchain"
	"github.com/pleed0215/nomadcoin/utils"
)

var port string
const BASE_URL string = "http://localhost"

type url string

func (u url) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("%s%s%s", BASE_URL, port, u)), nil
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

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []urlDescription{
		{
			URL: "/",
			Method: "GET",
			Description: "See Documentation",
		},
		{
			URL: "/blocks",
			Method: "POST",
			Description: "Create a block",
			Payload:"data:string",
		},
		{
			URL: "/blocks",
			Method: "GET",
			Description: "See all blocks",
		},
		{
			URL: "/blocks/{id}",
			Method: "GET",
			Description: "See a block",
			Payload: "id:int",
		},
	}
	
	rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(data)
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rw.Header().Add("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(bc.GetBlockchain().AllBlock())

	case "POST":
		var addBlockBody addBlockBody
		utils.HandleError(json.NewDecoder(r.Body).Decode(&addBlockBody))
		bc.GetBlockchain().AddBlock(addBlockBody.Message)
		rw.WriteHeader(http.StatusCreated)
	}
}

func Start(aPort int) {
	newServer := http.NewServeMux()
	port=fmt.Sprintf(":%d", aPort)
	newServer.HandleFunc("/", documentation)
	newServer.HandleFunc("/blocks", blocks)
	fmt.Println("listening on http://localhost", port)
	log.Fatal(http.ListenAndServe(port, newServer))
}