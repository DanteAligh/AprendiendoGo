// Día 26 — Modelado de un mini-ERP: Producto, Cliente, Factura
//
// Cómo correrlo (solo librería estándar):
//
//	go run ./cmd/extra/modelado-erp
//
// ---------------------------------------------------------------------------
// ¿QUÉ ES "MODELAR" UN DOMINIO?
//
// Antes de escribir un solo endpoint HTTP, todo backend serio empieza por
// una pregunta: ¿cuáles son las ENTIDADES del negocio y cómo se relacionan
// entre sí? Esto es el "modelo de dominio". En un ERP, típicamente hay:
//
//   - Producto:  algo que se vende o se usa (tiene precio, tiene stock).
//   - Cliente:   quién compra.
//   - Factura:   un documento que registra una venta a un Cliente,
//     compuesta por varios ItemFactura.
//   - ItemFactura: una línea dentro de una factura ("2 unidades del
//     producto X a $10 cada una"). Fijate que ItemFactura NO
//     contiene una copia del Producto entero — contiene su ID.
//     Esto es clave: relacionamos entidades por ID, igual que
//     harías con claves foráneas (foreign keys) en una base de
//     datos relacional. Guardar una copia completa del producto
//     dentro de cada item sería duplicar datos que después se
//     desincronizan (¿qué pasa si el precio del producto
//     cambia? ¿la factura vieja debería reflejar el precio
//     nuevo o el que tenía al momento de la venta? modelar bien
//     esto es una decisión de negocio real).
//
// Este archivo se queda a propósito en el nivel de "structs + funciones",
// SIN HTTP y SIN base de datos — la idea es que primero domines cómo se
// relacionan los datos y qué reglas de negocio aplican, antes de meter la
// complejidad de una API completa (eso es exactamente lo que vas a construir
// vos mismo en proyecto-final/).
package main

import (
	"errors"
	"fmt"
)

// Producto representa algo que el ERP vende o controla en inventario.
type Producto struct {
	ID     int
	Nombre string
	Precio float64 // precio unitario
	Stock  int     // unidades disponibles
}

// Cliente representa a quién se le factura.
type Cliente struct {
	ID     int
	Nombre string
	Email  string
}

// ItemFactura es UNA LÍNEA de una factura: "esta cantidad de este
// producto". Fijate que solo guarda ProductoID (la relación por ID que
// mencionamos arriba), no un Producto completo.
type ItemFactura struct {
	ProductoID int
	Cantidad   int
	// PrecioUnitario se copia AL MOMENTO de facturar, no se recalcula
	// después. Esto es intencional: si mañana el producto sube de precio,
	// las facturas viejas no deben cambiar retroactivamente. Es una regla
	// de negocio real de cualquier sistema de facturación.
	PrecioUnitario float64
}

// Subtotal calcula cuánto suma esta línea: cantidad * precio unitario.
// Es un método en ItemFactura porque es un cálculo que depende SOLO de los
// datos de ItemFactura — no necesita mirar la Factura completa ni la base
// de datos. Mantener los cálculos lo más "locales" posible hace que sean
// triviales de testear (semana 3: table-driven tests) de forma aislada.
func (i ItemFactura) Subtotal() float64 {
	return float64(i.Cantidad) * i.PrecioUnitario
}

// Factura agrupa varios ItemFactura para un Cliente determinado.
type Factura struct {
	ID         int
	ClienteID  int
	Items      []ItemFactura
	Confirmada bool // false = borrador, true = ya se descontó el stock
}

// Total suma el subtotal de todos los items. Este es exactamente el tipo
// de lógica de negocio que en la arquitectura de capas (día 22) vive en la
// capa de "service", NUNCA en el handler HTTP: el handler solo debería
// llamar a factura.Total(), no calcularlo él mismo.
func (f Factura) Total() float64 {
	total := 0.0
	for _, item := range f.Items {
		total += item.Subtotal()
	}
	return total
}

// CantidadItems es un pequeño helper que cuenta cuántas líneas tiene la
// factura (no cuántas unidades — eso sería otra función). Sirve para
// mostrar que podés tener varias funciones de negocio pequeñas y con
// nombres claros, en vez de una sola función gigante que hace de todo.
func (f Factura) CantidadItems() int {
	return len(f.Items)
}

// ErrStockInsuficiente se declara como variable de paquete (siguiendo el
// patrón de manejo de errores que ya viste en la semana 2 con errors.Is)
// para que quien llame a esta función pueda distinguir ESTE error
// específico de otros errores posibles, sin depender de comparar strings.
var ErrStockInsuficiente = errors.New("stock insuficiente")

// ErrProductoNoExiste indica que se referenció un ProductoID que no está
// en el catálogo — un caso muy común al validar facturas.
var ErrProductoNoExiste = errors.New("producto no existe")

