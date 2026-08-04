// Tests table-driven para operaciones.go.
//
// Convención de Go: el archivo de tests vive junto al código, termina en
// _test.go, y cada función de prueba empieza con Test y recibe
// *testing.T. go test descubre estas funciones automáticamente.
package main

import "testing"

// TestDividir usa el patrón table-driven: en vez de una función TestX por
// caso, definimos una tabla de casos y un solo loop que los recorre. Cada
// caso corre con t.Run(nombre, ...) para que, si falla, sepas EXACTAMENTE
// cuál falló por su nombre en la salida de "go test -v".
func TestDividir(t *testing.T) {
	casos := []struct {
		nombre      string
		a, b        float64
		esperado    float64
		esperaError bool
	}{
		{nombre: "división normal", a: 10, b: 2, esperado: 5, esperaError: false},
		{nombre: "división con decimales", a: 7, b: 2, esperado: 3.5, esperaError: false},
		{nombre: "división entre cero", a: 5, b: 0, esperado: 0, esperaError: true},
		{nombre: "dividendo negativo", a: -10, b: 2, esperado: -5, esperaError: false},
		{nombre: "ambos negativos", a: -10, b: -2, esperado: 5, esperaError: false},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			resultado, err := Dividir(c.a, c.b)

			if c.esperaError {
				if err == nil {
					t.Errorf("Dividir(%v, %v): se esperaba un error, pero no hubo ninguno", c.a, c.b)
				}
				return
			}

			if err != nil {
				t.Fatalf("Dividir(%v, %v): error inesperado: %v", c.a, c.b, err)
			}
			if resultado != c.esperado {
				t.Errorf("Dividir(%v, %v) = %v; se esperaba %v", c.a, c.b, resultado, c.esperado)
			}
		})
	}
}

func TestEsPalindromo(t *testing.T) {
	casos := []struct {
		nombre   string
		entrada  string
		esperado bool
	}{
		{nombre: "palindromo simple", entrada: "ana", esperado: true},
		{nombre: "palindromo más largo", entrada: "reconocer", esperado: true},
		{nombre: "no es palindromo", entrada: "golang", esperado: false},
		{nombre: "string vacío es palindromo trivial", entrada: "", esperado: true},
		{nombre: "un solo caracter siempre es palindromo", entrada: "x", esperado: true},
		// Documentamos la decisión de diseño: esta función SÍ distingue
		// mayúsculas de minúsculas, así que "Ana" no es palindromo para
		// nuestra implementación (a != A).
		{nombre: "sensible a mayúsculas", entrada: "Ana", esperado: false},
		// Y tampoco ignora espacios: "a nana" no es palindromo porque el
		// espacio se compara tal cual, sin normalizar.
		{nombre: "sensible a espacios", entrada: "a nana", esperado: false},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			resultado := EsPalindromo(c.entrada)
			if resultado != c.esperado {
				t.Errorf("EsPalindromo(%q) = %v; se esperaba %v", c.entrada, resultado, c.esperado)
			}
		})
	}
}

func TestContarVocales(t *testing.T) {
	casos := []struct {
		nombre   string
		entrada  string
		esperado int
	}{
		{nombre: "todas consonantes", entrada: "xyz", esperado: 0},
		{nombre: "mezcla simple", entrada: "golang", esperado: 2},
		{nombre: "mayúsculas cuentan igual", entrada: "AEIOU", esperado: 5},
		{nombre: "string vacío", entrada: "", esperado: 0},
		{nombre: "repetidas", entrada: "aeiouaeiou", esperado: 10},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			resultado := ContarVocales(c.entrada)
			if resultado != c.esperado {
				t.Errorf("ContarVocales(%q) = %d; se esperaba %d", c.entrada, resultado, c.esperado)
			}
		})
	}
}
