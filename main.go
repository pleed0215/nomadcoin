package main

import (
	"fmt"

	Person "github.com/pleed0215/nomadcoin/person"
)




func main() {
	led := Person.Person{Name: "Eundeok", Age: 1}
	fmt.Println(led.SayHello())
	led.AddAge()
	fmt.Println(led)
}