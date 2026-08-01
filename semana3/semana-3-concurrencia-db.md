# Semana 3 · Concurrencia, canales y base de datos

**Meta de la semana:** migrar la API a PostgreSQL y añadirle un worker en segundo plano que procese trabajos por canal.

**Carpeta de trabajo:** `internal\db\`, `internal\worker\`

**Requisito previo:** instala Docker Desktop desde `https://www.docker.com/products/docker-desktop/`. Lo usarás hoy para la base de datos y la semana que viene para el despliegue.

---

## Día 1 · Goroutines y WaitGroup

**Conceptos:** `go`, `sync.WaitGroup`, por qué el orden no está garantizado.

Crea `cmd\dia15\main.go`:

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func trabajar(id int, wg *sync.WaitGroup) {
	// defer wg.Done() SIEMPRE lo primero: garantiza que se llama
	// aunque la función salga antes de tiempo.
	defer wg.Done()
	time.Sleep(time.Duration(id*100) * time.Millisecond)
	fmt.Printf("trabajador %d terminó\n", id)
}

func main() {
	inicio := time.Now()

	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1) // Add ANTES de lanzar la goroutine, nunca dentro
		go trabajar(i, &wg)
	}
	wg.Wait() // bloquea hasta que el contador llegue a 0

	fmt.Printf("todo listo en %v\n", time.Since(inicio))

	// Sin WaitGroup, main terminaría antes y estas goroutines
	// morirían sin ejecutarse. Pruébalo comentando wg.Wait().
	resultados := make([]int, 5)
	var wg2 sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			// Escribir cada goroutine en SU índice es seguro
			// sin mutex: no hay solapamiento de memoria.
			resultados[i] = i * i
		}()
	}
	wg2.Wait()
	fmt.Println("cuadrados:", resultados)
}
```

Ejecuta:

```powershell
go run ./cmd/dia15
go run -race ./cmd/dia15
```

> Nota histórica útil: antes de Go 1.22, la variable `i` del bucle se compartía entre iteraciones y este código daba resultados incorrectos. Verás muchos tutoriales antiguos con `go func(i int){...}(i)` para sortearlo. Con Go moderno ya no hace falta, pero reconocerás el patrón.

### Ejercicio

Escribe un programa que descargue 5 URLs en paralelo con `http.Get` y muestre el código de estado y el tiempo de cada una. Compara el tiempo total contra hacerlo en secuencia.

---

## Día 2 · Canales

**Conceptos:** canales sin y con búfer, dirección de canal, `close`, `range` sobre canal.

Crea `cmd\dia16\main.go`:

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

// chan<- solo escritura, <-chan solo lectura. Declarar la dirección
// documenta la intención y el compilador la verifica.
func productor(salida chan<- int, n int) {
	for i := 1; i <= n; i++ {
		salida <- i
		time.Sleep(50 * time.Millisecond)
	}
	close(salida) // cerrar avisa a los lectores de que no habrá más
}

func consumidor(id int, entrada <-chan int, resultados chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	// range sobre un canal termina solo cuando el canal se cierra.
	for v := range entrada {
		resultados <- fmt.Sprintf("consumidor %d procesó %d", id, v)
	}
}

func main() {
	// Sin búfer: el envío bloquea hasta que alguien recibe.
	trabajos := make(chan int)
	// Con búfer: acepta 20 envíos sin bloquear.
	resultados := make(chan string, 20)

	go productor(trabajos, 9)

	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go consumidor(i, trabajos, resultados, &wg)
	}

	// Cerrar resultados cuando todos los consumidores acaben.
	go func() {
		wg.Wait()
		close(resultados)
	}()

	for r := range resultados {
		fmt.Println(r)
	}
}
```

Ejecuta:

```powershell
go run ./cmd/dia16
```

**Reglas de canales que evitan el 90 % de los bloqueos:**
- Cierra siempre desde el productor, nunca desde el consumidor.
- Cerrar un canal ya cerrado, o escribir en uno cerrado, provoca panic.
- Leer de un canal cerrado devuelve el valor cero al instante; usa `v, ok := <-ch` para distinguirlo.
- Un canal nil bloquea para siempre — causa habitual de deadlock.

### Ejercicio

Implementa un *pipeline* de tres etapas: generar números del 1 al 20 → filtrar los pares → elevar al cuadrado. Cada etapa es una función que recibe un canal de entrada y devuelve uno de salida.

---

## Día 3 · `select` y `context`

**Conceptos:** `select`, timeouts, cancelación propagada.

