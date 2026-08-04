// Día 18: net/http básico — tu primer servidor HTTP desde cero.
//
// Este archivo es REFERENCIA DE SINTAXIS, no la solución de los ejercicios.
// Arranca un servidor real en el puerto 8081 y se queda escuchando (no
// termina solo). Corre con:
//
//	go run ./cmd/extra/http-servidor
//
// Y en otra terminal, pruébalo con:
//
//	curl http://localhost:8081/
//	curl "http://localhost:8081/saludo?nombre=Carlos"
//	curl -X POST -d "hola mundo" http://localhost:8081/eco
//	curl http://localhost:8081/saludo-json
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Saludo es la struct que devolvemos como JSON en /saludo-json.
type Saludo struct {
	Mensaje string `json:"mensaje"`
}

func main() {
	// Usamos un ServeMux explícito en vez del mux por defecto
	// (http.DefaultServeMux) para que las rutas registradas queden claras
	// y aisladas en esta variable, en vez de vivir en un estado global
	// implícito de la biblioteca estándar. En un backend real esto importa
	// porque evita colisiones si otro paquete también registra rutas en
	// el mux por defecto.
	mux := http.NewServeMux()

	mux.HandleFunc("/", manejarInicio)
	mux.HandleFunc("/saludo", manejarSaludo)
	mux.HandleFunc("/eco", manejarEco)
	mux.HandleFunc("/saludo-json", manejarSaludoJSON)

	fmt.Println("servidor escuchando en http://localhost:8081")
	// ListenAndServe bloquea aquí. Si falla (ej. el puerto ya está en
	// uso), devuelve un error que SIEMPRE hay que revisar.
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatal("error arrancando el servidor:", err)
	}
}

// w http.ResponseWriter es por donde ESCRIBES la respuesta.
// r *http.Request contiene todo sobre la petición entrante.
func manejarInicio(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hola, backend")
}

// r.URL.Query() da acceso a los query params (?nombre=valor).
func manejarSaludo(w http.ResponseWriter, r *http.Request) {
	nombre := r.URL.Query().Get("nombre") // devuelve "" si no viene el param
	if nombre == "" {
		nombre = "desconocido"
	}
	fmt.Fprintf(w, "Hola, %s\n", nombre)
}

// Solo aceptamos POST; para cualquier otro método respondemos 405.
func manejarEco(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "solo se acepta POST en /eco")
		return
	}

	// r.Body es un io.Reader: lo leemos completo con io.ReadAll.
	cuerpo, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "error leyendo el body")
		return
	}

	w.Write(cuerpo) // devolvemos exactamente lo que recibimos
}

// Responder JSON: hay que poner el header Content-Type antes de escribir
// el body, y json.NewEncoder(w).Encode(...) serializa directo al writer
// sin que tengas que hacer Marshal + Write por separado.
func manejarSaludoJSON(w http.ResponseWriter, r *http.Request) {
	s := Saludo{Mensaje: "Hola desde JSON"}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(s); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
