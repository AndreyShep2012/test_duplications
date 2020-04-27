package main

import (
	"fmt"

	"github.com/AndreyShep2012/test_duplications/sum"
	"github.com/AndreyShep2012/test_duplications/sum2"
)

func main() {
	v1 := 5
	v2 := 10
	fmt.Println("sum of", v1, "and", v2, "is", sum.Sum(v1, v2))
	fmt.Println("sum of", v1, "and", v2, "is", sum2.Sum(v1, v2))
}
