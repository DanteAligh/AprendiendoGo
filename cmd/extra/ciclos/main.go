// Día 4 — Ejemplo: las tres formas del for, range, break/continue, ciclos
// anidados y labels.
//
// Caso de uso de este ejemplo: simular un contador de inventario que se va
// agotando, y recorrer una cuadrícula de asientos de un cine saltándose una
// fila específica. Los ejercicios de ejercicios.md usan otros casos (sumas,
// pares, tabla de multiplicar) — este ejemplo NO los resuelve.
package main

import "fmt"

func main() {
	// --- Forma clásica: inicialización; condición; actualización ---
	fmt.Println("--- Forma clásica ---")
	for i := 1; i <= 5; i++ {
		fmt.Println("Encendiendo motor, intento número", i)
	}

	// --- Forma "while": solo la condición ---
	// Go no tiene la palabra "while": el mismo "for" cubre este caso
	// cuando solo escribes la condición.
	fmt.Println("--- Forma while ---")
	inventario := 5
	for inventario > 0 {
		fmt.Println("Quedan", inventario, "unidades en inventario")
		inventario--
	}
	fmt.Println("Inventario agotado")

	// --- Forma infinita con break ---
	// Sin condición alguna: se repite para siempre hasta que algo interno
	// (aquí, un break) lo detenga.
	fmt.Println("--- Forma infinita ---")
	intentos := 0
	for {
		intentos++
		if intentos == 3 {
			fmt.Println("Conexión exitosa en el intento", intentos)
			break // sin esto, el ciclo nunca terminaría
		}
		fmt.Println("Intento de conexión", intentos, "fallido, reintentando...")
	}

	// --- range sobre un slice ---
	// range te da (índice, valor) en cada vuelta. Si no necesitas el
	// índice, se descarta con "_" (blank identifier).
	fmt.Println("--- range sobre un slice ---")
	frutas := []string{"manzana", "pera", "uva"}
	for indice, fruta := range frutas {
		fmt.Printf("Posición %d: %s\n", indice, fruta)
	}

	// --- continue ---
	fmt.Println("--- continue ---")
	for _, fruta := range frutas {
		if fruta == "pera" {
			continue // se salta la impresión solo para "pera"
		}
		fmt.Println("Comprando:", fruta)
	}

	// --- Ciclos anidados + label + continue sobre el ciclo externo ---
	// Simulamos una sala de cine de 4 filas x 4 asientos, pero la fila 2
	// está en mantenimiento y debe saltarse POR COMPLETO (no solo un
	// asiento). Sin el label, "continue" dentro del ciclo interno solo
	// afectaría al ciclo interno, no lograríamos saltar la fila entera.
	fmt.Println("--- Ciclos anidados con label ---")
filas:
	for fila := 1; fila <= 4; fila++ {
		for asiento := 1; asiento <= 4; asiento++ {
			if fila == 2 {
				continue filas // salta directo a la siguiente fila
			}
			fmt.Printf("Fila %d, Asiento %d disponible\n", fila, asiento)
		}
	}
}
