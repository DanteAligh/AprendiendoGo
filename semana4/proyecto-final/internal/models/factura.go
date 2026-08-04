package models

// ItemFactura es una línea dentro de una Factura: "esta cantidad de este
// producto, a este precio". Fijate que se relaciona con Producto por ID
// (ProductoID), NO incrustando un Producto completo — el mismo patrón que
// viste en el ejemplo del día 26. PrecioUnitario se guarda en el momento
// de facturar (no se recalcula después contra el precio actual del
// producto): eso es una decisión de negocio deliberada.
type ItemFactura struct {
	ProductoID     int     `json:"producto_id"`
	Cantidad       int     `json:"cantidad"`
	PrecioUnitario float64 `json:"precio_unitario"`
}

// Factura es el segundo recurso principal de este ERP, y el que muestra la
// RELACIÓN entre recursos: una Factura agrupa varios ItemFactura, cada uno
// apuntando a un Producto existente.
//
// Confirmada distingue una factura recién creada (false) de una donde ya
// se validó stock y se descontó del inventario (true) — la misma idea que
// viste en el ejemplo del día 26.
type Factura struct {
	ID            int           `json:"id"`
	ClienteNombre string        `json:"cliente_nombre"`
	Items         []ItemFactura `json:"items"`
	Confirmada    bool          `json:"confirmada"`
}

// TODO (te toca a vos):
// El cálculo del total de una factura (sumar cantidad * precio_unitario de
// cada item) es lógica de negocio, no un detalle del modelo. Ejercicios.md
// te va a pedir implementarlo en internal/service, tal como lo practicaste
// (sin la solución) en el ejemplo del día 26. Resistí la tentación de
// agregar un método Total() acá mismo antes de leer la guía de los días
// 27-28: pensar POR QUÉ va en el service y no en el modelo es parte del
// ejercicio.
