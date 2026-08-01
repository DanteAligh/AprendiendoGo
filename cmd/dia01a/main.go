package main

import (
	"fmt"
)

func total(a, b float64) float64 {
	return a * b
}

func main() {
	var nombre string = "Teclado mecánico"
	precio := 1250.5
	unidad := 3
	disponible := true
	t := total(precio, float64(unidad))

	fmt.Printf("nombre: %s\n precio: %.2f\n unidad: %d\n disponible: %t\n total: %.2f\n", nombre, precio, unidad, disponible, t)
}
