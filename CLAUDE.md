# CLAUDE.md

Contexto permanente de este repositorio. Léelo antes de responder cualquier cosa aquí.

## 1. Qué es este proyecto

Es una **carpeta de aprendizaje**, no un producto. El objetivo es que el usuario aprenda Go
en 4 semanas (28 días), partiendo de **cero absoluto**: no sabe programar en ningún lenguaje.

- Idioma de trabajo: **español**. Código, comentarios, mensajes y explicaciones en español
  (los identificadores de Go también: `dividir`, `esFiebre`, `temperaturas`).
- Sistema: **Windows 11 + PowerShell**. Rutas en prosa con `\`, rutas de Go con `/`.
- Go instalado: **1.26.5**. Módulo: `aprendiendo-go` (ver `go.mod`).
- Sin frameworks: todo con la **librería estándar** hasta la semana 3 (ahí entra el driver de Postgres).
- No es un repositorio git todavía. Git aparece en la semana 4, día 5.

## 2. Cómo tratar al usuario (importante)

Está empezando desde cero. Al explicar:

- **Nunca dar por sabido vocabulario.** "Compilar", "puntero", "struct", "flag", "endpoint"
  se explican la primera vez que aparecen, con una analogía del mundo real.
- **Una idea nueva por vez.** Si un ejemplo necesita tres conceptos nuevos, se parte en tres.
- **Todo código va comentado** línea a línea cuando introduce algo nuevo.
- **Siempre incluir el comando exacto** para ejecutar lo que se escribió, y **qué debería salir**
  en pantalla. Si sale otra cosa, ya sabe que algo falló.
- **Explicar los errores del compilador**, no solo arreglarlos. El error es material de estudio.
- No adelantar temas de semanas futuras. Si aparece, se nombra y se dice "eso es de la semana N".
- No hacerle el ejercicio. Dar pistas, revisar lo que escribió, señalar el error y por qué.

## 3. Estructura del repositorio

```
aprendiendo-go\
├── go.mod                       identificación del módulo
├── CLAUDE.md                    este archivo
├── CONCEPTOS-BASICOS.md         glosario desde cero (archivo, terminal, variable, tipo, función)
├── semana1\semana-1-sintaxis.md          días 1-7   → sintaxis + CLI de estadísticas
│   └── diaNN-ejercicios.md               días 01-07: explicación + ejercicios A/B/C
│                                         (el 07 es el proyecto CLI, por fases)
├── semana2\semana-2-http.md              días 8-14  → API REST con net/http
├── semana3\semana-3-concurrencia-db.md   días 15-21 → goroutines, canales, Postgres
├── semana4\semana-4-docker-despliegue.md días 22-28 → Docker, CI, despliegue
└── cmd\
    ├── dia01\main.go            programa del día (copiado del material)
    └── dia01ej\main.go          resolución del ejercicio del día
```

**Convención de carpetas** (respetarla siempre):

| Qué | Dónde |
|---|---|
| Programa de ejemplo del día N | `cmd\diaNN\main.go` |
| Ejercicio del día N resuelto | `cmd\diaNNej\main.go` |
| Código reutilizable (paquetes) | `internal\<paquete>\` |

Cada carpeta bajo `cmd\` es un programa independiente con su propio `package main` y su `func main()`.
Por eso pueden convivir muchos `main` en el mismo módulo: **están en carpetas distintas**.

## 4. Comandos

```powershell
go run ./cmd/dia01          # compila y ejecuta, sin dejar archivo
go build -o bin/x.exe ./cmd/dia01
go test ./...               # todos los tests del módulo
go fmt ./...                # formatea (en Go no se discute el estilo)
go vet ./...                # detecta errores que sí compilan pero están mal
go doc fmt.Printf           # documentación sin salir de la terminal
```

## 5. Estado del progreso

> Mantener esta sección al día: marcar `[x]` cuando el día esté hecho y verificado.

- [x] Día 1 — variables, tipos, funciones con múltiples retornos (`cmd\dia01`, `cmd\dia01ej`)
- [ ] Días 2-7 — semana 1
- [ ] Días 8-14 — semana 2
- [ ] Días 15-21 — semana 3
- [ ] Días 22-28 — semana 4

## 6. Ejemplo canónico explicado

Este es el nivel de detalle esperado en cada explicación. Programa: `cmd\dia01\main.go`.

```go
package main            // (1)

