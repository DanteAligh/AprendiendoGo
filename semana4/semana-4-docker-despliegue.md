# Semana 4 · Contenedor y despliegue en la nube

**Meta de la semana:** tu API corriendo en internet, en una imagen de menos de 20 MB, con CI que ejecuta los tests en cada push.

**Carpeta de trabajo:** raíz del proyecto (`Dockerfile`, `.github\workflows\`)

---

## Día 1 · Configuración por entorno

Las aplicaciones en contenedor no leen archivos de configuración: leen variables de entorno. Es el punto 3 de los "12 factores" y lo que esperan todas las plataformas.

Crea `internal\config\config.go`:

```go
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Puerto       string
	DSN          string
	Entorno      string
	NumWorkers   int
	TimeoutLect  time.Duration
	TimeoutEscr  time.Duration
}

func Cargar() (Config, error) {
	c := Config{
		Puerto:      texto("PUERTO", "8080"),
		DSN:         texto("DATABASE_URL", ""),
		Entorno:     texto("ENTORNO", "desarrollo"),
		NumWorkers:  entero("NUM_WORKERS", 3),
		TimeoutLect: 5 * time.Second,
		TimeoutEscr: 10 * time.Second,
	}
	// Fallar al arrancar es mejor que fallar en la primera petición.
	if c.DSN == "" {
		return Config{}, fmt.Errorf("falta la variable DATABASE_URL")
	}
	return c, nil
}

func texto(clave, porDefecto string) string {
	if v := os.Getenv(clave); v != "" {
		return v
	}
	return porDefecto
}

func entero(clave string, porDefecto int) int {
	v, err := strconv.Atoi(os.Getenv(clave))
	if err != nil {
		return porDefecto
	}
	return v
}
```

En `main`, sustituye toda lectura suelta de `os.Getenv` por `config.Cargar()`.

Para desarrollo local, crea `.env.example` (este sí se sube a git):

```
PUERTO=8080
DATABASE_URL=postgres://postgres:secreto@localhost:5432/tareas?sslmode=disable
ENTORNO=desarrollo
NUM_WORKERS=3
```

Y verifica que `.env` esté en tu `.gitignore`. **Nunca subas credenciales reales al repositorio.**

En PowerShell, para probar:

```powershell
$env:DATABASE_URL="postgres://postgres:secreto@localhost:5432/tareas?sslmode=disable"
go run ./cmd/api
```

### Ejercicio

Añade validación: si `ENTORNO=produccion`, exige que exista `API_KEY` y que el DSN no apunte a `localhost`. Escribe un test de tabla para `Cargar` usando `t.Setenv` (que restaura el valor al terminar el test).

---

## Día 2 · Compilación estática y primer Dockerfile

Go compila a un binario sin dependencias externas, lo que permite imágenes diminutas. La clave es `CGO_ENABLED=0`.

Crea `Dockerfile` en la raíz:

```dockerfile
# ---- Etapa de compilación ----
FROM golang:1.24-alpine AS build

WORKDIR /src

# Copiar primero los archivos de módulo aprovecha la caché de capas:
# si no cambian las dependencias, Docker no vuelve a descargarlas.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO_ENABLED=0 produce un binario estático sin dependencias de libc.
# -ldflags "-s -w" quita la tabla de símbolos: ~30% menos de tamaño.
RUN CGO_ENABLED=0 GOOS=linux go build \
	-ldflags="-s -w" \
	-o /bin/api ./cmd/api

# ---- Etapa final ----
# distroless/static no tiene shell, ni gestor de paquetes, ni nada:
# solo tu binario. Menos superficie de ataque, imagen mínima.
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /bin/api /api

USER nonroot:nonroot
EXPOSE 8080

ENTRYPOINT ["/api"]
```

Y `.dockerignore` (evita copiar basura al contexto de build):

```
.git
.env
bin/
*.md
*_test.go
.github
```

Construye y mira el tamaño:

```powershell
docker build -t tareas-api:local .
docker images tareas-api
```

Deberías ver alrededor de 15–20 MB. Compara con lo que ocuparía si la etapa final fuera `golang:1.24` (más de 800 MB).

### Ejercicio

Prueba tres variantes de la etapa final y anota los tamaños: `scratch`, `gcr.io/distroless/static-debian12` y `alpine:3.20`. Con Alpine podrás entrar con `docker run -it --entrypoint sh` para depurar; con las otras no. Piensa en qué escenario compensa cada una.

---

## Día 3 · Docker Compose

Un contenedor solo no sirve: necesitas la base de datos junto a la aplicación.

Crea `docker-compose.yml`:

```yaml
services:
  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_PASSWORD: secreto
      POSTGRES_DB: tareas
    ports:
      - "5432:5432"
    volumes:
      - datos-pg:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 3s
      retries: 5

  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      # Dentro de la red de compose, el host es el nombre del servicio.
      DATABASE_URL: postgres://postgres:secreto@db:5432/tareas?sslmode=disable
      PUERTO: "8080"
      ENTORNO: produccion
    depends_on:
      db:
        # Sin esto la API arranca antes que Postgres y falla al conectar.
        condition: service_healthy

