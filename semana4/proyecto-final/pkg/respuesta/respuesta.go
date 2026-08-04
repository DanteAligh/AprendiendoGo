// Package respuesta vive en pkg/ (no en internal/) porque es un helper
// genérico de propósito general para responder JSON por HTTP: no tiene
// nada específico de este ERP, así que en teoría podrías copiarlo a
// cualquier otro proyecto Go tuyo (repasá la diferencia entre pkg/ e
// internal/ del día 22). Está completamente implementado como andamiaje:
// tu trabajo en los handlers es USARLO, no reescribirlo.
package respuesta

import (
	"encoding/json"
	"net/http"
)

// errorJSON es la forma estándar en la que este API devuelve errores:
// {"error": "mensaje legible"}. Tener un formato consistente en TODA la
// API hace que quien la consuma (tu propio frontend, un curl, otro
// servicio) sepa siempre dónde buscar el mensaje de error.
type errorJSON struct {
	Error string `json:"error"`
}

// JSON escribe cualquier valor como JSON en la respuesta HTTP, con el
// código de estado indicado y el header Content-Type correcto. Usalo así
// desde un handler:
//
//	respuesta.JSON(w, http.StatusOK, productos)
//	respuesta.JSON(w, http.StatusCreated, productoCreado)
func JSON(w http.ResponseWriter, codigoEstado int, datos interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(codigoEstado)
	if err := json.NewEncoder(w).Encode(datos); err != nil {
		// Si falló CODIFICAR la respuesta (raro, pero posible), ya
		// escribimos el header y el código de estado, así que lo único
		// que podemos hacer es loguearlo — no podemos "deshacer" la
		// respuesta a esta altura.
		return
	}
}

// Error escribe una respuesta de error consistente: {"error": "..."} con
// el código de estado HTTP indicado. Usalo así:
//
//	respuesta.Error(w, http.StatusNotFound, "producto no encontrado")
//	respuesta.Error(w, http.StatusBadRequest, "el campo nombre es obligatorio")
func Error(w http.ResponseWriter, codigoEstado int, mensaje string) {
	JSON(w, codigoEstado, errorJSON{Error: mensaje})
}
