package main

import "fmt"

func main() {
	var numero1 int = 10
	var numero2 float64 = 10.5

	// Corrección: Usamos %d para enteros, %f (o %.1f) para float64 y %T para los tipos
	fmt.Printf("numero1: %d (%T), numero2: %.1f (%T)\n", numero1, numero1, numero2, numero2)

	// Corrección: Usamos 10.0 para asegurar la compatibilidad con float64
	igual := (numero1 > 5 && numero2 < 10.0) || numero1 == 100

	// Aquí %t es correcto porque 'igual' es una variable booleana
	fmt.Printf("igual: %t\n", igual)
}
