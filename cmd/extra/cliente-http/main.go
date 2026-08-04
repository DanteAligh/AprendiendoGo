// Día 25 — Consumir una API externa desde Go (cliente HTTP)
//
// Cómo correrlo (solo librería estándar, sin dependencias externas):
//
//	go run ./cmd/extra/cliente-http
//
// ---------------------------------------------------------------------------
// CONTEXTO
//
// Hasta ahora escribiste servidores HTTP (tu API responde peticiones). Hoy
// escribimos el lado opuesto: tu programa Go actúa como CLIENTE y le hace
// peticiones a OTRO servidor. Esto es exactamente lo que vas a hacer cuando
// tu ERP necesite hablar con un servicio externo: una pasarela de pagos, un
// servicio de envío de emails, o —lo que te interesa para tu meta final—
// una API de un modelo de IA (como la API de Anthropic o de OpenAI) para,
// por ejemplo, generar un resumen automático de una factura o clasificar
// un ticket de soporte.
//
// El flujo siempre es el mismo, sin importar qué API externa sea:
//  1. Armar el cuerpo de la petición (normalmente JSON).
//  2. Crear un *http.Request con el método, la URL, el cuerpo y las
//     cabeceras (headers) necesarias — casi siempre "Content-Type" y
//     "Authorization".
//  3. Enviarlo con un *http.Client.
//  4. Leer el cuerpo de la respuesta y decodificar el JSON a un struct Go.
//  5. Manejar errores: la red puede fallar, el servidor externo puede
//     responder con un código de error (4xx, 5xx), el JSON puede venir
//     distinto a lo esperado.
//
// ---------------------------------------------------------------------------
// ¿POR QUÉ USAMOS httptest.NewServer EN VEZ DE UNA API REAL?
//
// Para que este ejemplo sea 100% reproducible y no dependa de que tengas
// una API key válida de un servicio de IA real (ni de que ese servicio
// esté disponible ahora mismo), levantamos un servidor HTTP DE PRUEBA
// dentro del mismo programa con net/http/httptest. Ese servidor de prueba
// simula cómo respondería una API real de IA: recibe un POST con un
// "prompt" y devuelve JSON con una "respuesta" generada.
//
// La parte importante: el CÓDIGO DEL CLIENTE que escribimos (la función
// consultarAsistenteIA) es EXACTAMENTE el mismo código que usarías contra
// una API real. Lo único que cambia entre este ejemplo y un caso real es
// la URL a la que apunta y que necesitarías una API key de verdad en el
// header Authorization. Todo lo demás —armar el request, mandar headers,
// parsear JSON, manejar errores— es idéntico.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"time"
)

// PeticionIA es lo que le mandamos a la API externa. En una API real de IA
// esto tendría más campos (modelo, temperatura, máximo de tokens...), pero
// la idea es la misma: un struct Go que se serializa a JSON.
type PeticionIA struct {
	Prompt string `json:"prompt"`
}

// RespuestaIA es lo que esperamos que la API externa nos devuelva. Definir
// un struct para la respuesta (en vez de trabajar con map[string]interface{})
// nos da chequeo de tipos y autocompletado — el mismo patrón que ya usaste
// con encoding/json en la semana 3.
type RespuestaIA struct {
	Respuesta string `json:"respuesta"`
	Modelo    string `json:"modelo"`
}

// crearServidorDePruebaIA levanta un servidor HTTP falso que imita cómo se
// comporta una API de IA real: exige un header Authorization, exige que el
// método sea POST, y responde con un JSON generado a partir del prompt
// recibido. httptest.NewServer arranca este servidor en un puerto libre de
// localhost y nos da su URL real — es un servidor HTTP de verdad, no un
// simulacro en memoria, así que el cliente que escribamos habla HTTP real
// con él.
func crearServidorDePruebaIA() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "método no permitido", http.StatusMethodNotAllowed)
			return
		}
		if r.Header.Get("Authorization") != "Bearer clave-de-prueba-123" {
			http.Error(w, "no autorizado", http.StatusUnauthorized)
			return
		}

		var peticion PeticionIA
		if err := json.NewDecoder(r.Body).Decode(&peticion); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		respuesta := RespuestaIA{
			Respuesta: fmt.Sprintf("Respuesta simulada para el prompt: %q", peticion.Prompt),
			Modelo:    "modelo-de-prueba-v1",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(respuesta)
	}))
}