// ValidarStock recorre los items de una factura y verifica, contra un
// catálogo de productos (acá representado simplemente como
// map[int]Producto), que haya stock suficiente para CADA línea. Fijate que
// esta función NO modifica nada — solo valida. Separar "validar" de
// "aplicar el descuento de stock" es otra decisión de diseño deliberada:
// nunca querés descontar stock a medias si la línea 3 de 5 falla.
func ValidarStock(factura Factura, catalogo map[int]Producto) error {
	for _, item := range factura.Items {
		producto, existe := catalogo[item.ProductoID]
		if !existe {
			return fmt.Errorf("item con producto ID %d: %w", item.ProductoID, ErrProductoNoExiste)
		}
		if producto.Stock < item.Cantidad {
			return fmt.Errorf("producto %q (stock disponible: %d, pedido: %d): %w",
				producto.Nombre, producto.Stock, item.Cantidad, ErrStockInsuficiente)
		}
	}
	return nil
}

// AplicarDescuentoStock asume que ValidarStock ya pasó sin error, y recién
// ahí descuenta las cantidades del catálogo. Nota: catalogo se recibe como
// map[int]Producto (no un puntero al mapa) porque los mapas en Go ya son
// tipos de referencia — modificar catalogo[id] adentro de esta función SÍ
// afecta al mapa original que le pasó quien llamó. Esto es un detalle sutil
// de Go que vale la pena recordar de la semana 1/2.
func AplicarDescuentoStock(factura Factura, catalogo map[int]Producto) {
	for _, item := range factura.Items {
		producto := catalogo[item.ProductoID]
		producto.Stock -= item.Cantidad
		catalogo[item.ProductoID] = producto // los structs son por valor: hay que reasignar
	}
}

func main() {
	// Catálogo de productos, simulando lo que en proyecto-final/ vendría de
	// una base de datos o de un repositorio en memoria.
	catalogo := map[int]Producto{
		1: {ID: 1, Nombre: "Teclado mecánico", Precio: 45.00, Stock: 10},
		2: {ID: 2, Nombre: "Mouse inalámbrico", Precio: 25.50, Stock: 3},
		3: {ID: 3, Nombre: "Monitor 24\"", Precio: 180.00, Stock: 0},
	}

	cliente := Cliente{ID: 1, Nombre: "Comercial Andina S.A.", Email: "compras@andina.com"}

	fmt.Println("=== Caso 1: factura válida ===")
	facturaValida := Factura{
		ID:        100,
		ClienteID: cliente.ID,
		Items: []ItemFactura{
			{ProductoID: 1, Cantidad: 2, PrecioUnitario: catalogo[1].Precio},
			{ProductoID: 2, Cantidad: 1, PrecioUnitario: catalogo[2].Precio},
		},
	}

	if err := ValidarStock(facturaValida, catalogo); err != nil {
		fmt.Println("Error inesperado:", err)
	} else {
		AplicarDescuentoStock(facturaValida, catalogo)
		facturaValida.Confirmada = true
		fmt.Printf("Factura #%d confirmada para %s\n", facturaValida.ID, cliente.Nombre)
		fmt.Printf("  Items: %d | Total: $%.2f\n", facturaValida.CantidadItems(), facturaValida.Total())
		fmt.Printf("  Stock de teclados después de facturar: %d\n", catalogo[1].Stock)
		fmt.Printf("  Stock de mouses después de facturar: %d\n", catalogo[2].Stock)
	}

	fmt.Println("\n=== Caso 2: factura con stock insuficiente (monitor, stock=0) ===")
	facturaSinStock := Factura{
		ID:        101,
		ClienteID: cliente.ID,
		Items: []ItemFactura{
			{ProductoID: 3, Cantidad: 1, PrecioUnitario: catalogo[3].Precio},
		},
	}
	err := ValidarStock(facturaSinStock, catalogo)
	if errors.Is(err, ErrStockInsuficiente) {
		fmt.Println("Como se esperaba, se rechazó por stock insuficiente:", err)
	} else if err != nil {
		fmt.Println("Error inesperado:", err)
	}

	fmt.Println("\n=== Caso 3: factura con un producto que no existe ===")
	facturaProductoInexistente := Factura{
		ID:        102,
		ClienteID: cliente.ID,
		Items: []ItemFactura{
			{ProductoID: 999, Cantidad: 1, PrecioUnitario: 10.0},
		},
	}
	err = ValidarStock(facturaProductoInexistente, catalogo)
	if errors.Is(err, ErrProductoNoExiste) {
		fmt.Println("Como se esperaba, se rechazó por producto inexistente:", err)
	} else if err != nil {
		fmt.Println("Error inesperado:", err)
	}
}
