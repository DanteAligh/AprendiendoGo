package handler

import (
	"net/http"

	"aprendiendo-go/semana4/proyecto-final/internal/service"
	"aprendiendo-go/semana4/proyecto-final/pkg/respuesta"
)

// FacturaHandler agrupa los handlers HTTP relacionados a Factura.
type FacturaHandler struct {
	servicio service.FacturaService
}

func NewFacturaHandler(servicio service.FacturaService) *FacturaHandler {
	return &FacturaHandler{servicio: servicio}
}

// Listar debería responder GET /facturas.
func (h *FacturaHandler) Listar(w http.ResponseWriter, r *http.Request) {
	// TODO: implementar (igual patrón que ProductoHandler.Listar).
	respuesta.Error(w, http.StatusNotImplemented, "TODO: implementar Listar facturas")
}

// Crear debería responder POST /facturas. Este es el endpoint más
// importante del proyecto integrador: acá es donde tu API expone la
// lógica de negocio que escribiste en facturaServicio.Crear (validación de
// stock, congelar precios, descuento de inventario).
//
// Pasos sugeridos:
//  1. Decodificar el body JSON a un models.Factura (con su slice de Items).
//  2. Llamar a h.servicio.Crear(factura).
//  3. Distinguir en la respuesta al menos estos tres casos de error:
//     - producto no existe en algún item -> 400 o 404 (justificá tu elección)
//     - stock insuficiente -> 409 Conflict (es el código semánticamente
//     correcto para "tu petición es válida pero choca con el estado
//     actual del servidor")
//     - error interno inesperado -> 500
//  4. Si todo sale bien, responder 201 con la factura creada (confirmada,
//     con los precios ya congelados).
func (h *FacturaHandler) Crear(w http.ResponseWriter, r *http.Request) {
	// TODO: implementar (ver pasos sugeridos arriba). Vas a necesitar
	// errors.Is contra service.ErrProductoNoExiste y
	// service.ErrStockInsuficiente para distinguir los casos de error.
	respuesta.Error(w, http.StatusNotImplemented, "TODO: implementar Crear factura")
}

// ObtenerPorID debería responder GET /facturas/{id}.
func (h *FacturaHandler) ObtenerPorID(w http.ResponseWriter, r *http.Request) {
	// TODO: implementar (mismo patrón que ProductoHandler.ObtenerPorID,
	// usando r.PathValue("id")).
	respuesta.Error(w, http.StatusNotImplemented, "TODO: implementar ObtenerPorID factura")
}
