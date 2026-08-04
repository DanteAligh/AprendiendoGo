# erp-api — Proyecto integrador (días 27-28)

Este es tu proyecto final del plan de 4 semanas: una API REST en Go, con
arquitectura de capas, para un mini-ERP con dos recursos relacionados:
**productos** y **facturas** (una factura contiene items que referencian
productos por ID, valida stock y descuenta inventario — el mismo dominio
que modelaste en `ejemplos/dia26_modelado_erp.go`, ahora expuesto por
HTTP con capas de verdad).

**No te doy este proyecto resuelto.** Lo que hay acá es un **andamiaje**:
estructura de carpetas, `go.mod`, modelos, interfaces, wiring de
middlewares y rutas — todo lo que necesitás para arrancar sin perder
tiempo en decisiones de "plomería". La lógica de negocio y los handlers
son tu trabajo. La guía paso a paso está en `../ejercicios.md`, sección
"Días 27-28".

## Qué SÍ está implementado (andamiaje)

- `cmd/api/main.go` — conecta todas las capas (repository -> service ->
  handler), registra las rutas, aplica los middlewares y levanta el
  servidor. **Este archivo ya funciona tal cual.**
- `internal/middleware/middleware.go` — logging, recuperación de pánicos
  (`recover`) y CORS. Igual que en `ejemplos/dia23_middleware_env.go`.
- `pkg/respuesta/respuesta.go` — helpers `respuesta.JSON(...)` y
  `respuesta.Error(...)` para responder JSON de forma consistente.
- `internal/models/*.go` — los structs `Producto`, `Factura` e
  `ItemFactura`, con sus tags JSON.
- Las **interfaces** de `internal/repository` y `internal/service`
  (los contratos que tus implementaciones deben cumplir).
- Los **constructores** de los repositorios en memoria
  (`NewRepositorioProductoMemoria`, `NewRepositorioFacturaMemoria`) — el
  mapa ya está inicializado, el mutex ya está declarado.

## Qué te toca implementar a vos

Todo lo marcado con `// TODO:` en:

- `internal/repository/producto_memoria.go` y `factura_memoria.go` —
  las operaciones CRUD reales sobre los mapas en memoria.
- `internal/service/producto_servicio.go` — validaciones de negocio de
  Producto (nombre no vacío, precio > 0, stock >= 0...).
- `internal/service/factura_servicio.go` — **el corazón del proyecto**:
  validar stock, congelar precios al momento de facturar, descontar
  inventario.
- `internal/handler/producto_handler.go` y `factura_handler.go` —
  decodificar JSON de entrada, llamar al service correspondiente,
  responder con el código de estado HTTP correcto.

Ahora mismo, si corrés el servidor, **compila y arranca perfectamente**,
pero cada endpoint de negocio responde `501 Not Implemented` — ese es tu
punto de partida verificable. Solo `GET /salud` responde de verdad, porque
no tiene lógica de negocio.

## Cómo correrlo

```bash
cd proyecto-final
go run ./semana4/proyecto-final/cmd/api
```

Por defecto escucha en el puerto 8080. Podés cambiarlo con la variable de
entorno `PUERTO`:

```bash
PUERTO=9090 go run ./semana4/proyecto-final/cmd/api
```

## Cómo probarlo mientras avanzás (ejemplos de curl)

Estos son los comandos que vas a usar para probar tu propia implementación
a medida que la vayas completando. El comportamiento descrito es el que
DEBERÍAS obtener una vez implementado — ahora mismo todos (salvo
`/salud`) te van a devolver `501`.

```bash
# Salud del servidor (ya funciona)
curl http://localhost:8080/salud

# Crear un producto
curl -X POST http://localhost:8080/productos \
  -H "Content-Type: application/json" \
  -d '{"nombre":"Teclado mecánico","precio":45.00,"stock":10}'

# Listar productos
curl http://localhost:8080/productos

# Obtener un producto por ID
curl http://localhost:8080/productos/1

# Actualizar un producto
curl -X PUT http://localhost:8080/productos/1 \
  -H "Content-Type: application/json" \
  -d '{"id":1,"nombre":"Teclado mecánico RGB","precio":55.00,"stock":8}'

# Eliminar un producto
curl -X DELETE http://localhost:8080/productos/1

# Crear una factura (usa productos que ya existan, y que tengan stock)
curl -X POST http://localhost:8080/facturas \
  -H "Content-Type: application/json" \
  -d '{
        "cliente_nombre": "Comercial Andina S.A.",
        "items": [
          {"producto_id": 1, "cantidad": 2},
          {"producto_id": 2, "cantidad": 1}
        ]
      }'

# Debería fallar con 409 si pedís más cantidad de la que hay en stock
curl -i -X POST http://localhost:8080/facturas \
  -H "Content-Type: application/json" \
  -d '{"cliente_nombre":"Cliente X","items":[{"producto_id":1,"cantidad":9999}]}'

# Listar facturas
curl http://localhost:8080/facturas

# Obtener una factura por ID
curl http://localhost:8080/facturas/1
```

## Después de terminarlo: hacia tu ERP + IA real

Este proyecto es deliberadamente pequeño, pero la arquitectura (capas,
middleware, inyección de dependencias por interfaces) es la MISMA que vas
a usar para un ERP real, mucho más grande. Los siguientes pasos naturales,
una vez que termines lo de arriba, son los "retos extra" que están en
`../ejercicios.md`: persistencia real con SQLite (día 20 de la semana 3 +
patrón de `go get` del día 24 de esta semana), autenticación con JWT (día
24), y un endpoint que llame a una API de IA real usando el patrón del día
25 (por ejemplo, generar automáticamente una descripción de producto o un
resumen de una factura). Ese es, literalmente, el siguiente proyecto que
vas a construir después de este plan de estudio.
