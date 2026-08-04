// Día 9: Métodos — receptor por valor vs receptor por puntero
//
// El caso de uso de los ejercicios es Contador, CuentaBancaria, Rectangulo y
// ListaTareas. Aquí usamos un caso distinto: un "Termostato" con historial,
// para mostrar la misma lección con otro ejemplo.
package main

import "fmt"

// Termostato guarda la temperatura actual y cuántas veces se ha ajustado.
type Termostato struct {
	TemperaturaActual float64
	Ajustes           int
}

// Leer es un método con receptor POR VALOR: solo lee datos, no necesita
// modificar el struct. Usar valor aquí es correcto y además dice, con solo
// leer la firma, "este método no va a cambiar tu Termostato".
func (t Termostato) Leer() string {
	return fmt.Sprintf("Temperatura actual: %.1f°C (ajustado %d veces)", t.TemperaturaActual, t.Ajustes)
}

// AjustarPorPuntero SÍ modifica el struct, así que el receptor DEBE ser
// puntero (*Termostato). Si fuera por valor, la modificación desaparecería
// en cuanto el método terminara: solo cambiaríamos una copia temporal.
func (t *Termostato) AjustarPorPuntero(nuevaTemp float64) {
	t.TemperaturaActual = nuevaTemp
	t.Ajustes++
}

// AjustarPorValorMal es el "anti-ejemplo" a propósito: parece que ajusta la
// temperatura, compila sin ningún error, pero el cambio nunca sale del
// método porque "t" aquí es una COPIA del Termostato original.
func (t Termostato) AjustarPorValorMal(nuevaTemp float64) {
	t.TemperaturaActual = nuevaTemp
	t.Ajustes++
	// t se descarta al terminar el método; el Termostato original ni se entera.
}

func main() {
	term := Termostato{TemperaturaActual: 20.0}

	fmt.Println(term.Leer())

	// Llamamos el método "malo": compila perfecto, pero no cambia nada afuera.
	term.AjustarPorValorMal(30.0)
	fmt.Println("Después de AjustarPorValorMal(30.0):")
	fmt.Println(term.Leer()) // sigue en 20.0 -> la prueba de que no funcionó

	// Llamamos el método correcto, con puntero.
	term.AjustarPorPuntero(30.0)
	fmt.Println("Después de AjustarPorPuntero(30.0):")
	fmt.Println(term.Leer()) // ahora sí cambió a 30.0

	// Nota importante: aunque "term" es una variable normal (no un puntero),
	// Go nos deja llamar term.AjustarPorPuntero(...) directamente. Por
	// debajo, Go convierte esto a (&term).AjustarPorPuntero(...) porque
	// "term" es una variable direccionable. Esta conveniencia solo aplica a
	// variables; no funcionaría sobre un literal temporal como
	// Termostato{}.AjustarPorPuntero(10) (eso SÍ daría error de compilación).
	term.AjustarPorPuntero(22.5)
	fmt.Println("Después de otro ajuste:")
	fmt.Println(term.Leer())

	// Si tenemos explícitamente un puntero, la sintaxis es idéntica: Go
	// desreferencia por nosotros al leer campos o llamar métodos de valor.
	ptr := &term
	fmt.Println("Leyendo a través de un puntero explícito:", ptr.Leer())
}
