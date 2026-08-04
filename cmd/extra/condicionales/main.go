// Día 3 — Ejemplo: if/else, switch (con y sin expresión), y por qué Go no
// tiene operador ternario.
//
// Caso de uso de este ejemplo: clasificar el "nivel de riesgo" de una
// temperatura corporal y decidir un mensaje según el día de la semana.
// Los ejercicios de ejercicios.md usan otros casos (par/impar, mayor de 3
// números, calificaciones, calculadora) — este ejemplo NO los resuelve.
package main

import "fmt"

func main() {
	temperatura := 38.4 // en grados Celsius, valor fijo para el ejemplo

	// --- if / else if / else clásico ---
	// Fíjate: sin paréntesis alrededor de la condición, y las llaves son
	// obligatorias aunque el cuerpo sea de una sola línea.
	if temperatura < 37.5 {
		fmt.Println("Temperatura normal")
	} else if temperatura < 38.0 {
		fmt.Println("Febrícula leve")
	} else if temperatura < 39.0 {
		fmt.Println("Fiebre moderada")
	} else {
		fmt.Println("Fiebre alta, buscar atención médica")
	}

	// --- if con declaración corta previa ---
	// La variable "diferencia" solo existe dentro de este bloque if/else.
	// Este patrón es MUY común en Go real, sobre todo con funciones que
	// devuelven (valor, error).
	if diferencia := temperatura - 37.0; diferencia > 1.5 {
		fmt.Printf("Está %.1f grados por encima de lo normal\n", diferencia)
	} else {
		fmt.Println("Diferencia manejable respecto a lo normal")
	}

	// --- switch sin expresión ("switch verdadero") ---
	// Cada case es una condición booleana independiente. Es una alternativa
	// más legible a una cadena larga de if/else if cuando hay varias
	// condiciones excluyentes entre sí.
	var nivelRiesgo string
	switch {
	case temperatura >= 39.5:
		nivelRiesgo = "crítico"
	case temperatura >= 38.0:
		nivelRiesgo = "alto"
	case temperatura >= 37.5:
		nivelRiesgo = "moderado"
	default:
		nivelRiesgo = "bajo"
	}
	fmt.Println("Nivel de riesgo:", nivelRiesgo)

	// --- switch con expresión (comparación directa de un valor) ---
	// A diferencia de C/Java/JavaScript, NO hace falta "break": cada case
	// termina solo al ejecutarse, no hay fall-through automático.
	diaSemana := 3
	switch diaSemana {
	case 1:
		fmt.Println("Lunes: arranca la semana")
	case 2, 3, 4: // varios valores en un mismo case, separados por coma
		fmt.Println("Mitad de semana: sigue así")
	case 5:
		fmt.Println("Viernes: ya casi")
	case 6, 7:
		fmt.Println("Fin de semana: a descansar")
	default:
		fmt.Println("Número de día no válido")
	}

	// --- Go no tiene operador ternario ---
	// En otros lenguajes escribirías algo como:
	//   estado = temperatura >= 38.0 ? "con fiebre" : "sin fiebre"
	// En Go, la forma idiomática es un if/else normal que asigna a una
	// variable ya declarada. Se ve "más largo", pero es más explícito y
	// fácil de leer de un vistazo, incluso anidado con otras condiciones.
	var estado string
	if temperatura >= 38.0 {
		estado = "con fiebre"
	} else {
		estado = "sin fiebre"
	}
	fmt.Println("Estado del paciente:", estado)
}
