package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
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
			URL: "/blocks/{height}",
			Method: "GET",
			Description: "See a block",
			Payload: "height:int",
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
		bc.BC().AddBlock(addBlockBody.Message)
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

func Start(aPort int) {
	newServer := mux.NewRouter()
	port=fmt.Sprintf(":%d", aPort)
	newServer.Use(jsonMiddleware)
	newServer.HandleFunc("/", documentation).Methods("GET")
	newServer.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	newServer.HandleFunc("/block/{hash:[a-f0-9]+}", block).Methods("GET")
	fmt.Println("listening on http://localhost", port)
	log.Fatal(http.ListenAndServe(port, newServer))
}