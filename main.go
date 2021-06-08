package main

import (
	"flag"
	"fmt"
	"os"
)

func usage() {
	fmt.Println("Welcome to 노마드 코인")
	fmt.Println("Please use the following commands:")
	fmt.Println("explorer: Start the HTML explorer.")
	fmt.Println("rest: Start the rest api server.")
	os.Exit(0)
}

func main () {
	if len(os.Args) < 2  {
		usage()
	}

	rest:=flag.NewFlagSet("rest", flag.ExitOnError)
	portFlag:=rest.Int("port", 4000, "Sets the port of the server.")
	flags:=os.Args[2:]

	switch os.Args[1] {
	case "explorer":
		fmt.Println("Start explorer")
	case "rest":
		rest.Parse(flags)
	default:
		usage()
	}

	if rest.Parsed() {
		fmt.Printf("Start with port: %d", *portFlag)
	}
}