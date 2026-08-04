// Día 1 — Ejemplo: estructura de un programa, variables, tipos, constantes
// y conversión de tipos con strconv.
//
// Caso de uso de este ejemplo: un "recibo simple" de una compra, para
// mostrar variables de distintos tipos trabajando juntas. (Los ejercicios
// de ejercicios.md usan otros casos: círculo, edad, moneda, etc. — este
// ejemplo NO los resuelve, solo enseña la sintaxis).
package main

// "fmt" es el paquete estándar para entrada/salida formateada (imprimir,
// leer, formatear strings). "strconv" sirve para convertir entre strings
// y tipos numéricos.
import (
	"fmt"
	"strconv"
)

// func main() es el punto de entrada obligatorio de cualquier programa Go
// ejecutable. Sin esta función (dentro de package main), "go run" no sabría
// por dónde empezar.
func main() {
	// Declaración corta con ":=" — Go infiere el tipo automáticamente.
	// Aquí "producto" es string y "cantidad" es int, sin que lo escribamos.
	producto := "Teclado mecánico"
	cantidad := 2

	// Declaración explícita con "var tipo" — útil cuando quieres dejar el
	// tipo bien claro, o cuando aún no tienes un valor inicial.
	var precioUnitario float64 = 89.5

	// Una constante: su valor NO puede cambiar después de esta línea.
	// La usamos para el porcentaje de impuesto, que en este programa es fijo.
	const IMPUESTO = 0.19 // 19%

	// Tipado fuerte en acción: no podemos multiplicar un int por un float64
	// directamente. Convertimos "cantidad" (int) a float64 explícitamente
	// con float64(...) para poder operar con precioUnitario.
	subtotal := float64(cantidad) * precioUnitario
	impuestoCalculado := subtotal * IMPUESTO
	total := subtotal + impuestoCalculado

	// %s para strings, %d para enteros, %.2f para decimales con 2 posiciones.
	fmt.Printf("Producto: %s\n", producto)
	fmt.Printf("Cantidad: %d\n", cantidad)
	fmt.Printf("Subtotal: %.2f\n", subtotal)
	fmt.Printf("Impuesto (19%%): %.2f\n", impuestoCalculado) // %% imprime un % literal
	fmt.Printf("Total: %.2f\n", total)

	// --- Valor cero de las variables ---
	// Si declaras una variable con "var" sin inicializarla, Go le asigna
	// automáticamente su "valor cero": 0 para números, "" para strings,
	// false para bools. Esto evita el problema de variables "sin definir"
	// que existe en otros lenguajes.
	var codigoDescuento string
	var descuentoAplicado bool
	fmt.Println("Código de descuento (valor cero):", codigoDescuento)
	fmt.Println("¿Descuento aplicado? (valor cero):", descuentoAplicado)

	// --- Conversión string -> número con strconv ---
	// Imaginemos que el número de factura llegó como texto (por ejemplo,
	// desde un formulario web). strconv.Atoi devuelve DOS valores: el
	// número convertido y un posible error. Por ahora guardamos el error
	// en "_" (el "blank identifier") porque todavía no vamos a manejarlo -
	// eso lo vemos más adelante en la semana de manejo de errores.
	numeroFacturaTexto := "1024"
	numeroFactura, _ := strconv.Atoi(numeroFacturaTexto)
	numeroFacturaSiguiente := numeroFactura + 1

	fmt.Println("Número de factura actual:", numeroFactura)
	fmt.Println("Número de la próxima factura:", numeroFacturaSiguiente)

	// --- Conversión número -> string con strconv ---
	// strconv.Itoa hace el camino inverso: de int a string.
	totalComoTexto := strconv.Itoa(int(total))
	fmt.Println("El total, como texto plano: " + totalComoTexto)
}
