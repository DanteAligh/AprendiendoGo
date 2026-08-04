// Día 21: funciones simples para practicar testing table-driven.
//
// Este archivo es REFERENCIA DE SINTAXIS. El código aquí es intencionalmente
// simple: lo importante hoy no es la lógica sino CÓMO se prueba.
package main

import (
	"errors"
	"fmt"
	"strings"
)

// Dividir devuelve un error explícito en vez de dejar que el programa
// entre en pánico (dividir entre cero con float64 no hace panic en Go,
// da +Inf/NaN — pero para una operación de negocio real, tratarlo como
// error es la forma correcta de comunicarlo a quien llama).
func Dividir(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("no se puede dividir entre cero")
	}
	return a / b, nil
}

// EsPalindromo verifica si un string se lee igual al derecho y al revés.
// Decisión de diseño (documentada para que los tests la reflejen):
// esta función NO ignora mayúsculas ni espacios. "Ana" y "ana" son
// palindromos distintos para efectos de esta función; "a nana" no cuenta
// como palindromo porque el espacio se compara tal cual.
func EsPalindromo(s string) bool {
	n := len(s)
	for i := 0; i < n/2; i++ {
		if s[i] != s[n-1-i] {
			return false
		}
	}
	return true
}

// ContarVocales cuenta cuántas vocales (a, e, i, o, u) tiene un string,
// sin distinguir mayúsculas de minúsculas.
func ContarVocales(s string) int {
	vocales := "aeiouAEIOU"
	total := 0
	for _, c := range s {
		if strings.ContainsRune(vocales, c) {
			total++
		}
	}
	return total
}

func main() {
	// main existe solo para que el paquete compile como programa ejecutable
	// (go run .). La parte interesante de este día vive en los tests.
	resultado, err := Dividir(10, 2)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Printf("10 / 2 = %.2f\n", resultado)
	}

	fmt.Println(`EsPalindromo("reconocer"):`, EsPalindromo("reconocer"))
	fmt.Println(`ContarVocales("murciélago"):`, ContarVocales("murciélago"))
}
