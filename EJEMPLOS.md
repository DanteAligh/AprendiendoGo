# Ejemplos de referencia (`cmd\extra\`)

Estos son **programas completos, comentados y que compilan**, importados del plan
de estudio `delicop/aprender_go` y adaptados a este repositorio.

## Qué son y qué NO son

Son **referencia de sintaxis**: cada uno enseña un concepto con un caso de uso
*distinto* al del ejercicio de ese día, precisamente para que no puedas copiarlo
y pegarlo como solución. Cuando te atores en un ejercicio y no recuerdes "cómo se
escribía esto en Go", abre el ejemplo del concepto, míralo, ciérralo y vuelve a
tu ejercicio.

**No son la solución de nada.** La solución la escribes tú.

## Cómo ejecutar cualquiera

Desde la raíz del repositorio, en PowerShell:

```powershell
go run ./cmd/extra/<carpeta>
```

Por ejemplo:

```powershell
go run ./cmd/extra/punteros
```

Cada ejemplo vive en su **propia carpeta** porque todos son `package main` con su
propia `func main()`: si estuvieran juntos en un mismo directorio, Go diría
`main redeclared in this block` y no compilaría ninguno.

## Índice: concepto → día nuestro → carpeta

### Semana 1 — sintaxis

| Concepto | Día | Carpeta | Notas |
|---|---|---|---|
| Variables, tipos, constantes, `strconv` | 1 | `cmd\extra\variables` | |
| Leer de la consola, operadores | 1 | `cmd\extra\entrada` | espera que escribas algo |
| Funciones (retornos múltiples y con nombre, variádicas) | 1 | `cmd\extra\funciones` | |
| `if` / `else` / `switch` | 2 | `cmd\extra\condicionales` | |
| Las tres formas del `for` | 2 | `cmd\extra\ciclos` | |
| Arrays y slices, `append` | 3 | `cmd\extra\slices` | |
| Maps y el idiom "comma-ok" | 3 | `cmd\extra\maps` | |
| Structs | 4 | `cmd\extra\structs` | |
| Métodos (receptor valor vs puntero) | 4 | `cmd\extra\metodos` | |
| Punteros | 4 | `cmd\extra\punteros` | |
| Interfaces | 4 | `cmd\extra\interfaces` | |
| Errores: `errors.New`, `%w`, `errors.Is` | 5 | `cmd\extra\errores` | |
| Crear tu propio paquete | 6 | `cmd\extra\paquetes` + `internal\figuras` | |
| Tests de tabla | 6 | `internal\operaciones` | `go test ./internal/operaciones` |
| Leer y escribir archivos | 7 | `cmd\extra\archivos` | |

### Semana 2 — HTTP

| Concepto | Día | Carpeta |
|---|---|---|
| Servidor HTTP básico, `ServeMux`, handlers | 8 | `cmd\extra\http-servidor` |
| JSON: `Marshal`, `Unmarshal`, struct tags | 9 | `cmd\extra\json` |
| API REST CRUD con almacén en memoria | 11 | `cmd\extra\api-memoria` |
| Middleware encadenado + variables de entorno | 12 | `cmd\extra\middleware-env` |

### Semana 3 — concurrencia y base de datos

| Concepto | Día | Carpeta |
|---|---|---|
| Goroutines y canales | 15 | `cmd\extra\goroutines` |
| `sync.Mutex`, `WaitGroup`, `select` | 16-17 | `cmd\extra\sync-select` |
| `database/sql` con SQLite | 18 | `cmd\extra\sqlite` |

### Semana 4 — producción

| Concepto | Día | Carpeta |
|---|---|---|
| Consumir una API externa (`http.Client`, timeouts) | 25 | `cmd\extra\cliente-http` |
| Autenticación: bcrypt + JWT | extra | `cmd\extra\auth-jwt` |
| Modelado de datos de un ERP | 26 | `cmd\extra\modelado-erp` |
| Andamiaje de API en capas | 27-28 | `semana4\proyecto-final\` |

## Sobre la base de datos: SQLite primero, Postgres después

El ejemplo `cmd\extra\sqlite` usa **SQLite** con el driver `modernc.org/sqlite`,
que es 100% Go: no necesita CGO, ni un compilador de C, ni Docker. En Windows eso
significa que funciona escribiendo `go run` y ya.

Aprendes SQL ahí (día 18-21), y en la semana 4 el mismo código se mueve a
**PostgreSQL dentro de Docker**, que es lo que se usa en producción de verdad.
Lo único que cambia entre uno y otro es el driver y que los placeholders pasan de
`?` a `$1`. El concepto —abrir conexión, `Exec`, `Query`, consultas
parametrizadas— es idéntico.

## `semana4\proyecto-final\`

Es una API REST tipo mini-ERP con arquitectura en capas
(`handler` → `service` → `repository` → `models`). **Compila y arranca**, pero
cada endpoint de negocio responde `501 Not Implemented` a propósito: ese es tu
punto de partida para los días 27-28, no un ejemplo terminado.

```powershell
go run ./semana4/proyecto-final/cmd/api
```