// consultarAsistenteIA es el CLIENTE HTTP en sí. Este es el código que se
// parece a lo que vas a escribir cuando integres tu ERP con una API de IA
// real: solo tendrías que cambiar urlBase y la claveAPI por los valores
// reales del proveedor que uses.
func consultarAsistenteIA(urlBase, claveAPI, prompt string) (*RespuestaIA, error) {
	// Paso 1: armar el cuerpo JSON de la petición.
	cuerpoPeticion := PeticionIA{Prompt: prompt}
	cuerpoJSON, err := json.Marshal(cuerpoPeticion)
	if err != nil {
		return nil, fmt.Errorf("codificando la petición: %w", err)
	}

	// Paso 2: crear el *http.Request explícitamente (en vez de usar
	// http.Post directamente) porque necesitamos agregar el header
	// Authorization. http.NewRequest no envía nada todavía, solo arma el
	// objeto en memoria.
	peticion, err := http.NewRequest(http.MethodPost, urlBase+"/v1/generar", bytes.NewReader(cuerpoJSON))
	if err != nil {
		return nil, fmt.Errorf("creando el request: %w", err)
	}
	peticion.Header.Set("Content-Type", "application/json")
	// Este es el patrón estándar de autenticación con la mayoría de las
	// APIs de IA (Anthropic, OpenAI, etc.): un token en el header
	// Authorization con el prefijo "Bearer ".
	peticion.Header.Set("Authorization", "Bearer "+claveAPI)

	// Un cliente HTTP con timeout es OBLIGATORIO en código de producción:
	// sin timeout, si el servidor externo se cuelga, tu programa podría
	// quedarse esperando para siempre, con una goroutine y una conexión
	// bloqueadas indefinidamente.
	cliente := &http.Client{Timeout: 10 * time.Second}

	// Paso 3: enviar la petición de verdad.
	respuestaHTTP, err := cliente.Do(peticion)
	if err != nil {
		return nil, fmt.Errorf("llamando a la API externa: %w", err)
	}
	defer respuestaHTTP.Body.Close()

	// Paso 4 (parte 1): revisar el código de estado ANTES de intentar
	// decodificar el cuerpo. Muchas APIs devuelven JSON de error también en
	// respuestas 4xx/5xx, pero tiene una forma distinta al JSON de éxito.
	if respuestaHTTP.StatusCode != http.StatusOK {
		cuerpoError, _ := io.ReadAll(respuestaHTTP.Body)
		return nil, fmt.Errorf("la API externa respondió %d: %s",
			respuestaHTTP.StatusCode, string(cuerpoError))
	}

	// Paso 4 (parte 2): decodificar el JSON de éxito a nuestro struct.
	var resultado RespuestaIA
	if err := json.NewDecoder(respuestaHTTP.Body).Decode(&resultado); err != nil {
		return nil, fmt.Errorf("decodificando la respuesta: %w", err)
	}

	return &resultado, nil
}

func main() {
	// Levantamos nuestro servidor de prueba que simula la API de IA.
	servidorPrueba := crearServidorDePruebaIA()
	defer servidorPrueba.Close()

	fmt.Println("Servidor de prueba (simula una API de IA) escuchando en:", servidorPrueba.URL)

	fmt.Println("\n--- Caso 1: petición correcta ---")
	resultado, err := consultarAsistenteIA(servidorPrueba.URL, "clave-de-prueba-123",
		"Resumime esta factura en una línea")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Respuesta recibida: %s\n", resultado.Respuesta)
		fmt.Printf("Modelo usado: %s\n", resultado.Modelo)
	}

	fmt.Println("\n--- Caso 2: clave de API incorrecta (debe fallar con 401) ---")
	_, err = consultarAsistenteIA(servidorPrueba.URL, "clave-incorrecta",
		"Este prompt no debería procesarse")
	if err != nil {
		fmt.Println("Como se esperaba, falló:", err)
	}
}