Crea `cmd\dia17\main.go`:

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func tareaLenta(ctx context.Context, duracion time.Duration) (string, error) {
	select {
	case <-time.After(duracion):
		return "terminada", nil
	case <-ctx.Done():
		// ctx.Err() dice si fue cancelación o timeout.
		return "", ctx.Err()
	}
}

func main() {
	// select elige la primera operación lista. Si hay varias, al azar.
	a := make(chan string)
	b := make(chan string)
	go func() { time.Sleep(100 * time.Millisecond); a <- "de A" }()
	go func() { time.Sleep(200 * time.Millisecond); b <- "de B" }()

	for i := 0; i < 2; i++ {
		select {
		case m := <-a:
			fmt.Println(m)
		case m := <-b:
			fmt.Println(m)
		case <-time.After(time.Second):
			fmt.Println("timeout")
		}
	}

	// Contexto con timeout: se cancela solo a los 300 ms.
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel() // llamar a cancel SIEMPRE, o filtras recursos

	if _, err := tareaLenta(ctx, 1*time.Second); errors.Is(err, context.DeadlineExceeded) {
		fmt.Println("la tarea excedió el plazo")
	}

	ctx2, cancel2 := context.WithCancel(context.Background())
	go func() { time.Sleep(50 * time.Millisecond); cancel2() }()
	if _, err := tareaLenta(ctx2, time.Second); errors.Is(err, context.Canceled) {
		fmt.Println("la tarea fue cancelada")
	}
}
```

Ejecuta:

```powershell
go run ./cmd/dia17
```

**Por qué importa:** en tu API, `r.Context()` se cancela solo cuando el cliente cierra la conexión. Si pasas ese contexto a las consultas de base de datos, una petición abandonada deja de consumir recursos automáticamente.

### Ejercicio

Escribe `buscarEnParalelo(ctx, consultas []string) (string, error)` que lance varias búsquedas simuladas y devuelva la **primera** que responda, cancelando las demás. Añade un timeout global de 2 segundos.

---

## Día 4 · Postgres con Docker y `database/sql`

**Levanta la base de datos:**

```powershell
docker run --name pg-go -e POSTGRES_PASSWORD=secreto -e POSTGRES_DB=tareas -p 5432:5432 -d postgres:16
docker ps
```

**Instala el driver:**

```powershell
go get github.com/jackc/pgx/v5/stdlib
```

Crea `internal\db\db.go`:

```go
package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // import en blanco: registra el driver
)

// Conectar devuelve un POOL, no una conexión. Créalo una sola vez
// al arrancar y compártelo: sql.DB es seguro para uso concurrente.
func Conectar(ctx context.Context, dsn string) (*sql.DB, error) {
	pool, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("abriendo pool: %w", err)
	}

	pool.SetMaxOpenConns(25)
	pool.SetMaxIdleConns(25)
	pool.SetConnMaxLifetime(5 * time.Minute)

	// sql.Open no conecta de verdad; Ping sí.
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}
	return pool, nil
}

const esquema = `
CREATE TABLE IF NOT EXISTS tareas (
	id     SERIAL PRIMARY KEY,
	titulo TEXT NOT NULL,
	hecha  BOOLEAN NOT NULL DEFAULT FALSE,
	creada TIMESTAMPTZ NOT NULL DEFAULT NOW()
);`

func Migrar(ctx context.Context, pool *sql.DB) error {
	_, err := pool.ExecContext(ctx, esquema)
	return err
}
```

Crea `internal\db\tareas.go`:

```go
package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var ErrNoEncontrada = errors.New("tarea no encontrada")

type Tarea struct {
	ID     int       `json:"id"`
	Titulo string    `json:"titulo"`
	Hecha  bool      `json:"hecha"`
	Creada time.Time `json:"creada"`
}

type Repo struct{ pool *sql.DB }

func NuevoRepo(p *sql.DB) *Repo { return &Repo{pool: p} }

func (r *Repo) Crear(ctx context.Context, titulo string) (Tarea, error) {
	// $1 es un placeholder: el driver escapa el valor.
	// NUNCA construyas SQL concatenando cadenas.
	const q = `INSERT INTO tareas (titulo) VALUES ($1) RETURNING id, titulo, hecha, creada`
	var t Tarea
	err := r.pool.QueryRowContext(ctx, q, titulo).Scan(&t.ID, &t.Titulo, &t.Hecha, &t.Creada)
	if err != nil {
		return Tarea{}, fmt.Errorf("insertando tarea: %w", err)
	}
	return t, nil
}

