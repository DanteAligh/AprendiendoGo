// Package handler es la capa que SÍ sabe de HTTP: decodifica JSON de
// entrada, llama a la capa de service, y codifica la respuesta JSON con el
// código de estado adecuado (usando pkg/respuesta). No debería contener
// reglas de negocio (eso vive en internal/service) ni saber cómo se
// guardan los datos (internal/repository).
//
// Cada método de acá abajo está registrado como una ruta en
// cmd/api/main.go. Ahora mismo todos responden 501 Not Implemented: ese es
// tu punto de partida. Guiate por ejercicios.md (días 27-28) para
// implementarlos en el orden sugerido.
package handler

import (
	"net/http"

	"aprendiendo-go/semana4/proyecto-final/internal/service"
	"aprendiendo-go/semana4/proyecto-final/pkg/respuesta"
)

// ProductoHandler agrupa los handlers HTTP relacionados a Producto. Guarda
// una referencia a la INTERFAZ service.ProductoService (no a un struct
// concreto), inyectada desde cmd/api/main.go.
type ProductoHandler struct {
	servicio service.ProductoService
}

// NewProductoHandler construye el handler inyectando su dependencia.
func NewProductoHandler(servicio service.ProductoService) *ProductoHandler {
	return &ProductoHandler{servicio: servicio}
}

// Listar debería responder GET /productos con la lista completa.
//
// Pasos sugeridos:
//  1. Llamar a h.servicio.Listar().
//  2. Si hay error, responder con respuesta.Error(w, http.StatusInternalServerError, err.Error()).
//  3. Si no, responder con respuesta.JSON(w, http.StatusOK, productos).
func (h *ProductoHandler) Listar(w http.ResponseWriter, r *http.Request) {
	// TODO: implementar (ver pasos sugeridos arriba).
	respuesta.Error(w, http.StatusNotImplemented, "TODO: implementar Listar productos")
}

// Crear debería responder POST /productos.
//
// Pasos sugeridos:
//  1. Decodificar el body JSON a un models.Producto con json.NewDecoder(r.Body).Decode(&producto).
//     Si el JSON es inválido, responder http.StatusBadRequest.
//  2. Llamar a h.servicio.Crear(producto).
//  3. Si el service devuelve un error de validación, responder
//     http.StatusBadRequest o http.StatusUnprocessableEntity (pensá cuál
//     tiene más sentido y por qué).
//  4. Si todo salió bien, responder http.StatusCreated (201) con el
//     producto ya creado (con su ID asignado).
func (h *ProductoHandler) Crear(w http.ResponseWriter, r *http.Request) {
	// TODO: implementar (ver pasos sugeridos arriba).
	respuesta.Error(w, http.StatusNotImplemented, "TODO: implementar Crear producto")
}

// ObtenerPorID debería responder GET /productos/{id}.
//
// Pista: con el enrutador estándar de Go 1.22+ (el que ya está registrado
// en cmd/api/main.go con patrones como "GET /productos/{id}"), podés leer
// el valor de {id} desde el propio request con:
//
//	idTexto := r.PathValue("id")
//
// Vas a necesitar convertir idTexto (string) a int con strconv.Atoi, y
// manejar el error si no es un número válido (responder 400).
func (h *ProductoHandler) ObtenerPorID(w http.ResponseWriter, r *http.Request) {
	// TODO: implementar.
	// Recordá distinguir el caso "no encontrado" (404) del caso "error
	// interno" (500) — vas a necesitar errors.Is contra el error que
	// definiste en la capa de repository/service para eso.
	respuesta.Error(w, http.StatusNotImplemented, "TODO: implementar ObtenerPorID producto")
}

// Actualizar debería responder PUT /productos/{id}.
func (h *ProductoHandler) Actualizar(w http.ResponseWriter, r *http.Request) {
	// TODO: implementar. Combina lo de Crear (decodificar body) con lo de
	// ObtenerPorID (leer {id} del path). Pensá: ¿qué pasa si el id del
	// path no coincide con el id del body? ¿cuál debería ganar?
	respuesta.Error(w, http.StatusNotImplemented, "TODO: implementar Actualizar producto")
}

// Eliminar debería responder DELETE /productos/{id}.
//
// Convención común en APIs REST: si el borrado fue exitoso, responder
// http.StatusNoContent (204) SIN cuerpo (no llames a respuesta.JSON en ese
// caso, usá w.WriteHeader(http.StatusNoContent) directamente).
func (h *ProductoHandler) Eliminar(w http.ResponseWriter, r *http.Request) {
	// TODO: implementar.
	respuesta.Error(w, http.StatusNotImplemented, "TODO: implementar Eliminar producto")
}
