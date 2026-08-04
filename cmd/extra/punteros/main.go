// Día 10: Punteros
//
// Los ejercicios trabajan con intercambiar(), BuscarPrimero(), Config y
// ContarNil(). Aquí usamos otro caso: un "carrito de compras" simulado y una
// búsqueda de descuentos, para no repetir el mismo ejemplo que resolverías
// tú mismo.
package main

import "fmt"

// Producto representa un artículo con su precio. Es un struct pequeño mas
// lo usamos para mostrar cómo un puntero permite modificar el original.
type Producto struct {
	Nombre string
	Precio float64
}

// aplicarDescuento recibe un PUNTERO a Producto y modifica su precio
// directamente. Si recibiera un Producto por valor, el descuento se
// aplicaría solo a una copia y se perdería al salir de la función — el
// mismo problema que vimos ayer con receptores por valor.
func aplicarDescuento(p *Producto, porcentaje float64) {
	p.Precio = p.Precio - (p.Precio * porcentaje / 100)
}

// buscarMasCaro recorre un slice de productos y devuelve un PUNTERO al más
// caro, o nil si el slice está vacío. Devolver un puntero aquí evita copiar
// el struct completo, y nil nos permite representar limpiamente "no hay
// nada que devolver" (muy común en funciones que buscan algo en un backend:
// "no encontrado" en vez de un valor inventado).
func buscarMasCaro(productos []Producto) *Producto {
	if len(productos) == 0 {
		return nil
	}
	masCaro := &productos[0]
	for i := 1; i < len(productos); i++ {
		if productos[i].Precio > masCaro.Precio {
			masCaro = &productos[i]
		}
	}
	return masCaro
}

func main() {
	// --- & y * básicos ---
	precioBase := 100.0
	p := &precioBase // p es *float64, guarda la dirección de precioBase
	fmt.Println("Valor de precioBase:", precioBase)
	fmt.Println("Dirección guardada en p:", p)
	fmt.Println("Valor desreferenciado (*p):", *p)

	*p = 120.0 // modificamos precioBase indirectamente
	fmt.Println("Después de *p = 120.0, precioBase vale:", precioBase)

	// --- Puntero a struct que modifica el original ---
	producto := Producto{Nombre: "Teclado mecánico", Precio: 800.0}
	fmt.Printf("\nAntes del descuento: %+v\n", producto)
	aplicarDescuento(&producto, 15) // pasamos la dirección con &
	fmt.Printf("Después de 15%% de descuento: %+v\n", producto)

	// --- Puntero devuelto por una función, y el chequeo obligatorio de nil ---
	catalogo := []Producto{
		{Nombre: "Mouse", Precio: 250.0},
		{Nombre: "Monitor", Precio: 3200.0},
		{Nombre: "Cable USB", Precio: 50.0},
	}

	masCaro := buscarMasCaro(catalogo)
	// SIEMPRE hay que verificar nil antes de desreferenciar un puntero que
	// puede no apuntar a nada. Si nos saltamos esto y masCaro fuera nil,
	// *masCaro provocaría un panic en tiempo de ejecución.
	if masCaro != nil {
		fmt.Printf("\nEl producto más caro es: %s ($%.2f)\n", masCaro.Nombre, masCaro.Precio)
	} else {
		fmt.Println("\nEl catálogo está vacío, no hay producto más caro.")
	}

	// Probamos también el caso vacío para ver el otro camino:
	vacio := []Producto{}
	resultado := buscarMasCaro(vacio)
	if resultado == nil {
		fmt.Println("Catálogo vacío: buscarMasCaro devolvió nil, como se esperaba.")
	}
}