volumes:
  datos-pg:
```

Levanta todo:

```powershell
docker compose up --build
```

En otra terminal:

```powershell
curl.exe http://localhost:8080/salud
curl.exe -X POST http://localhost:8080/tareas -H "Content-Type: application/json" -d "{\"titulo\":\"desde el contenedor\"}"
docker compose logs -f api
docker compose down
```

**El error clásico:** poner `localhost` en el `DATABASE_URL` de la API. Dentro de un contenedor, `localhost` es ese contenedor, no tu máquina. El host correcto es el nombre del servicio, `db`.

### Ejercicio

Añade un servicio `migrate` que use la imagen `migrate/migrate` para aplicar las migraciones antes de que arranque la API. Añade también un perfil de desarrollo que monte tu código y recompile al guardar (busca `air` o `reflex` para recarga en caliente).

---

## Día 4 · Observabilidad y `/salud`

Antes de desplegar, la aplicación tiene que poder decir si está sana.

```go
func (api *API) salud(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	estado := map[string]any{
		"estado":  "ok",
		"version": Version, // inyectada al compilar
		"hora":    time.Now().UTC(),
	}

	if err := api.pool.PingContext(ctx); err != nil {
		estado["estado"] = "degradado"
		estado["db"] = "sin conexión"
		// 503 le dice al balanceador que no le mande tráfico.
		escribirJSON(w, http.StatusServiceUnavailable, estado)
		return
	}
	estado["db"] = "ok"
	escribirJSON(w, http.StatusOK, estado)
}
```

Para inyectar la versión al compilar, declara en el paquete:

```go
var Version = "dev"
```

Y compila con:

```powershell
go build -ldflags="-X 'aprendiendo-go/internal/tareas.Version=1.0.0'" -o bin/api.exe ./cmd/api
```

En el Dockerfile puedes pasarlo como argumento:

```dockerfile
ARG VERSION=dev
RUN CGO_ENABLED=0 go build -ldflags="-s -w -X 'aprendiendo-go/internal/tareas.Version=${VERSION}'" -o /bin/api ./cmd/api
```

**Logs estructurados.** En producción los logs los lee una máquina, no tú. Usa `slog` con salida JSON:

```go
var log *slog.Logger
if cfg.Entorno == "produccion" {
	log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
} else {
	log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
}
slog.SetDefault(log)
```

Escribe siempre a `os.Stdout`, nunca a un archivo: el contenedor es efímero y la plataforma recoge la salida estándar.

### Ejercicio

Distingue dos endpoints: `/salud/vivo` (responde 200 siempre que el proceso esté en pie) y `/salud/listo` (comprueba base de datos y cola). Es la distinción entre *liveness* y *readiness* que usan Kubernetes y Cloud Run. Detén el contenedor de Postgres con `docker stop` y comprueba que `/salud/listo` pasa a 503 mientras `/salud/vivo` sigue en 200.

---

## Día 5 · Git y CI con GitHub Actions

Si aún no lo has hecho, sube el proyecto:

```powershell
git add .
git commit -m "API de tareas en Go"
git branch -M main
git remote add origin https://github.com/TU_USUARIO/aprendiendo-go.git
git push -u origin main
```

Crea `.github\workflows\ci.yml`:

```yaml
name: CI

on:
  push:
    branches: [main]
  pull_request:

jobs:
  test:
    runs-on: ubuntu-latest

    services:
      postgres:
        image: postgres:16
        env:
          POSTGRES_PASSWORD: secreto
          POSTGRES_DB: tareas_test
        ports:
          - 5432:5432
        options: >-
          --health-cmd pg_isready
          --health-interval 5s
          --health-timeout 3s
          --health-retries 5

    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
          cache: true

      - name: Formato
        run: test -z "$(gofmt -l .)"

      - name: Vet
        run: go vet ./...

      - name: Tests
        env:
          DATABASE_URL: postgres://postgres:secreto@localhost:5432/tareas_test?sslmode=disable
        run: go test -race -coverprofile=cobertura.out ./...

      - name: Compilar imagen
        run: docker build -t tareas-api:${{ github.sha }} .
