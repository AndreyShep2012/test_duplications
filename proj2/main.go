package main

import (
	"fmt"
)

func sum(v1, v2 int) int {
	return v1 + v2
}

func main() {
	v1 := 5
	v2 := 10

	test1 := "hello"
	test2 := "hello"

	fmt.Println("sum of", v1, "and", v2, "is", sum(v1, v2))

	fmt.Println(test1, test2)
}
