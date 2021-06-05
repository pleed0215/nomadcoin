package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	bc "github.com/pleed0215/nomadcoin/blockchain"
)

const PORT string = ":4000"

type homeData struct {
	Title string
	Blocks []*bc.Block
}

func home(rw http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/home.gohtml")) 
	data := homeData{"GoGOGOOO", bc.GetBlockchain().AllBlock()}
	tmpl.Execute(rw, data)
}

func main() {
	http.HandleFunc("/", home)
	chain := bc.GetBlockchain()
	chain.AddBlock("Genesis block")
	chain.AddBlock("Second block")
	chain.AddBlock("Third block")
	fmt.Printf("Listening on http://localhost%s\n", PORT)
	log.Fatal(http.ListenAndServe(PORT, nil)) 
	

	chain.ListBlocks()
}