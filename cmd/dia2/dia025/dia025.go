package main

import "fmt"

func main() {
	var precio float64
	var cantidad int

	fmt.Println("Ingresa el precio y la cantidad (puedes usar Enter entre ellos):")
	// Si escribes "Juan", presionas Enter, y luego "25", funcionará perfectamente.
	fmt.Scan(&precio, &cantidad)

	total := precio * float64(cantidad)
	fmt.Printf("total : %.2f", total)

}
