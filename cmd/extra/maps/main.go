// Día 7 — Ejemplo: maps (CRUD, range, delete, comma-ok idiom) + un repaso
// que combina variables, condicionales, ciclos, funciones y slices.
//
// Caso de uso de este ejemplo: un mini "catálogo de precios" de una tienda
// con reglas de descuento. El mini-reto integrador de ejercicios.md (sistema
// de calificaciones) es un caso DISTINTO — este ejemplo no lo resuelve.
package main

import "fmt"

// Una función normal, como las del día 5, que usa un switch sin expresión
// (como en el día 3) para decidir el porcentaje de descuento según la
// categoría de un producto.
func descuentoPorCategoria(categoria string) float64 {
	switch categoria {
	case "electronica":
		return 10.0
	case "ropa":
		return 20.0
	case "alimentos":
		return 0.0
	default:
		return 5.0
	}
}

func main() {
	// --- Declarar un map ---
	// map[string]float64: la llave es el nombre del producto, el valor es
	// su precio.
	precios := map[string]float64{
		"laptop":   3200000,
		"camiseta": 45000,
		"arroz":    6000,
	}
	// Y un segundo map que asocia cada producto con su categoría, para
	// poder calcular el descuento con la función de arriba.
	categorias := map[string]string{
		"laptop":   "electronica",
		"camiseta": "ropa",
		"arroz":    "alimentos",
	}

	// --- Crear / actualizar (CRUD) ---
	precios["mouse"] = 55000 // crear una llave nueva
	categorias["mouse"] = "electronica"
	precios["camiseta"] = 42000 // actualizar una llave existente (rebaja de precio)

	// --- range sobre un map ---
	// Recuerda: el ORDEN no está garantizado. Cada ejecución puede variar.
	fmt.Println("--- Catálogo con descuento aplicado ---")
	var totalConDescuento float64
	var productosProcesados []string // slice para ir guardando nombres, como en el día 6
	for producto, precio := range precios {
		descuento := descuentoPorCategoria(categorias[producto])
		precioFinal := precio - (precio * descuento / 100)
		totalConDescuento += precioFinal
		productosProcesados = append(productosProcesados, producto)
		fmt.Printf("%-10s | precio: %.0f | descuento: %.0f%% | precio final: %.0f\n",
			producto, precio, descuento, precioFinal)
	}
	fmt.Printf("\nProductos procesados: %v\n", productosProcesados)
	fmt.Printf("Total del catálogo (con descuentos): %.0f\n", totalConDescuento)

	// --- comma-ok idiom ---
	// Verificamos si "audifonos" existe en el map de precios. Si accediéramos
	// directamente con precios["audifonos"] sin comprobar, Go nos daría 0.0
	// (el valor cero de float64), lo cual sería indistinguible de un
	// producto que de verdad cuesta 0. El segundo valor booleano resuelve
	// esa ambigüedad.
	fmt.Println("\n--- comma-ok idiom ---")
	precio, existe := precios["audifonos"]
	if !existe {
		fmt.Println("'audifonos' no está en el catálogo")
	} else {
		fmt.Println("Precio de audífonos:", precio)
	}

	precio, existe = precios["laptop"]
	if existe {
		// Usamos %.0f en vez de Println directo: los float64 grandes se
		// imprimen en notación científica (3.2e+06) con el formato %v por
		// defecto, y %.0f nos da el número completo, más legible aquí.
		fmt.Printf("Precio de laptop confirmado: %.0f\n", precio)
	}

	// --- delete ---
	fmt.Println("\n--- delete ---")
	fmt.Println("Antes de eliminar 'arroz', ¿existe?:", func() bool {
		_, ok := precios["arroz"]
		return ok
	}())
	delete(precios, "arroz")
	_, existeArroz := precios["arroz"]
	fmt.Println("Después de eliminar 'arroz', ¿existe?:", existeArroz)
}