```

Haz push y mira la pestaña Actions del repositorio. **Rompe un test a propósito**, haz push, y comprueba que el workflow falla. Esa es la señal de que la red de seguridad funciona.

### Ejercicio

Añade un paso que publique la imagen en GitHub Container Registry (`ghcr.io`) solo cuando el push sea a `main`. Necesitarás `docker/login-action` con `${{ secrets.GITHUB_TOKEN }}`. Añade también `golangci-lint` como paso de análisis estático.

---

## Día 6 · Despliegue

Elige una plataforma. Las tres funcionan bien con un contenedor de Go; **Fly.io** es la más directa para empezar.

**Opción A — Fly.io**

```powershell
# Instalar (PowerShell)
iwr https://fly.io/install.ps1 -useb | iex

fly auth signup
fly launch --no-deploy    # detecta el Dockerfile y genera fly.toml
fly postgres create       # crea la base de datos gestionada
fly postgres attach <nombre-db>   # inyecta DATABASE_URL automáticamente
fly deploy
fly logs
fly open
```

**Opción B — Google Cloud Run**

```powershell
gcloud auth login
gcloud run deploy tareas-api --source . --region us-central1 --allow-unauthenticated
```

Cloud Run compila desde el fuente y escala a cero cuando no hay tráfico. Para la base de datos usarías Cloud SQL o un Postgres gestionado externo como Neon o Supabase.

**Opción C — Railway**: conecta el repositorio de GitHub desde la web, añade un plugin de Postgres y despliega. Es la ruta con menos comandos.

**Detalle que rompe despliegues en Cloud Run:** la plataforma asigna el puerto por la variable `PORT`, no `PUERTO`. Asegúrate de que tu configuración lee ambas, o renombra la tuya:

```go
Puerto: texto("PORT", texto("PUERTO", "8080")),
```

Y escucha en `:8080` (todas las interfaces), no en `127.0.0.1:8080`, o el tráfico externo nunca llegará.

### Ejercicio

Despliega, y desde tu máquina ejecuta el CRUD completo contra la URL pública. Luego mira los logs en producción y confirma que salen en JSON con el ID de petición que añadiste en la semana 2.

---

## Día 7 · Cierre y repaso

**Documenta el proyecto.** Crea un `README.md` con: qué hace, cómo levantarlo en local (`docker compose up`), las variables de entorno necesarias, la lista de endpoints y la URL desplegada. Un README claro es lo primero que mira alguien que evalúa tu código.

**Revisa tu propio código con esta lista:**

- [ ] No hay credenciales en el repositorio (`git log -p | Select-String "password"` para comprobarlo)
- [ ] Todos los errores se manejan; no hay `_ = err` sin justificación
- [ ] Todo `rows` y `f` tiene su `defer Close()`
- [ ] Todas las consultas usan placeholders
- [ ] `go vet ./...` y `gofmt -l .` salen limpios
- [ ] `go test -race ./...` pasa
- [ ] El servidor tiene timeouts y apagado ordenado
- [ ] La imagen pesa menos de 25 MB

**Qué viene después.** Tres direcciones según lo que te interese:

*Profundizar en Go:* genéricos, `errgroup` para concurrencia con errores, `sync.Once` y `atomic`, perfilado con `pprof`, benchmarks con `go test -bench`.

*Backend serio:* autenticación con JWT, límite de peticiones por IP, OpenAPI para documentar, tests de integración con `testcontainers-go`, métricas con Prometheus y trazas con OpenTelemetry.

*Ecosistema:* gRPC con Protocol Buffers, `sqlc` para consultas tipadas, Kubernetes si vas hacia infraestructura, o herramientas de CLI con `cobra`.

### Ejercicio final

Añade una funcionalidad completa de principio a fin y llévala a producción: modelo, migración, repositorio, endpoints, tests, despliegue. Por ejemplo, etiquetas para las tareas con relación muchos a muchos. Ese ciclo completo, hecho solo y sin guía, es la verdadera prueba de que la semana 4 quedó aprendida.

---

## Comandos de referencia

| Comando | Qué hace |
|---|---|
| `docker build -t nombre .` | Construir la imagen |
| `docker compose up --build` | Levantar todo el stack |
| `docker compose down -v` | Bajar y borrar volúmenes |
| `docker exec -it contenedor sh` | Entrar (solo imágenes con shell) |
| `docker system prune -a` | Liberar espacio en disco |
| `fly deploy` / `fly logs` | Desplegar y ver logs |
| `go build -ldflags="-s -w"` | Binario reducido |
