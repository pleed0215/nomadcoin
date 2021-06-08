package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/pleed0215/nomadcoin/explorer"
	"github.com/pleed0215/nomadcoin/rest"
)

func usage() {
	fmt.Println("Welcome to 노마드 코인")
	fmt.Println("Please use the following commands:")
	fmt.Println("-port=[PORT]		Must be integer. Set port to use.")
	fmt.Println("-mode=[MODE]		Choose html, rest or both. Set which server to use.")
	os.Exit(0)
}

func Start() {

	port:=flag.Int("port", 4000, "Sets the port of the server.")
	mode:=flag.String("mode", "both", "Choose html or rest. Set which server to use.\nMode both will start html server first with port given, and later start rest api server with port+1.")

	if len(os.Args) ==1 {
		flag.Usage()
		os.Exit(0)
	}

	flag.Parse()
	
	if flag.Parsed() {
		switch *mode {
		case "html":
			fmt.Println("Start html server.")
			explorer.Start(*port)
		case "rest":
			fmt.Println("Start rest api server.")
			rest.Start(*port)
		case "both":
			fmt.Println("Start server.")
			go explorer.Start(*port)
			rest.Start(*port+1)

		default:
			flag.Usage()
		}
	}
}