import "fmt"            // (2)

var version = "1.0"     // (3)

func dividir(a, b float64) (float64, bool) {   // (4)
	if b == 0 {
		return 0, false                        // (5)
	}
	return a / b, true
}

func main() {                                  // (6)
	var nombre string = "Go"                   // (7)
	edad := 16                                 // (8)
	pi := 3.14159

	fmt.Printf("Lenguaje: %s, edad: %d, pi: %.2f, version: %s\n",
		nombre, edad, pi, version)             // (9)

	resultado, ok := dividir(10, 3)            // (10)
	if !ok {
		fmt.Println("no se puede dividir entre cero")
		return                                 // (11)
	}
	fmt.Printf("10/3 = %.4f\n", resultado)
}
```

1. **`package main`** — todo archivo `.go` pertenece a un paquete (una caja de código).
   El nombre `main` es especial: significa "esto es un programa que se puede ejecutar",
   no una librería. Sin él, `go run` no sabría qué arrancar.
2. **`import "fmt"`** — pide prestada la caja `fmt` (de *format*), que trae las funciones
   para escribir en pantalla. Go es estricto: si importas algo y no lo usas, **no compila**.
3. **`var version = "1.0"`** — variable a nivel de paquete: existe en todo el archivo.
   Fuera de una función **no se puede usar `:=`**, hay que escribir `var`.
4. **La firma de la función** — `a, b float64` es abreviatura de `a float64, b float64`.
   `(float64, bool)` significa que devuelve **dos valores**: el resultado y un aviso de si salió bien.
   Esta pareja *(valor, ¿ok?)* es la base del manejo de errores en Go; más adelante el `bool`
   se convierte en un `error` de verdad (día 5).
5. **`return 0, false`** — dividir entre cero no tiene respuesta, así que devuelve un cero
   de relleno y un `false` que dice "no me hagas caso". Go **obliga** a devolver los dos valores.
6. **`func main()`** — el punto de entrada. Cuando el programa arranca, empieza aquí.
   Cuando esta función termina, el programa termina.
7. **`var nombre string = "Go"`** — forma larga: nombre, tipo y valor. Correcta pero rara en Go
   dentro de una función; se usa cuando quieres dejar el tipo explícito.
8. **`edad := 16`** — forma corta. El `:=` **crea** la variable e **infiere** el tipo
   mirando el valor: `16` es `int`, `3.14159` es `float64`. Solo funciona dentro de funciones.
9. **`Printf`** — imprime con formato. Cada `%` es un hueco que se rellena en orden con los
   argumentos: `%s` texto, `%d` entero, `%.2f` decimal con 2 cifras, `\n` salto de línea.
   (`Println` sería la versión sin huecos ni control de decimales.)
10. **`resultado, ok := dividir(10, 3)`** — recoge los dos valores en dos variables.
    Si solo te interesara uno, el otro se descarta con `_`, que es el "cubo de basura" de Go.
11. **`return` a secas** — como `main` no devuelve nada, `return` significa "termina aquí".

Ejecutar:

```powershell
go run ./cmd/dia01
```

Salida esperada:

```
Lenguaje: Go, edad: 16, pi: 3.14, version: 1.0
10/3 = 3.3333
```

## 7. Ejercicios diarios (28 días, dificultad creciente)

Un ejercicio por día, **además** del programa de ejemplo que trae cada archivo de semana.
Cada uno se guarda en `cmd\diaNNej\main.go` (o donde se indique) y debe compilar y correr.

### Semana 1 · Sintaxis (días 1-7) — nivel: pasos firmes

| Día | Ejercicio | Concepto que ejercita |
|---|---|---|
| 1 | Convertir Celsius a Fahrenheit sobre una lista de temperaturas y marcar cuáles son fiebre (≥38). **Hecho.** | funciones, `float64`, `bool`, `Printf` |
| 2 | Clasificador de notas 0-100: `for` del 0 al 100 de 10 en 10, y un `switch` que imprima "suspenso / aprobado / notable / sobresaliente". Añadir `continue` para saltarse los múltiplos de 30. | `for`, `if/else`, `switch`, `continue` |
| 3 | Dado un slice de palabras, devolver un `map[string]int` con cuántas veces aparece cada una, e imprimir cuál es la más repetida. | slices, maps, `range`, `append` |
| 4 | `type Producto struct{ Nombre string; Precio float64; Stock int }`. Método `Valor()` que devuelva `Precio*Stock`, y método con **puntero** `Vender(n int)` que reste stock. Interfaz `Describible` con `Describir() string` implementada por `Producto`. | structs, métodos, valor vs puntero, interfaces |
| 5 | Función `leerPrecio(txt string) (float64, error)` que devuelva un error propio (`errors.New` + `fmt.Errorf` con `%w`) si el texto no es un número o es negativo. Comprobarlo con `errors.Is`. Usar `defer` para imprimir "fin de la lectura" pase lo que pase. | `error`, `%w`, `errors.Is`, `defer` |
| 6 | Paquete `internal\stats`: `Media`, `Mediana`, `Max`, `Min` sobre `[]float64`, cada una con su error para slice vacío. Test de tabla que cubra: vacío, un elemento, negativos, pares e impares. | paquetes, tests de tabla, `t.Run` |
| 7 | **Proyecto:** CLI `cmd\stats` que reciba `-archivo ventas.csv -columna 2` con `flag`, lea el CSV con `encoding/csv`, ignore la cabecera, y muestre media, mediana, máximo y mínimo. Debe fallar con mensaje claro (y código de salida 1) si el archivo no existe o la columna no es numérica. | `flag`, `os`, `encoding/csv`, errores reales |

### Semana 2 · HTTP (días 8-14) — nivel: primeras piezas de verdad

| Día | Ejercicio | Concepto que ejercita |
|---|---|---|
| 8 | Servidor con tres rutas: `/`, `/hora` (devuelve la hora actual) y `/saludo?nombre=X`. Si falta `nombre`, responder **400** con un texto explicativo. | `ServeMux`, `HandlerFunc`, query params, códigos HTTP |
| 9 | `POST /validar` que reciba `{"email":"...","edad":N}` y devuelva la lista de errores de validación en un `RespuestaError struct{ Error string; Detalles []string }` usado en **todos** los fallos. | `encoding/json`, tags, decodificar cuerpo |
| 10 | Almacén en memoria de tareas (`internal\tareas`) protegido con `sync.RWMutex`: `Crear`, `Listar`, `Buscar`, `Borrar`. Correr `go test -race` y explicar qué detecta. | mapas compartidos, mutex, `-race` |
| 11 | Convertir los handlers sueltos en un `type Servidor struct` con el almacén dentro y un método `Rutas()`. CRUD completo de `/tareas` y `/tareas/{id}` (Go 1.22+ soporta parámetros en el patrón). | inyección de dependencias, organización |
| 12 | Middleware `Log` (método, ruta, duración, código de estado) y middleware `Auth` que exija la cabecera `X-Token`. Encadenarlos. Requiere envolver el `ResponseWriter` para capturar el código. | funciones que devuelven funciones, composición |
| 13 | Tests con `httptest` de todo el CRUD: crear, listar, 404 al buscar inexistente, 401 sin token, 400 con JSON inválido. | `httptest`, tests de handlers |
| 14 | **Proyecto:** servidor de producción con timeouts (`ReadTimeout`, `WriteTimeout`) y **apagado ordenado** con `signal.NotifyContext` + `srv.Shutdown(ctx)`. Comprobar que una petición en curso termina antes de cerrar. | `http.Server`, señales, `context` |

### Semana 3 · Concurrencia y base de datos (días 15-21) — nivel: exigente

| Día | Ejercicio | Concepto que ejercita |
|---|---|---|
| 15 | Lanzar 10 goroutines que descarguen (simulen) una URL cada una; esperar a todas con `sync.WaitGroup` y medir el tiempo total. Compararlo con hacerlo en serie. | `go`, `WaitGroup`, orden no garantizado |
| 16 | Pipeline de tres etapas conectadas por canales: generar números → elevar al cuadrado → sumar. Cerrar cada canal en su sitio y consumir con `range`. | canales, `close`, dirección de canal |
| 17 | Worker que consulta un servicio lento; cancelarlo con `context.WithTimeout` de 2s usando `select`. Distinguir en la salida "resultado" de "cancelado por timeout". | `select`, `context`, cancelación |
| 18 | Levantar Postgres con Docker, crear la tabla `tareas`, e implementar `internal\db` con `database/sql`: `Crear`, `Listar`, `Buscar`, `Actualizar`, `Borrar`, siempre con consultas parametrizadas (`$1`) y `context`. | SQL, `database/sql`, inyección SQL |
| 19 | `MoverTarea(origen, destino int)` dentro de una transacción: si el segundo `UPDATE` falla, nada se guarda. Demostrarlo forzando el fallo. Patrón `defer tx.Rollback()`. | `Begin`, `Commit`, `Rollback` |
| 20 | Worker en segundo plano que lee trabajos de un canal con buffer, los procesa con N goroutines, reintenta 3 veces si falla, y se apaga limpiamente cuando se cierra el canal. | pool de workers, reintentos, apagado |
| 21 | **Proyecto:** la API de la semana 2 usando Postgres en vez de memoria, con el worker enganchado: al crear una tarea se encola un trabajo que la marca como procesada. Tests de integración contra la base real. | integración de todo |

### Semana 4 · Docker y despliegue (días 22-28) — nivel: producción

| Día | Ejercicio | Concepto que ejercita |
|---|---|---|
| 22 | `internal\config`: leer `PUERTO`, `DATABASE_URL`, `NIVEL_LOG` del entorno con valores por defecto, y **fallar al arrancar** con mensaje claro si falta una obligatoria. | variables de entorno, 12 factores |
| 23 | `Dockerfile` multi-etapa: compilar con `CGO_ENABLED=0`, imagen final `scratch` o `distroless`, usuario no root. Objetivo: **menos de 20 MB**. Verificar con `docker images`. | compilación estática, capas, tamaño |
| 24 | `docker-compose.yml` con la API y Postgres, `healthcheck` en la base y `depends_on: condition: service_healthy`, más volumen para que los datos sobrevivan a `docker compose down`. | Compose, redes, volúmenes |
| 25 | Endpoint `/salud` que compruebe de verdad la base (`db.PingContext`) y devuelva 200/503, más logs estructurados con `log/slog` en JSON incluyendo un id de petición por request. | observabilidad, `slog` |
| 26 | `git init`, `.gitignore` correcto (binarios, `.env`), primer commit, y `.github\workflows\ci.yml` que corra `go vet`, `go test -race ./...` y construya la imagen en cada push. | git, CI |
| 27 | Desplegar la imagen en un servicio gratuito, con las variables de entorno configuradas allí, y comprobar `/salud` desde internet. | despliegue |
| 28 | **Cierre:** `README.md` que explique cómo levantarlo en local y en Docker, más un repaso escrito: qué es una goroutine, por qué `append` se reasigna, qué hace `defer`, cuándo un receptor va por puntero. Si no se puede explicar, se vuelve a ese día. | consolidación |

## 8. Reglas para el asistente al revisar ejercicios

1. Ejecutar `go run` / `go test` de verdad antes de decir que algo funciona.
2. Correr `go vet ./...` y `go fmt ./...` sobre lo escrito.
3. Al corregir: primero **qué falla y por qué**, después el código arreglado.
4. Si el usuario pide "hazlo tú", hacerlo, pero explicando cada línea nueva.
5. Al terminar un día, marcar la casilla de la sección 5.
