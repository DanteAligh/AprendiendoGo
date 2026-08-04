// Día 2 — Ejemplo: operadores, entrada por consola (fmt.Scan y bufio.Scanner)
// y formateo con fmt.Printf.
//
// Caso de uso de este ejemplo: un pequeño cálculo de "combustible restante
// en un viaje" (distancia, consumo, litros disponibles). Los ejercicios de
// ejercicios.md usan otros casos (suma/resta de dos números, nombre+edad,
// precio*cantidad) — este ejemplo NO los resuelve, solo enseña la sintaxis.
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// --- Entrada simple con fmt.Scan ---
	// fmt.Scan lee "hasta el siguiente espacio en blanco" (espacio, tab o
	// salto de línea). Necesita la DIRECCIÓN de memoria de la variable
	// (por eso el "&"), porque tiene que escribir el valor leído allí
	// mismo, no simplemente devolverlo.
	var kmARecorrer float64
	var kmPorLitro float64
	fmt.Print("¿Cuántos km vas a recorrer? ")
	fmt.Scan(&kmARecorrer)
	fmt.Print("¿Cuántos km hace tu carro por litro? ")
	fmt.Scan(&kmPorLitro)

	// --- Operadores aritméticos ---
	// División entre dos float64 sí da decimales (a diferencia de int/int,
	// que trunca). Aquí ambos operandos ya son float64, así que no hace
	// falta convertir nada.
	litrosNecesarios := kmARecorrer / kmPorLitro

	// --- Entrada de línea completa con bufio.Scanner ---
	// A diferencia de fmt.Scan, esto lee la línea COMPLETA, incluyendo
	// espacios. Útil para leer texto libre como un nombre de ciudad
	// compuesto ("Bogotá D.C.") o, en este caso, el nombre del conductor.
	lector := bufio.NewScanner(os.Stdin)
	fmt.Print("¿Cuál es tu nombre? ")
	lector.Scan()
	nombreConductor := lector.Text()

	// --- Operadores de comparación y lógicos ---
	tanqueLleno := 40.0 // litros que caben en el tanque, como ejemplo fijo
	alcanzaElViaje := litrosNecesarios <= tanqueLleno
	viajeLargo := kmARecorrer > 300
	// && (Y lógico): ambas condiciones deben cumplirse.
	// || (O lógico): basta con que una se cumpla.
	necesitaParadaExtra := viajeLargo && !alcanzaElViaje

	// --- Formateo con Printf ---
	// %s: string | %.1f: decimal con 1 posición | %v: valor "genérico"
	// %T: el TIPO de la variable (muy útil para depurar)
	// %t: booleano
	fmt.Printf("\nConductor: %s\n", nombreConductor)
	fmt.Printf("Litros necesarios: %.1f (tipo: %T)\n", litrosNecesarios, litrosNecesarios)
	fmt.Printf("¿Alcanza con un tanque lleno de %v litros?: %t\n", tanqueLleno, alcanzaElViaje)
	fmt.Printf("¿Necesita parada extra de gasolina?: %t\n", necesitaParadaExtra)
}
