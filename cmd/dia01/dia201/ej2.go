package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	var nombre string
	var edad int

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("ingresa nombre completo")
	scanner.Scan()
	nombre = scanner.Text()

	fmt.Println("ingresa edad")
	fmt.Scan(&edad)

	fmt.Printf("Hola, %s, tienes %d años\n", nombre, edad)
}
