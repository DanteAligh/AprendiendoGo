// Día 15: sync.WaitGroup, sync.Mutex, race conditions, select y sync.Once.
//
// Este archivo es REFERENCIA DE SINTAXIS, no la solución de los ejercicios.
// Corre este archivo con:
//
//	go run ./cmd/extra/sync-select
//
// Para ver el race detector en acción (sección 2), prueba también:
//
//	go run -race ./cmd/extra/sync-select
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	demoWaitGroup()
	demoMutex()
	demoSelect()
	demoOnce()
}

// -----------------------------------------------------------------------
// 1. sync.WaitGroup: esperar a que N goroutines terminen.
// -----------------------------------------------------------------------
//
// WaitGroup es un contador. Add(n) suma trabajo pendiente, cada goroutine
// llama Done() al terminar, y Wait() bloquea hasta que el contador llega a
// cero. No hay canal de por medio: solo te interesa "¿ya acabaron todos?".
func demoWaitGroup() {
	fmt.Println("--- 1. sync.WaitGroup ---")

	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1) // avisamos que hay una goroutine más pendiente

		// OJO: pasamos i como parámetro (numero) en vez de capturar la
		// variable del for directamente. Si no lo hiciéramos así, todas
		// las goroutines podrían ver el mismo valor final de i (bug muy
		// común). Desde Go 1.22 el for ya crea una variable nueva por
		// iteración, pero pasar el valor explícito sigue siendo la forma
		// más clara y portable de expresar la intención.
		go func(numero int) {
			defer wg.Done() // se ejecuta al terminar la goroutine, pase lo que pase
			time.Sleep(10 * time.Millisecond)
			fmt.Printf("  goroutine %d terminó\n", numero)
		}(i)
	}

	wg.Wait() // bloquea hasta que las 5 goroutines llamaron Done()
	fmt.Println("  todas las goroutines terminaron")
	fmt.Println()
}

// -----------------------------------------------------------------------
// 2. sync.Mutex: proteger un dato compartido contra race conditions.
// -----------------------------------------------------------------------
//
// Una race condition ocurre cuando dos o más goroutines leen/escriben la
// misma variable sin coordinación. contador++ NO es atómico: es
// "leer, sumar 1, escribir", y sin protección, incrementos concurrentes se
// pueden perder.
func demoMutex() {
	fmt.Println("--- 2. sync.Mutex (contador protegido) ---")

	var mu sync.Mutex
	var wg sync.WaitGroup
	contador := 0

	const numGoroutines = 50
	const incrementosPorGoroutine = 1000

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < incrementosPorGoroutine; j++ {
				mu.Lock() // solo una goroutine a la vez entra aquí
				contador++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	esperado := numGoroutines * incrementosPorGoroutine
	fmt.Printf("  contador final: %d (esperado: %d)\n", contador, esperado)
	if contador == esperado {
		fmt.Println("  correcto: el mutex evitó que se perdieran incrementos")
	}
	fmt.Println()
}

// -----------------------------------------------------------------------
// 3. select: elegir entre varios channels, el que esté listo primero.
// -----------------------------------------------------------------------
func demoSelect() {
	fmt.Println("--- 3. select con channels y timeout ---")

	canalA := make(chan string)
	canalB := make(chan string)

	go func() {
		time.Sleep(30 * time.Millisecond)
		canalA <- "mensaje desde A"
	}()

	go func() {
		time.Sleep(60 * time.Millisecond)
		canalB <- "mensaje desde B"
	}()

	// select espera a que CUALQUIERA de los cases esté listo.
	// time.After crea un channel que "envía" después de la duración dada:
	// es la forma idiomática de poner un timeout a una espera.
	for i := 0; i < 2; i++ {
		select {
		case msg := <-canalA:
			fmt.Println("  recibido:", msg)
		case msg := <-canalB:
			fmt.Println("  recibido:", msg)
		case <-time.After(200 * time.Millisecond):
			fmt.Println("  timeout: nadie respondió a tiempo")
		}
	}
	fmt.Println()
}

// -----------------------------------------------------------------------
// 4. sync.Once: ejecutar algo una sola vez, sin importar cuántas
//    goroutines lo pidan concurrentemente.
// -----------------------------------------------------------------------

type conexion struct {
	conectado bool
	once      sync.Once
}

func (c *conexion) Conectar() {
	c.once.Do(func() {
		fmt.Println("  conectando... (esto debe imprimirse UNA sola vez)")
		time.Sleep(20 * time.Millisecond)
		c.conectado = true
	})
}

func demoOnce() {
	fmt.Println("--- 4. sync.Once ---")

	c := &conexion{}
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Conectar()
		}()
	}

	wg.Wait()
	fmt.Println("  conectado:", c.conectado)
}
