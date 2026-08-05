package main

import (
	"fmt"
	"math"
)

func main() {
	PI := 3.1416
	Radio := 4

	resultado := math.Pow(float64(Radio), 2)
	Area := PI * resultado

	fmt.Println(Area)
}
