// Package middleware contiene los middlewares HTTP transversales del
// proyecto (logging, recuperación de pánicos, CORS), tal como los
// practicaste sin arquitectura de capas en ejemplos/dia23_middleware_env.go.
// Este archivo está COMPLETAMENTE IMPLEMENTADO como andamiaje: es
// infraestructura reutilizable, no lógica de negocio del ERP, así que no
// forma parte de lo que tenés que resolver vos. Tu trabajo es en
// internal/handler, internal/service e internal/repository.
package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logueo registra método, ruta, código de estado y duración de cada
// request. Para poder capturar el código de estado (net/http no lo expone
// directamente) envolvemos http.ResponseWriter en un wrapper mínimo que
// intercepta la llamada a WriteHeader.
func Logueo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inicio := time.Now()
		envoltorio := &escritorConEstado{ResponseWriter: w, estado: http.StatusOK}
		next.ServeHTTP(envoltorio, r)
		log.Printf("%s %s -> %d (%v)", r.Method, r.URL.Path, envoltorio.estado, time.Since(inicio))
	})
}

// escritorConEstado envuelve http.ResponseWriter para recordar qué código
// de estado se envió, sin cambiar el comportamiento real de la escritura.
type escritorConEstado struct {
	http.ResponseWriter
	estado int
}

func (e *escritorConEstado) WriteHeader(codigo int) {
	e.estado = codigo
	e.ResponseWriter.WriteHeader(codigo)
}

// RecuperarPanico evita que un panic() en cualquier handler tumbe todo el
// proceso: lo atrapa, responde 500 a ESA petición puntual, y el servidor
// sigue en pie para el resto de las peticiones. Repasá el comentario
// equivalente en ejemplos/dia23_middleware_env.go si querés el detalle de
// por qué esto es crítico en un servidor concurrente.
func RecuperarPanico(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC recuperado en %s %s: %v", r.Method, r.URL.Path, err)
				http.Error(w, `{"error":"error interno del servidor"}`, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// CORS habilita que un frontend en otro origen (ej. http://localhost:3000)
// pueda llamar a esta API durante desarrollo.
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Encadenar aplica una lista de middlewares a un handler base, en el orden
// en que aparecen en la lista (el primero de la lista es el más "externo").
func Encadenar(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}
