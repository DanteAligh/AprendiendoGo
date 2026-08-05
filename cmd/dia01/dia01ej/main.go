package main

import "fmt"

// Recibe grados Celsius (un número con decimales) y devuelve Fahrenheit.
func celsiusAFahrenheit(c float64) float64 {
	return c*9.0/5.0 + 32.0
}

// Devuelve true (sí) o false (no). A partir de 38 grados hay fiebre.
func esFiebre(tempC float64) bool {
	return tempC >= 38.0
}

func main() {
	temperaturas := []float64{36.5, 37.2, 38.4, 39.0}

	for _, t := range temperaturas {
		f := celsiusAFahrenheit(t)
		if esFiebre(t) {
			fmt.Printf("%.1f °C = %.1f °F -> FIEBRE\n", t, f)
		} else {
			fmt.Printf("%.1f °C = %.1f °F -> normal\n", t, f)
		}
	}
}
