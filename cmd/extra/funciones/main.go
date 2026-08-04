// Día 5 — Ejemplo: funciones con parámetros, retorno múltiple, named
// returns, funciones variádicas y closures.
//
// Caso de uso de este ejemplo: validar contraseñas y generar "creadores de
// descuento" para una tienda. Los ejercicios de ejercicios.md usan otros
// casos (mayoría de edad, calculadora, promedio, multiplicador) — este
// ejemplo NO los resuelve.
package main

import (
	"fmt"
	"strings"
)

// Retorno simple: una función que recibe un parámetro y devuelve un solo
// valor. Nota que el tipo va DESPUÉS del nombre del parámetro.
func longitudValida(clave string) bool {
	return len(clave) >= 8
}

// Retorno múltiple: devuelve el resultado Y un booleano que indica si la
// validación pasó. Este es EL patrón que vas a ver por todo Go (a menudo
// el segundo valor es un "error" en vez de un bool, pero el patrón es igual).
func validarClave(clave string) (bool, string) {
	if len(clave) < 8 {
		return false, "la clave debe tener al menos 8 caracteres"
	}
	if !strings.ContainsAny(clave, "0123456789") {
		return false, "la clave debe contener al menos un número"
	}
	return true, "clave válida"
}

// Named returns: nombramos los valores de retorno en la firma
// (esCorta, motivo), lo que los "pre-declara" como variables dentro de la
// función. El "return" desnudo al final devuelve lo que tengan en ese momento.
func analizarLongitud(clave string) (categoria string, longitud int) {
	longitud = len(clave)
	switch {
	case longitud < 6:
		categoria = "muy corta"
	case longitud < 10:
		categoria = "aceptable"
	default:
		categoria = "larga"
	}
	return // equivale a "return categoria, longitud"
}

// Función variádica: acepta cualquier cantidad de strings. Dentro de la
// función, "claves" se comporta como un slice de strings.
func contarClavesValidas(claves ...string) int {
	total := 0
	for _, clave := range claves {
		if ok, _ := validarClave(clave); ok {
			total++
		}
	}
	return total
}

// Función de orden superior que devuelve un closure: la función devuelta
// "recuerda" el valor de "porcentaje" que recibió al ser creada, incluso
// después de que creadorDeDescuento ya terminó de ejecutarse.
func creadorDeDescuento(porcentaje float64) func(float64) float64 {
	return func(precioOriginal float64) float64 {
		return precioOriginal - (precioOriginal * porcentaje / 100)
	}
}

func main() {
	fmt.Println("--- Retorno simple (bool) ---")
	fmt.Println("¿'abc' tiene longitud válida?:", longitudValida("abc"))
	fmt.Println("¿'contraseña1' tiene longitud válida?:", longitudValida("contraseña1"))

	fmt.Println("\n--- Retorno múltiple ---")
	valida, motivo := validarClave("abc123")
	fmt.Println("¿'abc123' es válida?:", valida, "->", motivo)
	valida, motivo = validarClave("miClave2024")
	fmt.Println("¿'miClave2024' es válida?:", valida, "->", motivo)

	fmt.Println("\n--- Named returns ---")
	categoria, longitud := analizarLongitud("hola")
	fmt.Printf("'hola' -> categoría: %s, longitud: %d\n", categoria, longitud)
	categoria, longitud = analizarLongitud("estaEsUnaClaveLarga")
	fmt.Printf("'estaEsUnaClaveLarga' -> categoría: %s, longitud: %d\n", categoria, longitud)

	fmt.Println("\n--- Función variádica ---")
	// Se puede llamar con cualquier cantidad de argumentos.
	fmt.Println("Claves válidas:", contarClavesValidas("abc", "clave123", "otraClave9", "1234567"))

	fmt.Println("\n--- Closures ---")
	// Dos closures creados con distinto "porcentaje", cada uno recuerda el
	// suyo de forma completamente independiente.
	descuentoVerano := creadorDeDescuento(20)  // 20% de descuento
	descuentoNavidad := creadorDeDescuento(35) // 35% de descuento
	fmt.Printf("Precio 100 con descuento de verano: %.2f\n", descuentoVerano(100))
	fmt.Printf("Precio 100 con descuento de navidad: %.2f\n", descuentoNavidad(100))
}
