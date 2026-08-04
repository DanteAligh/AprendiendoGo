// Día 6 — Ejemplo: arrays vs slices, len, cap, append, slicing [a:b], copy.
//
// Caso de uso de este ejemplo: manejar los últimos resultados de un sensor
// de temperatura (tamaño fijo con array) frente a una lista dinámica de
// pedidos de un restaurante (tamaño variable con slice). Los ejercicios de
// ejercicios.md usan otros casos (nombres, filtrado de pares, etc.) — este
// ejemplo NO los resuelve.
package main

import "fmt"

func main() {
	// --- Array: tamaño FIJO, parte del tipo ---
	// [3]float64 y [5]float64 son tipos distintos. Un array de temperaturas
	// de las últimas 3 lecturas de un sensor es un buen caso: sabemos de
	// antemano que SIEMPRE van a ser exactamente 3.
	var ultimasLecturas [3]float64
	ultimasLecturas[0] = 21.5
	ultimasLecturas[1] = 22.0
	ultimasLecturas[2] = 21.8
	fmt.Println("--- Array de tamaño fijo ---")
	fmt.Println("Lecturas del sensor:", ultimasLecturas)
	fmt.Println("Longitud del array (siempre 3):", len(ultimasLecturas))

	// --- Slice: tamaño flexible ---
	// A diferencia del array, no sabemos cuántos pedidos va a tener un
	// restaurante en una noche — por eso un slice es la elección natural.
	pedidos := []string{"Mesa 1: Hamburguesa", "Mesa 3: Ensalada"}
	fmt.Println("\n--- Slice inicial ---")
	fmt.Println("Pedidos:", pedidos)
	fmt.Println("len:", len(pedidos), "cap:", cap(pedidos))

	// --- append: agregar elementos ---
	// IMPORTANTE: append puede devolver un slice que apunta a un array
	// subyacente DISTINTO (si no había espacio suficiente), por eso
	// SIEMPRE se reasigna el resultado a la variable.
	pedidos = append(pedidos, "Mesa 5: Pizza")
	pedidos = append(pedidos, "Mesa 2: Sopa")
	fmt.Println("\n--- Después de 2 append ---")
	fmt.Println("Pedidos:", pedidos)
	fmt.Println("len:", len(pedidos), "cap:", cap(pedidos))

	// --- make: crear un slice con capacidad inicial conocida ---
	// Si ya sabemos aproximadamente cuántos pedidos esperamos, podemos
	// reservar espacio de una vez con make, evitando que Go tenga que
	// redimensionar el array subyacente varias veces.
	mesasReservadas := make([]int, 0, 10) // longitud 0, capacidad 10
	mesasReservadas = append(mesasReservadas, 1, 4, 7)
	fmt.Println("\n--- make con capacidad reservada ---")
	fmt.Println("Mesas reservadas:", mesasReservadas)
	fmt.Println("len:", len(mesasReservadas), "cap:", cap(mesasReservadas))

	// --- Slicing [a:b]: sub-vista de un slice ---
	// pedidos[1:3] toma desde el índice 1 (incluido) hasta el 3 (excluido):
	// o sea, las posiciones 1 y 2.
	primerosDosNuevos := pedidos[1:3]
	fmt.Println("\n--- Slicing [1:3] ---")
	fmt.Println("Sub-slice:", primerosDosNuevos)

	// OJO: el sub-slice comparte memoria con el original. Si modificamos
	// un elemento del sub-slice, el original también cambia.
	primerosDosNuevos[0] = "Mesa 5: Pizza (¡cambio de pedido!)"
	fmt.Println("Slice original después de modificar el sub-slice:", pedidos)

	// --- copy: duplicar SIN compartir memoria ---
	// A diferencia del slicing, copy crea una copia independiente: modificar
	// la copia NO afecta al original.
	copiaIndependiente := make([]string, len(pedidos))
	copy(copiaIndependiente, pedidos)
	copiaIndependiente[0] = "Esto NO afecta al slice original"
	fmt.Println("\n--- copy ---")
	fmt.Println("Copia independiente:", copiaIndependiente)
	fmt.Println("Slice original, intacto:", pedidos)
}
