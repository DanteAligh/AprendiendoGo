// main.go es el punto de entrada de este mini-proyecto. Vive en la raíz del
// módulo (junto al go.mod) e importa nuestro propio paquete local "figuras"
// usando el module path completo: "aprendiendo-go/internal/figuras", donde
// "aprendiendo-go" es el nombre declarado en go.mod y "figuras" es el nombre
// de la carpeta/paquete.
//
// Este archivo demuestra cómo se organiza un proyecto real en varios
// paquetes: main.go no sabe (ni le importa) CÓMO están implementadas las
// funciones de figuras, solo conoce su firma pública (exportada).
package main

import (
	"fmt"

	"aprendiendo-go/internal/figuras"
)

func main() {
	fmt.Println("=== Cálculos con el paquete local 'figuras' ===")

	areaCirculo := figuras.AreaCirculo(5)
	fmt.Printf("Área de un círculo de radio 5: %.2f\n", areaCirculo)

	areaCirculoInvalido := figuras.AreaCirculo(-3)
	fmt.Printf("Área de un círculo de radio -3 (inválido): %.2f\n", areaCirculoInvalido)

	areaRect := figuras.AreaRectangulo(4, 6)
	fmt.Printf("Área de un rectángulo de 4x6: %.2f\n", areaRect)

	perimetroRect := figuras.PerimetroRectangulo(4, 6)
	fmt.Printf("Perímetro de un rectángulo de 4x6: %.2f\n", perimetroRect)

	// Lo siguiente, si lo descomentas, NO compila: esValida no está
	// exportada, así que main.go (otro paquete) no puede verla.
	// fmt.Println(figuras.esValida(5))
}