func (r *Repo) Obtener(ctx context.Context, id int) (Tarea, error) {
	const q = `SELECT id, titulo, hecha, creada FROM tareas WHERE id = $1`
	var t Tarea
	err := r.pool.QueryRowContext(ctx, q, id).Scan(&t.ID, &t.Titulo, &t.Hecha, &t.Creada)
	if errors.Is(err, sql.ErrNoRows) {
		return Tarea{}, ErrNoEncontrada
	}
	if err != nil {
		return Tarea{}, fmt.Errorf("consultando tarea: %w", err)
	}
	return t, nil
}

func (r *Repo) Listar(ctx context.Context) ([]Tarea, error) {
	const q = `SELECT id, titulo, hecha, creada FROM tareas ORDER BY id`
	filas, err := r.pool.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("listando: %w", err)
	}
	defer filas.Close() // imprescindible: si no, filtras conexiones del pool

	out := []Tarea{}
	for filas.Next() {
		var t Tarea
		if err := filas.Scan(&t.ID, &t.Titulo, &t.Hecha, &t.Creada); err != nil {
			return nil, fmt.Errorf("escaneando: %w", err)
		}
		out = append(out, t)
	}
	// Comprobar el error del propio bucle: se olvida muy a menudo.
	return out, filas.Err()
}

