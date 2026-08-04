// Package figuras es un paquete propio de ejemplo: agrupa funciones para
// calcular áreas y perímetros de figuras geométricas simples.
//
// Fíjate en la convención de nombres: las funciones que empiezan con
// MAYÚSCULA (AreaCirculo, AreaRectangulo, PerimetroRectangulo) son
// EXPORTADAS: cualquier otro paquete que importe "figuras" puede usarlas.
// La función esValida (minúscula) es PRIVADA: solo el código dentro de este
// mismo paquete puede llamarla. Desde main.go, en la raíz del proyecto,
// intentar llamar a figuras.esValida daría un error de compilación.
package figuras

import "math"

// esValida es un helper interno: valida que una medida sea positiva antes
// de usarla en un cálculo. No exportada porque es un detalle de
// implementación, no algo que quien use este paquete necesite conocer.
func esValida(medida float64) bool {
	return medida > 0
}

// AreaCirculo devuelve el área de un círculo dado su radio.
// Si el radio no es válido (<= 0), devuelve 0.
func AreaCirculo(radio float64) float64 {
	if !esValida(radio) {
		return 0
	}
	return math.Pi * radio * radio
}

// AreaRectangulo devuelve el área de un rectángulo dado base y altura.
// Si alguna medida no es válida (<= 0), devuelve 0.
func AreaRectangulo(base, altura float64) float64 {
	if !esValida(base) || !esValida(altura) {
		return 0
	}
	return base * altura
}

// PerimetroRectangulo devuelve el perímetro de un rectángulo dado base y
// altura. Si alguna medida no es válida (<= 0), devuelve 0.
func PerimetroRectangulo(base, altura float64) float64 {
	if !esValida(base) || !esValida(altura) {
		return 0
	}
	return 2 * (base + altura)
}
