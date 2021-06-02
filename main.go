package main

import "fmt"

func plus(a ...int) int {
	var sum int = 0
	for _, item := range(a) {
		sum += item
	}
	return sum
}



func main() {
	result := plus(2,3, 4, 5, 6, 7, 8, 9, 10)
	foods := [] string {"potato", "pizza", "chicken"}

	fmt.Println(result)

	foods = append(foods, "tomato",)
	for _, food := range foods {
		fmt.Println(food);
	}

}