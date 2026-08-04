// Día 23 — Middleware HTTP y variables de entorno
//
// Cómo correrlo (sin dependencias externas, solo librería estándar):
//
//	go run ./cmd/extra/middleware-env
//
// Luego, en otra terminal:
//
//	curl http://localhost:8080/productos
//	curl http://localhost:8080/falla     <- este endpoint hace panic a propósito
//
// ---------------------------------------------------------------------------
// ¿QUÉ ES UN MIDDLEWARE?
//
// Un middleware en Go es, literalmente, una función que RECIBE un
// http.Handler y DEVUELVE otro http.Handler. Es el patrón "decorator":
// envolvés tu handler original con capas de comportamiento adicional sin
// tocar su código.
//
//	func(next http.Handler) http.Handler
//
// ¿Por qué importa esto para un ERP? Porque cosas como "loguear cada
// petición", "no tumbar el servidor si un handler entra en pánico",
// "verificar que venga un token válido" o "permitir peticiones desde tu
// frontend (CORS)" son necesidades TRANSVERSALES: las necesitás en decenas
// de endpoints (productos, clientes, facturas...). Sin middleware, tendrías
// que copiar/pegar ese código en cada handler. Con middleware, lo escribís
// una vez y lo aplicás a todos.
//
// ---------------------------------------------------------------------------
// ¿QUÉ SON LAS VARIABLES DE ENTORNO Y POR QUÉ IMPORTAN?
//
// Un ERP real corre en distintos entornos: tu máquina, un servidor de
// pruebas, producción. Cada entorno necesita configuración distinta (puerto,
// credenciales de base de datos, claves de APIs externas...). Si esa
// configuración estuviera escrita directamente en el código (hardcodeada),
// tendrías que recompilar el programa cada vez que cambia algo, y peor:
// terminarías subiendo contraseñas y claves secretas a tu repositorio git.
//
// La solución estándar es leer esa configuración desde VARIABLES DE ENTORNO
// del sistema operativo, con os.Getenv. En desarrollo, en vez de escribir
// `export PUERTO=8080` a mano cada vez, es común guardar esas variables en
// un archivo .env (que NO se sube a git, se agrega a .gitignore) y cargarlo
// al arrancar el programa. La librería más usada para esto es
// github.com/joho/godotenv (requiere `go get github.com/joho/godotenv` con
// conexión a internet). Como esta librería es opcional y no queremos que
// este ejemplo dependa de red para poder correr, más abajo implementamos
// una versión MUY simplificada de "cargar un .env a mano" usando solo
// bufio y os — así entendés el mecanismo antes de usar la librería.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// ---------------------------------------------------------------------------
// Middleware 1: logueo de cada request.
//
// Envolvemos "next" (el handler real) en una función que:
//  1. anota la hora de inicio,
//  2. deja que next atienda la petición normalmente,
//  3. después de que next termina, loguea método, ruta y cuánto tardó.
//
// Esto es exactamente lo que hacen frameworks como Gin o Echo por debajo:
// no hay magia, es una función que envuelve a otra.
func middlewareLogueo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inicio := time.Now()
		next.ServeHTTP(w, r)
		duracion := time.Since(inicio)
		log.Printf("%s %s -> %v", r.Method, r.URL.Path, duracion)
	})
}

