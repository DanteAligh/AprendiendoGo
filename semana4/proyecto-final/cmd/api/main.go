// Comando api: el punto de entrada de este ERP mínimo.
//
// Este archivo es ANDAMIAJE YA TERMINADO: arma (conecta) las 4 capas —
// repository, service, handler y middleware — y levanta el servidor HTTP.
// Es "cableado" (wiring), no lógica de negocio, así que ya viene resuelto
// para que puedas concentrarte en internal/repository, internal/service e
// internal/handler, que es donde vas a escribir código de verdad siguiendo
// ejercicios.md (días 27-28).
//
// Ahora mismo, el servidor COMPILA Y CORRE, pero cada endpoint responde
// 501 Not Implemented, porque los métodos de internal/handler todavía
// están en TODO. A medida que vayas implementando capa por capa
// (repository -> service -> handler, ese orden, recurso por recurso), vas
// a ir viendo cómo cada endpoint deja de devolver 501 y empieza a
// funcionar de verdad.
package main

import (
	"log"
	"net/http"
	"os"

	"aprendiendo-go/semana4/proyecto-final/internal/handler"
	"aprendiendo-go/semana4/proyecto-final/internal/middleware"
	"aprendiendo-go/semana4/proyecto-final/internal/repository"
	"aprendiendo-go/semana4/proyecto-final/internal/service"
)

func main() {
	// --- Capa de repositorio (persistencia) ---
	// Ahora mismo solo existe la implementación en memoria. El día que
	// quieras agregar una implementación con SQLite (repasá el día 20 de
	// la semana 3, y el ejemplo de bcrypt/JWT del día 24 para el patrón de
	// go get + go.mod de una dependencia externa), solo tendrías que:
	//   1. Escribir un nuevo struct en internal/repository que también
	//      implemente ProductoRepository/FacturaRepository.
	//   2. Cambiar estas dos líneas de acá abajo para instanciar el nuevo
	//      struct en vez del de memoria.
	// Ninguna otra capa (service, handler) debería necesitar cambios.
	repoProductos := repository.NewRepositorioProductoMemoria()
	repoFacturas := repository.NewRepositorioFacturaMemoria()

	// --- Capa de servicio (reglas de negocio) ---
	servicioProductos := service.NewProductoService(repoProductos)
	servicioFacturas := service.NewFacturaService(repoFacturas, repoProductos)

	// --- Capa de handlers (HTTP) ---
	handlerProductos := handler.NewProductoHandler(servicioProductos)
	handlerFacturas := handler.NewFacturaHandler(servicioFacturas)

	// --- Rutas ---
	// Usamos el enrutador estándar de Go (http.ServeMux), que desde Go 1.22
	// soporta patrones con método HTTP y variables de path (ej. "{id}") sin
	// necesidad de ninguna librería externa como gorilla/mux o chi. Por
	// eso este proyecto entero no requiere `go get` de nada.
	mux := http.NewServeMux()

	mux.HandleFunc("GET /productos", handlerProductos.Listar)
	mux.HandleFunc("POST /productos", handlerProductos.Crear)
	mux.HandleFunc("GET /productos/{id}", handlerProductos.ObtenerPorID)
	mux.HandleFunc("PUT /productos/{id}", handlerProductos.Actualizar)
	mux.HandleFunc("DELETE /productos/{id}", handlerProductos.Eliminar)

	mux.HandleFunc("GET /facturas", handlerFacturas.Listar)
	mux.HandleFunc("POST /facturas", handlerFacturas.Crear)
	mux.HandleFunc("GET /facturas/{id}", handlerFacturas.ObtenerPorID)

	// Endpoint de salud: útil para verificar rápido que el servidor está
	// vivo (y, en un despliegue real, para que un balanceador de carga o
	// Docker/Kubernetes sepan si el proceso está sano). Este SÍ está
	// completamente implementado porque no tiene ninguna lógica de negocio.
	mux.HandleFunc("GET /salud", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"estado":"ok"}`))
	})

	// --- Middlewares ---
	// Repasá el día 23: el orden importa. RecuperarPanico va primero (más
	// "externo") para poder atrapar pánicos de TODO lo demás, incluido el
	// propio logging si hiciera falta.
	handlerFinal := middleware.Encadenar(mux,
		middleware.RecuperarPanico,
		middleware.Logueo,
		middleware.CORS,
	)

	// --- Configuración desde variables de entorno (día 23) ---
	puerto := os.Getenv("PUERTO")
	if puerto == "" {
		puerto = "8080"
	}

	direccion := ":" + puerto
	log.Printf("erp-api escuchando en http://localhost%s", direccion)
	log.Printf("probá: curl http://localhost%s/salud", direccion)
	log.Printf("probá: curl http://localhost%s/productos", direccion)

	if err := http.ListenAndServe(direccion, handlerFinal); err != nil {
		log.Fatal(err)
	}
}
