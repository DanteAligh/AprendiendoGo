package main

import "fmt"

// Go infiere tipos con :=, pero var es obligatorio a nivel de paquete.
var version = "1.0"

// Una función puede devolver varios valores. Es la base del manejo
// de errores en Go: (resultado, error).
func dividir(a, b float64) (float64, bool) {
	if b == 0 {
		return 0, false
	}
	return a / b, true
}

func main() {
	var nombre string = "Go"
	edad := 16 // int por inferencia
	pi := 3.14159

	fmt.Printf("Lenguaje: %s, edad: %d, pi: %.2f, version: %s\n", nombre, edad, pi, version)

	resultado, ok := dividir(10, 3)
	if !ok {
		fmt.Println("no se puede dividir entre cero")
		return
	}
	fmt.Printf("10/3 = %.4f\n", resultado)
}