// ---------------------------------------------------------------------------
// Middleware 2: recuperación de pánicos.
//
// En Go, un panic() no recuperado en una goroutine tumba TODO el proceso.
// Cada request HTTP en net/http corre en su propia goroutine, así que un
// solo bug en un handler (ej. una división por cero, o indexar un slice
// vacío) podría tumbar tu API entera para TODOS los usuarios conectados.
//
// defer + recover() dentro de este middleware evita eso: si next.ServeHTTP
// entra en pánico, lo atrapamos, respondemos 500 a ESE cliente en particular,
// y el servidor sigue vivo para atender al resto.
func middlewareRecuperarPanico(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC recuperado en %s %s: %v", r.Method, r.URL.Path, err)
				http.Error(w, "error interno del servidor", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// ---------------------------------------------------------------------------
// Middleware 3: CORS básico.
//
// CORS (Cross-Origin Resource Sharing) es lo que le permite a un frontend
// que corre en, por ejemplo, http://localhost:3000 (React) hacer peticiones
// a tu API que corre en http://localhost:8080. Sin estas cabeceras, el
// navegador BLOQUEA la respuesta por seguridad (el backend y el frontend
// están en "orígenes" distintos). Esta es la versión más permisiva posible
// (Access-Control-Allow-Origin: *) — sirve para desarrollo; en producción
// normalmente restringís el origen a tu dominio real.
func middlewareCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Los navegadores, antes de un POST/PUT/DELETE "cross-origin", mandan
		// una petición OPTIONS de "preflight" para preguntar si está permitido.
		// Respondemos 200 sin cuerpo y cortamos ahí, sin llegar al handler real.
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// encadenarMiddlewares aplica una lista de middlewares a un handler, en
// orden. Fijate que el ORDEN importa: si middlewareRecuperarPanico no
// envuelve a TODOS los demás (incluyendo el logueo), un panic dentro del
// logueo o el CORS no se recuperaría. Por eso lo ponemos primero en la
// lista (más "externo").
func encadenarMiddlewares(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

// ---------------------------------------------------------------------------
// cargarEnvSimple imita, de forma muy básica, lo que hace una librería como
// godotenv: lee un archivo (si existe) con líneas "CLAVE=valor" y las
// carga en el entorno del proceso vía os.Setenv, SOLO si esa variable no
// estaba ya definida (para que un `export PUERTO=9090` manual en tu shell
// siempre tenga prioridad sobre el archivo .env).
//
// En un proyecto real usarías:
//
//	go get github.com/joho/godotenv
//	godotenv.Load() // una sola línea, al principio de main()
//
// Acá lo escribimos a mano para que el ejemplo compile y corra sin
// necesidad de conexión a internet, y para que entiendas qué hace esa
// línea mágica por debajo.
func cargarEnvSimple(ruta string) {
	archivo, err := os.Open(ruta)
	if err != nil {
		// No es un error grave: si no hay .env, seguimos con lo que ya
		// esté exportado en el entorno del sistema (o los valores por
		// defecto del programa).
		return
	}
	defer archivo.Close()

	scanner := bufio.NewScanner(archivo)
	for scanner.Scan() {
		linea := strings.TrimSpace(scanner.Text())
		if linea == "" || strings.HasPrefix(linea, "#") {
			continue // línea vacía o comentario
		}
		partes := strings.SplitN(linea, "=", 2)
		if len(partes) != 2 {
			continue // línea mal formada, la ignoramos
		}
		clave := strings.TrimSpace(partes[0])
		valor := strings.TrimSpace(partes[1])
		if _, yaDefinida := os.LookupEnv(clave); !yaDefinida {
			os.Setenv(clave, valor)
		}
	}
}

func main() {
	// Este ejemplo no crea un archivo .env de verdad porque no queremos
	// escribir archivos fuera de este ejercicio, pero si creás uno junto
	// a este archivo con contenido como:
	//
	//   PUERTO=9090
	//   ENTORNO=desarrollo
	//
	// y lo cargás con cargarEnvSimple("archivo.env"), vas a ver que el
	// servidor arranca en el puerto 9090 en vez del 8080 por defecto.
	cargarEnvSimple(".env")

	puerto := os.Getenv("PUERTO")
	if puerto == "" {
		puerto = "8080" // valor por defecto si no hay variable de entorno
	}
	entorno := os.Getenv("ENTORNO")
	if entorno == "" {
		entorno = "desarrollo"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/productos", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `[{"id":1,"nombre":"teclado"},{"id":2,"nombre":"mouse"}]`)
	})

	// Este endpoint existe solo para demostrar que el middleware de
	// recover() funciona: provoca un panic a propósito indexando un slice
	// vacío. Sin middlewareRecuperarPanico, esto tumbaría el servidor
	// entero. Con él, solo esta petición falla con 500 y el servidor sigue
	// respondiendo a todo lo demás.
	mux.HandleFunc("/falla", func(w http.ResponseWriter, r *http.Request) {
		var vacio []int
		fmt.Fprintln(w, vacio[0]) // panic: index out of range
	})

	handlerFinal := encadenarMiddlewares(mux,
		middlewareRecuperarPanico, // primero: atrapa pánicos de TODO lo de adentro
		middlewareLogueo,
		middlewareCORS,
	)

	direccion := ":" + puerto
	log.Printf("entorno=%s | servidor escuchando en http://localhost%s", entorno, direccion)
	log.Printf("probá: curl http://localhost%s/productos", direccion)
	log.Printf("probá: curl http://localhost%s/falla   (el servidor NO debe caerse)", direccion)

	if err := http.ListenAndServe(direccion, handlerFinal); err != nil {
		log.Fatal(err)
	}
}
