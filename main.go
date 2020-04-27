package main

import (
	"fmt"

	"github.com/AndreyShep2012/test_duplications/sum"
)

func main() {
	v1 := 5
	v2 := 10
	fmt.Println("sum of", v1, "and", v2, "is", sum.Sum(v1, v2))
}
