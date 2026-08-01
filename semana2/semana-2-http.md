# Semana 2 · Servidor HTTP con la librería estándar

**Meta de la semana:** una API REST de tareas, con middleware y tests, escrita solo con `net/http`. Nada de frameworks.

**Carpeta de trabajo:** `cmd\api\` e `internal\tareas\`

---

## Día 1 · Tu primer servidor

**Conceptos:** `http.ListenAndServe`, `http.ServeMux`, `HandlerFunc`, `ResponseWriter` y `*Request`.

Crea `cmd\dia08\main.go`:

```go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// Desde Go 1.22 el patrón puede incluir método y variables de ruta.
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hola desde Go")
	})

	mux.HandleFunc("GET /saludo/{nombre}", func(w http.ResponseWriter, r *http.Request) {
		nombre := r.PathValue("nombre")
		fmt.Fprintf(w, "hola, %s\n", nombre)
	})

	// Parámetros de query: /suma?a=3&b=4
	mux.HandleFunc("GET /eco", func(w http.ResponseWriter, r *http.Request) {
		msg := r.URL.Query().Get("msg")
		if msg == "" {
			// El orden importa: cabeceras, luego código, luego cuerpo.
			http.Error(w, "falta el parámetro msg", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, msg)
	})

	log.Println("escuchando en http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
```

Ejecuta y prueba desde otra terminal:

```powershell
go run ./cmd/dia08
```

```powershell
curl http://localhost:8080/
curl http://localhost:8080/saludo/carlos
curl "http://localhost:8080/eco?msg=probando"
curl -i "http://localhost:8080/eco"
```

> En PowerShell `curl` es un alias de `Invoke-WebRequest` y se comporta distinto. Usa `curl.exe` explícitamente (viene con Windows 10+) para que los ejemplos funcionen igual que en Linux.

**Regla de oro:** una vez llamas a `w.WriteHeader` o escribes en el cuerpo, ya no puedes cambiar cabeceras ni el código de estado. Por eso `http.Error` siempre va seguido de `return`.

### Ejercicio

Añade `GET /calc/{op}/{a}/{b}` que soporte `suma`, `resta`, `mult` y `div`. Devuelve 400 si los números no parsean, y 400 con mensaje claro en la división entre cero. Usa `strconv.ParseFloat` y `r.PathValue`.

---

## Día 2 · JSON de entrada y salida

**Conceptos:** `encoding/json`, etiquetas de struct, `json.Decoder`, validación.

Crea `cmd\dia09\main.go`:

```go
package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type Tarea struct {
	ID     int    `json:"id"`
	Titulo string `json:"titulo"`
	Hecha  bool   `json:"hecha"`
	// omitempty omite el campo si está vacío; "-" lo excluye siempre.
	Nota string `json:"nota,omitempty"`
}

// Helper que usarás en todos los handlers.
func escribirJSON(w http.ResponseWriter, codigo int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(codigo)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Println("error codificando json:", err)
	}
}

func (t Tarea) validar() error {
	if t.Titulo == "" {
		return errors.New("el título es obligatorio")
	}
	if len(t.Titulo) > 200 {
		return errors.New("el título no puede superar 200 caracteres")
	}
	return nil
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /tareas", func(w http.ResponseWriter, r *http.Request) {
		lista := []Tarea{
			{ID: 1, Titulo: "Aprender Go", Hecha: true},
			{ID: 2, Titulo: "Escribir una API"},
		}
		escribirJSON(w, http.StatusOK, lista)
	})

	mux.HandleFunc("POST /tareas", func(w http.ResponseWriter, r *http.Request) {
		var t Tarea
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields() // rechaza campos no esperados
		if err := dec.Decode(&t); err != nil {
			escribirJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if err := t.validar(); err != nil {
			escribirJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
			return
		}
		t.ID = 99
		escribirJSON(w, http.StatusCreated, t)
	})

	log.Println("escuchando en :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

Pruébalo:

```powershell
curl.exe -X POST http://localhost:8080/tareas -H "Content-Type: application/json" -d "{\"titulo\":\"Comprar pan\"}"
curl.exe -X POST http://localhost:8080/tareas -H "Content-Type: application/json" -d "{\"titulo\":\"\"}"
```

**Detalle importante:** solo los campos **exportados** (con mayúscula inicial) se serializan. Si `titulo` fuera minúscula en el struct, `encoding/json` no lo vería.

### Ejercicio

Crea un endpoint `POST /validar` que reciba `{"email":"...", "edad":N}` y devuelva un objeto con la lista de errores de validación. Añade un tipo `RespuestaError struct { Error string; Detalles []string }` y úsalo en todas las respuestas de fallo.

---

## Día 3 · Almacén en memoria y concurrencia segura

**Conceptos:** `sync.RWMutex`, encapsular el estado, inyectar dependencias en los handlers.

Crea `internal\tareas\almacen.go`:

```go
package tareas

import (
	"errors"
	"sync"
)

var ErrNoEncontrada = errors.New("tarea no encontrada")

type Tarea struct {
	ID     int    `json:"id"`
	Titulo string `json:"titulo"`
	Hecha  bool   `json:"hecha"`
}

// Almacen protege el mapa con un mutex: el servidor HTTP atiende
// cada petición en una goroutine distinta, así que varios handlers
// pueden tocar el mapa a la vez.
type Almacen struct {
	mu       sync.RWMutex
	datos    map[int]Tarea
	siguiente int
}

func NuevoAlmacen() *Almacen {
	return &Almacen{datos: make(map[int]Tarea), siguiente: 1}
}

func (a *Almacen) Crear(t Tarea) Tarea {
	a.mu.Lock()
	defer a.mu.Unlock()
	t.ID = a.siguiente
	a.siguiente++
	a.datos[t.ID] = t
	return t
}

func (a *Almacen) Obtener(id int) (Tarea, error) {
	a.mu.RLock() // RLock permite lecturas simultáneas
	defer a.mu.RUnlock()
	t, ok := a.datos[id]
	if !ok {
		return Tarea{}, ErrNoEncontrada
	}
	return t, nil
}

func (a *Almacen) Listar() []Tarea {
	a.mu.RLock()
	defer a.mu.RUnlock()
	out := make([]Tarea, 0, len(a.datos))
	for _, t := range a.datos {
		out = append(out, t)
	}
	return out
}

func (a *Almacen) Actualizar(id int, t Tarea) (Tarea, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if _, ok := a.datos[id]; !ok {
		return Tarea{}, ErrNoEncontrada
	}
	t.ID = id
	a.datos[id] = t
	return t, nil
}

func (a *Almacen) Borrar(id int) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if _, ok := a.datos[id]; !ok {
		return ErrNoEncontrada
	}
	delete(a.datos, id)
	return nil
}
```

Y `internal\tareas\almacen_test.go`:

```go
package tareas

import (
	"errors"
	"sync"
	"testing"
)

func TestCicloCompleto(t *testing.T) {
	a := NuevoAlmacen()

	creada := a.Crear(Tarea{Titulo: "una"})
	if creada.ID != 1 {
		t.Fatalf("ID = %d, quiero 1", creada.ID)
	}

	if _, err := a.Obtener(1); err != nil {
		t.Fatalf("Obtener: %v", err)
	}

	if err := a.Borrar(1); err != nil {
		t.Fatalf("Borrar: %v", err)
	}

	if _, err := a.Obtener(1); !errors.Is(err, ErrNoEncontrada) {
		t.Errorf("err = %v, quiero ErrNoEncontrada", err)
	}
}

// Ejecuta este test con -race para comprobar que el mutex funciona.
func TestConcurrencia(t *testing.T) {
	a := NuevoAlmacen()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			a.Crear(Tarea{Titulo: "x"})
		}()
	}
	wg.Wait()
	if n := len(a.Listar()); n != 100 {
		t.Errorf("hay %d tareas, quiero 100", n)
	}
}
```

Ejecuta:

```powershell
go test -race ./internal/tareas
```

### Ejercicio

Quita temporalmente los `Lock`/`Unlock` de `Crear` y vuelve a correr `go test -race`. Lee el informe de carrera que imprime — entender esa salida te ahorrará horas más adelante. Luego devuélvelos a su sitio.

---

## Día 4 · Estructurar la API con un tipo servidor

**Conceptos:** agrupar handlers como métodos, evitar variables globales.

Crea `internal\tareas\api.go`:

```go
package tareas

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type API struct {
	almacen *Almacen
}

func NuevaAPI(a *Almacen) *API { return &API{almacen: a} }

// Rutas devuelve el mux ya configurado. Así el servidor y los tests
// usan exactamente el mismo enrutado.
func (api *API) Rutas() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /tareas", api.listar)
	mux.HandleFunc("POST /tareas", api.crear)
	mux.HandleFunc("GET /tareas/{id}", api.obtener)
	mux.HandleFunc("PUT /tareas/{id}", api.actualizar)
	mux.HandleFunc("DELETE /tareas/{id}", api.borrar)
	mux.HandleFunc("GET /salud", func(w http.ResponseWriter, r *http.Request) {
		escribirJSON(w, http.StatusOK, map[string]string{"estado": "ok"})
	})
	return mux
}

func escribirJSON(w http.ResponseWriter, codigo int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(codigo)
	_ = json.NewEncoder(w).Encode(v)
}

func errorJSON(w http.ResponseWriter, codigo int, msg string) {
	escribirJSON(w, codigo, map[string]string{"error": msg})
}

func idDeRuta(r *http.Request) (int, error) {
	return strconv.Atoi(r.PathValue("id"))
}

func (api *API) listar(w http.ResponseWriter, r *http.Request) {
	escribirJSON(w, http.StatusOK, api.almacen.Listar())
}

func (api *API) crear(w http.ResponseWriter, r *http.Request) {
	var t Tarea
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		errorJSON(w, http.StatusBadRequest, "json inválido")
		return
	}
	if t.Titulo == "" {
		errorJSON(w, http.StatusUnprocessableEntity, "el título es obligatorio")
		return
	}
	escribirJSON(w, http.StatusCreated, api.almacen.Crear(t))
}

func (api *API) obtener(w http.ResponseWriter, r *http.Request) {
	id, err := idDeRuta(r)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "id inválido")
		return
	}
	t, err := api.almacen.Obtener(id)
	if errors.Is(err, ErrNoEncontrada) {
		errorJSON(w, http.StatusNotFound, "no encontrada")
		return
	}
	escribirJSON(w, http.StatusOK, t)
}

func (api *API) actualizar(w http.ResponseWriter, r *http.Request) {
	id, err := idDeRuta(r)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "id inválido")
		return
	}
	var t Tarea
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		errorJSON(w, http.StatusBadRequest, "json inválido")
		return
	}
	actualizada, err := api.almacen.Actualizar(id, t)
	if errors.Is(err, ErrNoEncontrada) {
		errorJSON(w, http.StatusNotFound, "no encontrada")
		return
	}
	escribirJSON(w, http.StatusOK, actualizada)
}

func (api *API) borrar(w http.ResponseWriter, r *http.Request) {
	id, err := idDeRuta(r)
	if err != nil {
		errorJSON(w, http.StatusBadRequest, "id inválido")
		return
	}
	if err := api.almacen.Borrar(id); errors.Is(err, ErrNoEncontrada) {
		errorJSON(w, http.StatusNotFound, "no encontrada")
		return
	}
	w.WriteHeader(http.StatusNoContent) // 204: sin cuerpo
}
```

Y `cmd\api\main.go`:

```go
package main

import (
	"log"
	"net/http"

	"aprendiendo-go/internal/tareas"
)

func main() {
	api := tareas.NuevaAPI(tareas.NuevoAlmacen())
	log.Println("API en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", api.Rutas()))
}
```

Ejecuta y prueba el CRUD completo:

```powershell
go run ./cmd/api
```

```powershell
curl.exe -X POST http://localhost:8080/tareas -H "Content-Type: application/json" -d "{\"titulo\":\"Estudiar\"}"
curl.exe http://localhost:8080/tareas
curl.exe http://localhost:8080/tareas/1
curl.exe -X PUT http://localhost:8080/tareas/1 -H "Content-Type: application/json" -d "{\"titulo\":\"Estudiar Go\",\"hecha\":true}"
curl.exe -i -X DELETE http://localhost:8080/tareas/1
curl.exe -i http://localhost:8080/tareas/999
```

### Ejercicio

Añade `GET /tareas?hecha=true` para filtrar, y paginación con `?limite=10&desde=0`. Devuelve la respuesta como `{"datos": [...], "total": N}`.

---

## Día 5 · Middleware

**Conceptos:** un middleware es una función que recibe un `http.Handler` y devuelve otro.

Crea `internal\tareas\middleware.go`:

```go
package tareas

import (
	"log"
	"net/http"
	"time"
)

// responseWriter envuelve el original para capturar el código de estado,
// que de otro modo no es accesible después de escribirlo.
type responseWriter struct {
	http.ResponseWriter
	codigo int
}

func (rw *responseWriter) WriteHeader(c int) {
	rw.codigo = c
	rw.ResponseWriter.WriteHeader(c)
}

func Logging(siguiente http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inicio := time.Now()
		rw := &responseWriter{ResponseWriter: w, codigo: http.StatusOK}
		siguiente.ServeHTTP(rw, r)
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, rw.codigo, time.Since(inicio))
	})
}

// Recuperar evita que un panic en un handler tumbe todo el servidor.
func Recuperar(siguiente http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC: %v", err)
				errorJSON(w, http.StatusInternalServerError, "error interno")
			}
		}()
		siguiente.ServeHTTP(w, r)
	})
}

// Encadenar aplica los middleware de forma que el primero de la lista
// sea el más externo.
func Encadenar(h http.Handler, ms ...func(http.Handler) http.Handler) http.Handler {
	for i := len(ms) - 1; i >= 0; i-- {
		h = ms[i](h)
	}
	return h
}
```

Actualiza `cmd\api\main.go`:

```go
handler := tareas.Encadenar(api.Rutas(), tareas.Recuperar, tareas.Logging)
log.Fatal(http.ListenAndServe(":8080", handler))
```

Para probar `Recuperar`, añade temporalmente una ruta que haga `panic("boom")` y comprueba que el servidor sigue vivo después de llamarla.

### Ejercicio

Escribe dos middleware más: uno que exija la cabecera `X-API-Key` con un valor fijo (401 si falta) y otro que limite el tamaño del cuerpo con `http.MaxBytesReader`. Aplica el de autenticación solo a las rutas de escritura.

---

## Día 6 · Tests con `httptest`

**Conceptos:** probar handlers sin abrir puertos.

Crea `internal\tareas\api_test.go`:

```go
package tareas

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func nuevaAPIDePrueba() http.Handler {
	return NuevaAPI(NuevoAlmacen()).Rutas()
}

func TestCrearTarea(t *testing.T) {
	h := nuevaAPIDePrueba()

	req := httptest.NewRequest("POST", "/tareas", strings.NewReader(`{"titulo":"probar"}`))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("código = %d, quiero 201. Cuerpo: %s", rec.Code, rec.Body)
	}

	var creada Tarea
	if err := json.NewDecoder(rec.Body).Decode(&creada); err != nil {
		t.Fatalf("decodificando: %v", err)
	}
	if creada.Titulo != "probar" || creada.ID == 0 {
		t.Errorf("tarea = %+v", creada)
	}
}

func TestErrores(t *testing.T) {
	h := nuevaAPIDePrueba()

	casos := []struct {
		nombre string
		metodo string
		ruta   string
		cuerpo string
		quiero int
	}{
		{"json roto", "POST", "/tareas", `{`, http.StatusBadRequest},
		{"sin titulo", "POST", "/tareas", `{"titulo":""}`, http.StatusUnprocessableEntity},
		{"id no numérico", "GET", "/tareas/abc", "", http.StatusBadRequest},
		{"inexistente", "GET", "/tareas/999", "", http.StatusNotFound},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			req := httptest.NewRequest(c.metodo, c.ruta, strings.NewReader(c.cuerpo))
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			if rec.Code != c.quiero {
				t.Errorf("código = %d, quiero %d", rec.Code, c.quiero)
			}
		})
	}
}
```

Ejecuta:

```powershell
go test -v -race ./internal/tareas
go test -cover ./...
```

### Ejercicio

Añade tests para actualizar y borrar, incluyendo el 204 sin cuerpo. Escribe un test que verifique que el middleware de logging no altera el código de estado.

---

## Día 7 · Servidor de producción y apagado ordenado

Reescribe `cmd\api\main.go` con la configuración que sí usarías en producción:

```go
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"aprendiendo-go/internal/tareas"
)

func main() {
	puerto := os.Getenv("PUERTO")
	if puerto == "" {
		puerto = "8080"
	}

	api := tareas.NuevaAPI(tareas.NuevoAlmacen())
	handler := tareas.Encadenar(api.Rutas(), tareas.Recuperar, tareas.Logging)

	// ListenAndServe sin timeouts deja el servidor expuesto a
	// conexiones lentas que agotan recursos. Siempre configúralos.
	srv := &http.Server{
		Addr:         ":" + puerto,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// El servidor arranca en su propia goroutine para que main
	// pueda quedarse esperando la señal de apagado.
	go func() {
		log.Println("escuchando en :" + puerto)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error del servidor: %v", err)
		}
	}()

	parar := make(chan os.Signal, 1)
	signal.Notify(parar, os.Interrupt, syscall.SIGTERM)
	<-parar // bloquea aquí hasta recibir Ctrl+C

	log.Println("apagando...")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("apagado forzado: %v", err)
	}
	log.Println("adiós")
}
```

Ejecuta, lanza una petición y pulsa `Ctrl+C` para ver el apagado ordenado:

```powershell
$env:PUERTO="3000"; go run ./cmd/api
```

### Ejercicio final de la semana

1. Mueve la configuración (puerto, timeouts, clave de API) a un struct `Config` que se lea de variables de entorno con valores por defecto.
2. Sustituye `log` por `log/slog` con salida en JSON: `slog.New(slog.NewJSONHandler(os.Stdout, nil))`.
3. Añade a cada petición un ID único (genera uno con `crypto/rand`), guárdalo en el `context` y sácalo en todos los logs de esa petición.
4. Comprueba que `go vet ./...` y `go test -race ./...` pasan limpios.

---

## Checklist antes de pasar a la semana 3

- [ ] Sé por qué `http.Error` va siempre seguido de `return`
- [ ] Entiendo que cada petición corre en su propia goroutine
- [ ] Escribo un middleware sin mirar el ejemplo
- [ ] Pruebo handlers con `httptest` sin levantar el servidor
- [ ] Configuro timeouts y apagado ordenado
