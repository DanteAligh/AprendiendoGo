package main

import (
	"fmt"
	"strconv"
)

const TASA_CAMBIO = 4.00

func main() {
	precioDolar := 19.99
	precioPeso := precioDolar * TASA_CAMBIO
	str := strconv.FormatFloat(precioPeso, 'f', 2, 64)
	fmt.Println("El precio en pesos es: " + str)
}
