// Día 17: JSON en Go con encoding/json.
//
// Este archivo es REFERENCIA DE SINTAXIS, no la solución de los ejercicios.
// Corre con:
//
//	go run ./cmd/extra/json
package main

import (
	"encoding/json"
	"fmt"
)

// Struct tags: le dicen a encoding/json cómo se llama cada campo en el
// JSON. Sin el tag, Go usaría el nombre del campo Go tal cual (Nombre, no
// nombre). omitempty le dice a Marshal que omita el campo si tiene el
// valor cero de su tipo (string vacío, 0, false, nil, slice vacío...).
type Producto struct {
	Nombre      string  `json:"nombre"`
	Precio      float64 `json:"precio"`
	Stock       int     `json:"stock"`
	Descripcion string  `json:"descripcion,omitempty"`
}

// Pedido demuestra JSON anidado: un slice de structs como campo.
type Pedido struct {
	Cliente   string     `json:"cliente"`
	Productos []Producto `json:"productos"`
}

func main() {
	demoMarshal()
	demoUnmarshal()
	demoOmitempty()
	demoAnidado()
	demoIndent()
}

// -----------------------------------------------------------------------
// 1. Marshal: de struct Go a JSON ([]byte).
// -----------------------------------------------------------------------
func demoMarshal() {
	fmt.Println("--- 1. json.Marshal ---")

	p := Producto{Nombre: "Teclado", Precio: 450.50, Stock: 12}

	datos, err := json.Marshal(p)
	if err != nil {
		fmt.Println("  error al convertir a JSON:", err)
		return
	}

	fmt.Println("  JSON:", string(datos))
	fmt.Println()
}

// -----------------------------------------------------------------------
//  2. Unmarshal: de JSON a struct Go. Siempre se pasa un puntero, porque
//     Unmarshal necesita ESCRIBIR en tu variable.
//
// -----------------------------------------------------------------------
func demoUnmarshal() {
	fmt.Println("--- 2. json.Unmarshal ---")

	jsonTexto := `{"nombre":"Teclado","precio":450.50,"stock":12}`

	var p Producto
	if err := json.Unmarshal([]byte(jsonTexto), &p); err != nil {
		fmt.Println("  error al parsear JSON:", err)
		return
	}

	fmt.Printf("  nombre=%s precio=%.2f stock=%d\n", p.Nombre, p.Precio, p.Stock)
	fmt.Println()
}

// -----------------------------------------------------------------------
// 3. omitempty: un campo vacío desaparece del JSON de salida.
// -----------------------------------------------------------------------
func demoOmitempty() {
	fmt.Println("--- 3. omitempty ---")

	conDescripcion := Producto{Nombre: "Mouse", Precio: 220, Stock: 30, Descripcion: "Óptico inalámbrico"}
	sinDescripcion := Producto{Nombre: "Cable HDMI", Precio: 90, Stock: 50} // Descripcion queda en ""

	datos1, _ := json.Marshal(conDescripcion)
	datos2, _ := json.Marshal(sinDescripcion)

	fmt.Println("  con descripción:", string(datos1))
	fmt.Println("  sin descripción:", string(datos2))
	fmt.Println("  nota: el segundo JSON no tiene la clave \"descripcion\" en absoluto")
	fmt.Println()
}

// -----------------------------------------------------------------------
// 4. JSON anidado: struct con slice de structs.
// -----------------------------------------------------------------------
func demoAnidado() {
	fmt.Println("--- 4. JSON anidado ---")

	pedido := Pedido{
		Cliente: "Ana Torres",
		Productos: []Producto{
			{Nombre: "Teclado", Precio: 450.50, Stock: 1},
			{Nombre: "Mouse", Precio: 220.00, Stock: 1},
		},
	}

	datos, err := json.Marshal(pedido)
	if err != nil {
		fmt.Println("  error:", err)
		return
	}
	fmt.Println("  JSON:", string(datos))
	fmt.Println()
}

// -----------------------------------------------------------------------
//  5. MarshalIndent: JSON legible para humanos (logs, debugging), en vez
//     del JSON compacto que normalmente mandas en una respuesta de API real
//     (donde cada byte extra es ancho de banda desperdiciado).
//
// -----------------------------------------------------------------------
func demoIndent() {
	fmt.Println("--- 5. json.MarshalIndent ---")

	pedido := Pedido{
		Cliente: "Luis Pérez",
		Productos: []Producto{
			{Nombre: "Monitor", Precio: 3200.00, Stock: 1},
		},
	}

	datos, err := json.MarshalIndent(pedido, "", "  ")
	if err != nil {
		fmt.Println("  error:", err)
		return
	}
	fmt.Println(string(datos))
}
