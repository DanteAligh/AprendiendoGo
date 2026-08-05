package main

import (
	"fmt"
	"strconv"
)

func main() {
	entero := "42"
	convert, err := strconv.Atoi(entero)

	if err != nil {
		fmt.Println(err)
	}
	suma := convert + 8
	fmt.Println(suma)

}