func (r *Repo) Borrar(ctx context.Context, id int) error {
	res, err := r.pool.ExecContext(ctx, `DELETE FROM tareas WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("borrando: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNoEncontrada
	}
	return nil
}
```

Prueba en `cmd\dia18\main.go`:

```go
package main

import (
	"context"
	"fmt"
	"log"

	"aprendiendo-go/internal/db"
)

func main() {
	ctx := context.Background()
	dsn := "postgres://postgres:secreto@localhost:5432/tareas?sslmode=disable"

	pool, err := db.Conectar(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	if err := db.Migrar(ctx, pool); err != nil {
		log.Fatal(err)
	}

	repo := db.NuevoRepo(pool)
	t, err := repo.Crear(ctx, "primera tarea en Postgres")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("creada: %+v\n", t)

	todas, err := repo.Listar(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("total:", len(todas))
}
```

```powershell
go run ./cmd/dia18
```

### Ejercicio

Añade `Actualizar` y un `Buscar(ctx, texto string)` que use `WHERE titulo ILIKE '%' || $1 || '%'`. Añade también `ContarPorEstado` que devuelva un map con cuántas tareas hay hechas y pendientes.

---

## Día 5 · Transacciones

**Conceptos:** `BeginTx`, el patrón `defer tx.Rollback()`.

Añade a `internal\db\tareas.go`:

```go
// TransferirTitulos mueve el título de una tarea a otra de forma atómica:
// o se aplican los dos cambios, o ninguno.
func (r *Repo) MarcarTodasHechas(ctx context.Context, ids []int) error {
	tx, err := r.pool.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("iniciando transacción: %w", err)
	}
	// Rollback tras un Commit exitoso no hace nada, así que este defer
	// es seguro y garantiza que nunca dejas una transacción abierta.
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `UPDATE tareas SET hecha = TRUE WHERE id = $1`)
	if err != nil {
		return fmt.Errorf("preparando: %w", err)
	}
	defer stmt.Close()

	for _, id := range ids {
		res, err := stmt.ExecContext(ctx, id)
		if err != nil {
			return fmt.Errorf("actualizando %d: %w", id, err)
		}
		if n, _ := res.RowsAffected(); n == 0 {
			// Al devolver aquí, el defer revierte TODO lo anterior.
			return fmt.Errorf("id %d: %w", id, ErrNoEncontrada)
		}
	}
	return tx.Commit()
}
```

Pruébalo pasando una lista donde el último ID no exista y comprueba que ninguna tarea quedó marcada.

### Ejercicio

Escribe una transacción que cree una tarea y registre la operación en una tabla `auditoria (id, accion, tarea_id, momento)`. Fuerza un fallo a mitad y verifica en la base de datos que no quedó nada escrito. Puedes inspeccionar con:

```powershell
docker exec -it pg-go psql -U postgres -d tareas -c "SELECT * FROM tareas;"
```

---

## Día 6 · Worker en segundo plano

**Conceptos:** pool de workers, apagado ordenado de goroutines.

Crea `internal\worker\worker.go`:

```go
package worker

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

type Trabajo struct {
	ID     int
	TareaID int
	Accion string
}

type Pool struct {
	trabajos chan Trabajo
	wg       sync.WaitGroup
	log      *slog.Logger
}

func NuevoPool(n int, buffer int, log *slog.Logger) *Pool {
	p := &Pool{trabajos: make(chan Trabajo, buffer), log: log}
	for i := 1; i <= n; i++ {
		p.wg.Add(1)
		go p.trabajador(i)
	}
	return p
}

func (p *Pool) trabajador(id int) {
	defer p.wg.Done()
	for t := range p.trabajos { // termina cuando se cierra el canal
		p.log.Info("procesando", "trabajador", id, "trabajo", t.ID, "accion", t.Accion)
		time.Sleep(200 * time.Millisecond) // simula trabajo real
	}
	p.log.Info("trabajador detenido", "id", id)
}

// Encolar no bloquea si el búfer está lleno: descarta y avisa.
// Así una avalancha de peticiones no tumba la API.
func (p *Pool) Encolar(t Trabajo) bool {
	select {
	case p.trabajos <- t:
		return true
	default:
		p.log.Warn("cola llena, trabajo descartado", "trabajo", t.ID)
		return false
	}
}

// Parar cierra la cola y espera a que se vacíe, con límite de tiempo.
func (p *Pool) Parar(ctx context.Context) error {
	close(p.trabajos)
	hecho := make(chan struct{})
	go func() { p.wg.Wait(); close(hecho) }()

	select {
	case <-hecho:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
```

Pruébalo en `cmd\dia20\main.go`:

```go
package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"aprendiendo-go/internal/worker"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	p := worker.NuevoPool(3, 10, log)

	for i := 1; i <= 12; i++ {
		p.Encolar(worker.Trabajo{ID: i, TareaID: i, Accion: "notificar"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := p.Parar(ctx); err != nil {
		log.Error("apagado incompleto", "err", err)
	}
}
```

```powershell
go run -race ./cmd/dia20
```

### Ejercicio

Añade reintentos: si el procesamiento falla, reencola el trabajo hasta 3 veces con espera creciente (100 ms, 200 ms, 400 ms). Añade un contador de trabajos procesados y fallidos protegido con `sync/atomic`.

---

## Día 7 · Integrar todo

Reescribe la API para que use el repositorio de Postgres en lugar del almacén en memoria, y encole un trabajo al crear cada tarea.

Puntos clave de la integración:

```go
// En cada handler, pasa el contexto de la petición hacia abajo.
// Si el cliente se desconecta, la consulta se cancela sola.
func (api *API) crear(w http.ResponseWriter, r *http.Request) {
	var entrada struct{ Titulo string `json:"titulo"` }
	if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
		errorJSON(w, http.StatusBadRequest, "json inválido")
		return
	}

	t, err := api.repo.Crear(r.Context(), entrada.Titulo)
	if err != nil {
		api.log.Error("creando tarea", "err", err)
		errorJSON(w, http.StatusInternalServerError, "error interno")
		return
	}

	api.pool.Encolar(worker.Trabajo{TareaID: t.ID, Accion: "notificar"})
	escribirJSON(w, http.StatusCreated, t)
}
```

En `main`, el orden de apagado importa: primero el servidor HTTP (deja de aceptar peticiones), luego el pool de workers (termina lo pendiente), y por último la base de datos.

```go
srv.Shutdown(ctx)
pool.Parar(ctx)
poolDB.Close()
```

### Ejercicio final de la semana

1. El endpoint `/salud` debe hacer `pool.PingContext` y devolver 503 si la base de datos no responde.
2. Añade un `docker-compose.yml` que levante Postgres para no depender de recordar el comando `docker run`.
3. Instala `golang-migrate` y mueve el esquema a archivos de migración numerados en `migrations\`.
4. Ejecuta `go test -race ./...` con todo integrado y déjalo en verde.

---

## Comandos de referencia

| Comando | Qué hace |
|---|---|
| `docker start pg-go` | Reinicia la base de datos ya creada |
| `docker logs pg-go` | Ver por qué no arranca |
| `docker exec -it pg-go psql -U postgres -d tareas` | Consola SQL |
| `go test -race ./...` | Detector de carreras (úsalo siempre esta semana) |
| `go get paquete@latest` | Añadir dependencia |
| `go mod tidy` | Limpiar dependencias sin usar |

## Checklist antes de pasar a la semana 4

- [ ] Sé cuándo un canal bloquea y cómo evitar un deadlock
- [ ] Paso `r.Context()` a todas las consultas de base de datos
- [ ] Uso siempre placeholders `$1`, nunca concateno SQL
- [ ] Cierro `rows` con `defer` y compruebo `rows.Err()`
- [ ] Entiendo el patrón `defer tx.Rollback()` seguido de `tx.Commit()`
