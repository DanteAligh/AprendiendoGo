# Anexo — Cómo investigar por tu cuenta (semana 3: base de datos y tests)

Esta semana el código deja de vivir solo en memoria: hablas con una base de
datos de verdad y escribes tests que la prueban. Dos áreas donde la
documentación que necesitas no siempre es la de Go.

## `database/sql`: la documentación del driver, no solo de Go

`database/sql` es una interfaz genérica — el comportamiento específico
(cómo conectarse, qué formato de connection string usa) depende del
**driver** (`modernc.org/sqlite`, `github.com/lib/pq` para Postgres, etc).
Cuando algo falle en la conexión, la documentación que necesitas es la del
driver en pkg.go.dev, no la de `database/sql`. Búscalo como
`pkg.go.dev/modernc.org/sqlite` (o el driver que estés usando).

## Testing: `go test` tiene flags que valen la pena conocer

```powershell
go test ./...              # corre todos los tests del módulo
go test -v ./...           # modo detallado: ves cada test por nombre
go test -run NombreTest     # corre solo un test específico
go test -cover ./...       # te dice qué % de tu código está cubierto
```

`go doc testing` te muestra el paquete completo, pero en la práctica vas
a aprender más leyendo tests reales de la librería estándar en GitHub que
la documentación sola — el patrón table-driven test se entiende mejor
viéndolo que leyéndolo.

## Regla de esta semana

Si el problema es de un **driver** externo, busca su documentación específica
en pkg.go.dev — no asumas que se comporta igual que otro driver que ya
conoces. Y antes de dar por bueno código nuevo, córrelo con `go test -race ./...`.
