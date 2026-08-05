package main

import (
	"fmt"
)

func main() {
	var num1, num2 int

	fmt.Println("Ingresa el primer numero: ")
	fmt.Scanln(&num1)
	fmt.Println("Ingresa el segundo numero: ")
	fmt.Scanln(&num2)

	suma := num1 + num2
	resta := num1 - num2
	multi := num1 * num2
	divi := num1 / num2
	residuo := num1 % num2

	fmt.Printf("Suma: %d\n", suma)
	fmt.Printf("Resta: %d\n", resta)
	fmt.Printf("Multiplicacion: %d\n", multi)
	fmt.Printf("Division: %d\n", divi)
	fmt.Printf("residuo: %d\n", residuo)
}
