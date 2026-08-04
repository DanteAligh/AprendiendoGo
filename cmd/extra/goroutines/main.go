// Día 14: Goroutines y channels básicos
//
// El reto integrador de los ejercicios pide un "procesador de pedidos" con
// structs+interfaces+errores+goroutines. Aquí mostramos, con un caso
// distinto (un "recolector de precios" simulando consultas a proveedores),
// la sintaxis base de goroutines y channels antes de que la uses en tu
// propia solución.
package main

import (
	"fmt"
	"time"
)

// consultarPrecio simula una consulta lenta (por ejemplo, a un proveedor
// externo o una API de IA) usando time.Sleep. Envía su resultado a un
// channel en vez de "retornarlo" normal, porque esta función va a correr
// dentro de una goroutine: una goroutine no puede "devolver" un valor con
// return de forma que otra goroutine lo reciba directamente, así que
// channels son el mecanismo para comunicar el resultado.
func consultarPrecio(proveedor string, resultado chan<- string) {
	// chan<- string en la firma dice: "esta función SOLO puede enviar a este
	// channel", no recibir. Es una restricción opcional pero útil: deja
	// clara la intención con solo leer la firma.
	tiempoSimulado := time.Duration(200+len(proveedor)*50) * time.Millisecond
	time.Sleep(tiempoSimulado)
	resultado <- fmt.Sprintf("%s respondió en %v", proveedor, tiempoSimulado)
}

func main() {
	fmt.Println("--- Parte 1: llamada secuencial vs goroutine ---")

	saludar := func(nombre string) {
		fmt.Println("Hola,", nombre)
	}

	// Llamadas secuenciales normales: cada una termina antes de que empiece
	// la siguiente. El orden de salida es 100% predecible.
	saludar("Ana")
	saludar("Beto")

	// Llamadas como goroutines: "go" arranca la función y el programa
	// continúa DE INMEDIATO, sin esperar a que termine. Por eso agregamos
	// un time.Sleep al final de main: si no le diéramos tiempo, el programa
	// podría terminar (y cerrar el proceso) antes de que las goroutines
	// alcancen a imprimir algo.
	go saludar("Carla (goroutine)")
	go saludar("Diego (goroutine)")
	time.Sleep(50 * time.Millisecond)

	fmt.Println("\n--- Parte 2: channel sin buffer, sincronización directa ---")

	listo := make(chan string) // sin buffer: enviar BLOQUEA hasta que alguien reciba
	go func() {
		fmt.Println("(goroutine) trabajando...")
		time.Sleep(300 * time.Millisecond)
		listo <- "trabajo terminado"
	}()

	fmt.Println("(main) esperando el resultado del channel...")
	mensaje := <-listo // esto bloquea aquí hasta que la goroutine envíe algo
	fmt.Println("(main) recibido:", mensaje)

	fmt.Println("\n--- Parte 3: channel CON buffer ---")

	numeros := make(chan int, 3) // buffer de 3: podemos enviar 3 veces sin que nadie reciba todavía
	numeros <- 10
	numeros <- 20
	numeros <- 30
	// Si intentáramos un cuarto "numeros <- 40" aquí sin haber recibido
	// nada, el envío bloquearía porque el buffer ya estaría lleno.
	close(numeros) // cerramos el channel: ya no se va a enviar nada más

	// range sobre un channel recibe valores hasta que el channel se cierra.
	for n := range numeros {
		fmt.Println("Recibido del buffer:", n)
	}

	fmt.Println("\n--- Parte 4: varias goroutines enviando al mismo channel ---")

	proveedores := []string{"ProveedorA", "ProveedorB", "ProveedorC"}
	resultados := make(chan string, len(proveedores)) // buffer del tamaño exacto que necesitamos

	for _, p := range proveedores {
		go consultarPrecio(p, resultados)
	}

	// Recibimos exactamente tantos valores como goroutines lanzamos.
	for i := 0; i < len(proveedores); i++ {
		fmt.Println("Respuesta recibida:", <-resultados)
	}

	fmt.Println("\nListo: todas las consultas concurrentes terminaron.")
}
