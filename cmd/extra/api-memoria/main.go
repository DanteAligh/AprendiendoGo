// Día 19: API REST en memoria — CRUD completo sobre un recurso ("tareas").
//
// Este archivo es REFERENCIA DE SINTAXIS, no la solución de los ejercicios
// (el recurso de los ejercicios es distinto: "tareas" con otros campos y
// reglas). Aquí modelamos "productos" para que veas el patrón completo:
// JSON (día 17) + HTTP (día 18) + almacenamiento en memoria protegido con
// sync.Mutex (día 15).
//
// Corre con:
//
//	go run ./cmd/extra/api-memoria
//
// Y pruébalo en otra terminal:
//
//	curl http://localhost:8082/productos
//	curl -X POST -d '{"nombre":"Teclado","precio":450.5,"stock":10}' http://localhost:8082/productos
//	curl http://localhost:8082/productos/1
//	curl -X PUT -d '{"nombre":"Teclado mecánico","precio":500,"stock":8}' http://localhost:8082/productos/1
//	curl -X DELETE http://localhost:8082/productos/1
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// Producto es el recurso que exponemos por HTTP.
type Producto struct {
	ID     int     `json:"id"`
	Nombre string  `json:"nombre"`
	Precio float64 `json:"precio"`
	Stock  int     `json:"stock"`
}

// almacen guarda los productos en memoria. Como varias peticiones HTTP se
// atienden en goroutines distintas y TODAS tocan este mismo map, protegemos
// cada acceso con un Mutex (día 15). Sin esto, dos POST simultáneos podrían
// pisarse el "siguiente id" y terminar con productos duplicados o
// sobrescritos.
type almacen struct {
	mu          sync.Mutex
	productos   map[int]Producto
	siguienteID int
}

func nuevoAlmacen() *almacen {
	return &almacen{
		productos: map[int]Producto{
			1: {ID: 1, Nombre: "Teclado", Precio: 450.50, Stock: 12},
			2: {ID: 2, Nombre: "Mouse", Precio: 220.00, Stock: 30},
		},
		siguienteID: 3,
	}
}

func (a *almacen) listar() []Producto {
	a.mu.Lock()
	defer a.mu.Unlock()

	lista := make([]Producto, 0, len(a.productos))
	for _, p := range a.productos {
		lista = append(lista, p)
	}
	return lista
}

func (a *almacen) obtener(id int) (Producto, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	p, ok := a.productos[id]
	return p, ok
}

func (a *almacen) crear(p Producto) Producto {
	a.mu.Lock()
	defer a.mu.Unlock()
	p.ID = a.siguienteID
	a.siguienteID++
	a.productos[p.ID] = p
	return p
}

func (a *almacen) actualizar(id int, p Producto) (Producto, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if _, ok := a.productos[id]; !ok {
		return Producto{}, false
	}
	p.ID = id
	a.productos[id] = p
	return p, true
}

func (a *almacen) eliminar(id int) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	if _, ok := a.productos[id]; !ok {
		return false
	}
	delete(a.productos, id)
	return true
}

func main() {
	al := nuevoAlmacen()

	mux := http.NewServeMux()
	// Ruta de colección: sin id, distinguimos la acción por método.
	mux.HandleFunc("/productos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			manejarListar(w, r, al)
		case http.MethodPost:
			manejarCrear(w, r, al)
		default:
			responderError(w, http.StatusMethodNotAllowed, "método no soportado en /productos")
		}
	})

	// Ruta de recurso individual: "/productos/" + id. El mux estándar no
	// separa automáticamente el {id} de la ruta (esto cambió en Go 1.22
	// con patrones como "GET /productos/{id}", que puedes explorar tú
	// mismo), así que aquí lo hacemos de forma explícita para que veas el
	// mecanismo de fondo.
	mux.HandleFunc("/productos/", func(w http.ResponseWriter, r *http.Request) {
		idTexto := strings.TrimPrefix(r.URL.Path, "/productos/")
		id, err := strconv.Atoi(idTexto)
		if err != nil {
			responderError(w, http.StatusBadRequest, "id inválido")
			return
		}

		switch r.Method {
		case http.MethodGet:
			manejarObtener(w, r, al, id)
		case http.MethodPut:
			manejarActualizar(w, r, al, id)
		case http.MethodDelete:
			manejarEliminar(w, r, al, id)
		default:
			responderError(w, http.StatusMethodNotAllowed, "método no soportado en /productos/{id}")
		}
	})

	fmt.Println("API REST en memoria escuchando en http://localhost:8082")
	if err := http.ListenAndServe(":8082", mux); err != nil {
		fmt.Println("error arrancando el servidor:", err)
	}
}

func manejarListar(w http.ResponseWriter, r *http.Request, al *almacen) {
	responderJSON(w, http.StatusOK, al.listar())
}

func manejarCrear(w http.ResponseWriter, r *http.Request, al *almacen) {
	var p Producto
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		responderError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	creado := al.crear(p)
	responderJSON(w, http.StatusCreated, creado)
}

func manejarObtener(w http.ResponseWriter, r *http.Request, al *almacen, id int) {
	p, ok := al.obtener(id)
	if !ok {
		responderError(w, http.StatusNotFound, "producto no encontrado")
		return
	}
	responderJSON(w, http.StatusOK, p)
}

func manejarActualizar(w http.ResponseWriter, r *http.Request, al *almacen, id int) {
	var p Producto
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		responderError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	actualizado, ok := al.actualizar(id, p)
	if !ok {
		responderError(w, http.StatusNotFound, "producto no encontrado")
		return
	}
	responderJSON(w, http.StatusOK, actualizado)
}

func manejarEliminar(w http.ResponseWriter, r *http.Request, al *almacen, id int) {
	if !al.eliminar(id) {
		responderError(w, http.StatusNotFound, "producto no encontrado")
		return
	}
	responderJSON(w, http.StatusOK, map[string]string{"mensaje": "producto eliminado"})
}

// responderJSON centraliza cómo respondemos JSON: mismo header, mismo
// código de estado, mismo Encode. Evita repetir estas tres líneas en cada
// handler.
func responderJSON(w http.ResponseWriter, codigo int, datos interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(codigo)
	json.NewEncoder(w).Encode(datos)
}

func responderError(w http.ResponseWriter, codigo int, mensaje string) {
	responderJSON(w, codigo, map[string]string{"error": mensaje})
}
