> **Banco de ejercicios extra (material importado).**
>
> Este archivo viene de otro plan de estudio y trae **muchos mas ejercicios**
> por concepto de los que tiene nuestro calendario. Usalo como cantera: cuando
> termines el ejercicio del dia y quieras mas repeticiones del mismo concepto,
> busca aqui ese concepto y haz dos o tres mas.
>
> **Dos avisos importantes:**
>
> 1. **La numeracion de dias de este archivo NO es la nuestra.** Aqui los
>    structs son "dia 8"; en nuestro calendario son el dia 4. Guiate por el
>    *concepto* (structs, maps, punteros...), nunca por el numero de dia.
> 2. **Da por sabida la logica basica de programacion.** Si una explicacion
>    supone que ya sabes que es un bucle, no te preocupes: eso lo cubre nuestro
>    material del dia correspondiente y CONCEPTOS-BASICOS.md.
>
> A proposito **no trae soluciones**. Esa es justo la idea.

---

# Ejercicios — Semana 4: Proyecto integrador

Este archivo es tu profesor de la semana. Leelo día por día, ANTES de
tocar código, y usá `ejemplos/` solo como referencia de sintaxis cuando te
atores — no como la solución a copiar. Los ejercicios de los días 22 a 26
no tienen solución escrita en ningún lado a propósito: la idea es que
llegues al proyecto integrador (días 27-28) habiendo peleado un poco con
cada pieza por tu cuenta.

---

## Día 22 — Arquitectura de un proyecto Go real

### Teoría

Hasta ahora escribiste programas de un archivo, o como mucho un par de
paquetes. Un ERP real tiene decenas (a veces cientos) de archivos, y con
el tiempo, varios desarrolladores tocándolos a la vez. Sin una estructura
clara, un proyecto así se convierte en lo que se conoce como "spaghetti
code": todo importa a todo, un cambio en una función rompe cosas en
lugares que no tienen relación aparente, y nadie se anima a tocar nada por
miedo a romper algo más.

Go no impone una estructura de carpetas obligatoria (a diferencia de
frameworks como Rails o Django), pero la comunidad convergió en una
convención muy estable, la que ya viste en `ejemplos/dia22_arquitectura/`:

- **`cmd/`** — un subdirectorio por cada binario ejecutable. Un ERP real
  típicamente tiene más de uno: la API HTTP, quizás un worker que procesa
  tareas en background, quizás una herramienta de línea de comandos para
  migraciones o importación de datos. Cada uno vive en su propio
  `cmd/nombre/main.go`.
- **`internal/`** — acá vive la lógica de tu aplicación. El compilador de
  Go IMPONE (no es solo convención) que nada fuera de tu módulo pueda
  importar paquetes bajo `internal/`. Esto evita acoplamientos accidentales.
- **`pkg/`** — código genérico, reutilizable, que en principio podría
  vivir en cualquier otro proyecto tuyo (un helper de respuestas HTTP, un
  cliente genérico...).
- **Capas dentro de `internal/`**: `handler` (HTTP) -> `service` (reglas
  de negocio) -> `repository` (persistencia). Cada capa solo conoce a la
  de abajo, nunca al revés, y nunca se salta capas (un handler no debería
  hablarle directo a un repository).

¿Por qué importa esto para vos? Porque en `proyecto-final/` vas a vivir
esta estructura de primera mano, y porque el día que empieces tu ERP real
con integración de IA, esta va a ser la primera decisión de diseño que
tomes.

### Ejercicios

1. Corré el ejemplo del día 22 (`cd ejemplos/dia22_arquitectura && go run
   ./cmd/servidor`) y modificá `internal/saludo/saludo.go` para que
   `Generar` también reciba una hora del día (mañana/tarde/noche) y ajuste
   el saludo ("Buenos días, ..." / "Buenas tardes, ..."). No toques
   `main.go` más de lo necesario para pasarle el nuevo parámetro.

2. Creá un SEGUNDO comando dentro del mismo módulo:
   `cmd/verificador/main.go`, que no levante ningún servidor HTTP, sino
   que simplemente llame a `saludo.Generar` con un par de valores fijos y
   los imprima por consola con `fmt.Println`. Esto te obliga a practicar
   la idea de "un módulo, varios binarios en `cmd/`". Corrélo con
   `go run ./cmd/verificador`.

3. Agregá un paquete nuevo bajo `internal/`, por ejemplo
   `internal/validacion`, con una función `EsNombreValido(nombre string)
   bool` que rechace nombres vacíos o que tengan más de 50 caracteres.
   Usala desde el handler de `cmd/servidor/main.go` antes de llamar a
   `saludo.Generar`: si el nombre no es válido, respondé con
   `http.StatusBadRequest` en vez de saludar.

4. Sin escribir código todavía: dibujá (en papel, en un comentario, o en
   un archivo de texto aparte) cómo organizarías `internal/` para un ERP
   con productos, clientes y facturas, usando la separación
   handler/service/repository. No hace falta que sea perfecto — la idea
   es que practiques pensar en capas ANTES de escribir el día 26 y el
   proyecto integrador.

5. **Reto extra (opcional):** Investigá qué pasa si intentás importar un
   paquete de `internal/` desde FUERA del módulo que lo contiene (por
   ejemplo, creando un módulo Go completamente aparte en otra carpeta e
   intentando `import "dia22_arquitectura/internal/saludo"`). Confirmá con
   tus propios ojos el error que da el compilador, y explicá con tus
   palabras por qué esa restricción es útil en un equipo grande.

---

## Día 23 — Middleware HTTP y variables de entorno

### Teoría

Repasá `ejemplos/dia23_middleware_env.go` con calma: un middleware es una
función `func(http.Handler) http.Handler`, es decir, una función que
envuelve un handler con comportamiento adicional sin tocar su código. Es
el mismo patrón "decorator" que probablemente ya conocés de otros
lenguajes, aplicado a HTTP.

En un ERP real vas a necesitar, como mínimo, tres middlewares desde el
primer día: logging (para poder auditar y depurar qué está pasando en
producción), recuperación de pánicos (para que un bug en un handler no
tumbe el servidor entero para todos los usuarios conectados) y, si tenés
un frontend separado, CORS.

Las variables de entorno son la forma estándar de configurar un programa
sin hardcodear valores (puertos, credenciales, claves de API) directamente
en el código — y sobre todo, sin subir secretos a tu repositorio git. Un
archivo `.env` es solo una conveniencia para no tener que exportar
variables a mano en cada sesión de terminal durante desarrollo.

### Ejercicios

1. Agregá un CUARTO middleware al ejemplo del día 23: uno que agregue un
   header de respuesta `X-Request-ID` con un identificador único por
   petición (podés generarlo con algo simple, como un contador atómico de
   `sync/atomic`, o con el paquete `crypto/rand` para algo más parecido a
   un UUID). Este patrón es real: en sistemas grandes, ese ID es lo que te
   permite rastrear una petición específica a través de logs de varios
   servicios distintos.

2. Modificá el middleware de logging para que también registre el código
   de estado HTTP de la respuesta (`200`, `404`, `500`...), no solo el
   método, la ruta y la duración. Pista: `http.ResponseWriter` no expone
   directamente el código de estado que se envió — vas a necesitar
   envolverlo en un wrapper propio que intercepte la llamada a
   `WriteHeader` (mirá cómo se resuelve exactamente este problema en
   `proyecto-final/internal/middleware/middleware.go`, que sí viene
   resuelto como andamiaje, DESPUÉS de intentarlo vos solo).

3. Creá un archivo `.env` de verdad junto al ejemplo del día 23, con al
   menos `PUERTO` y `ENTORNO`, y confirmá que `cargarEnvSimple` los carga
   correctamente (el servidor debería arrancar en el puerto que pusiste).
   Después, probá exportar `PUERTO` manualmente en tu shell ANTES de
   correr el programa (`export PUERTO=7000`) y confirmá que ese valor
   manual gana por encima del `.env` — es el comportamiento esperado y
   tenés que entender por qué `cargarEnvSimple` lo implementa así.

4. Si tenés conexión a internet, instalá `github.com/joho/godotenv`
   (`go get github.com/joho/godotenv`, en un módulo aparte para no romper
   el ejemplo original) y reescribí la carga de `.env` usando esa librería
   en vez de `cargarEnvSimple`. Compará cuánto código te ahorrás.

5. **Reto extra (opcional):** Escribí un middleware de "límite de tasa"
   (rate limiting) muy básico: si una misma IP hace más de N peticiones
   por segundo, respondé `429 Too Many Requests`. No hace falta que sea
   perfecto ni distribuido (con un mapa en memoria y un mutex alcanza para
   practicar la idea) — la versión de producción de esto normalmente usa
   Redis o un servicio dedicado, pero el concepto es el mismo.

---

## Día 24 — Autenticación: bcrypt y JWT

### Teoría

Repasá los comentarios de `ejemplos/dia24_auth_jwt_bcrypt.go` con
atención: ahí está explicado en detalle por qué bcrypt (hashing, no
cifrado, lento a propósito, con salt aleatorio) y qué es exactamente un
JWT (tres partes, la firma es lo que lo hace confiable, el contenido NO es
secreto).

Vale la pena remarcar la diferencia entre **autenticación** (¿quién sos?)
y **autorización** (¿qué podés hacer?). bcrypt + login resuelve la
autenticación. Un JWT que viaja en cada petición, más un middleware que lo
valida y expone el rol del usuario, es la base para resolver la
autorización (ej. "solo un admin puede eliminar productos").

### Ejercicios

1. Corré el ejemplo del día 24 y modificalo para que, en vez de un rol fijo
   ("admin"), la función `generarToken` reciba el rol como parámetro desde
   `main`, y generá DOS tokens distintos: uno para un usuario con rol
   `"empleado"` y otro con rol `"admin"`. Validá ambos e imprimí sus
   claims para confirmar que cada uno mantiene su propio rol.

2. Escribí una función `tieneRol(claims *ClaimsUsuario, rolRequerido
   string) bool` y usala para simular una autorización simple: antes de
   "ejecutar" una acción imaginaria como `eliminarProducto`, verificá que
   el rol del token sea `"admin"`; si no lo es, imprimí un mensaje de
   "acceso denegado" en vez de ejecutar la acción.

3. Cambiá la expiración del token a solo 2 segundos
   (`time.Now().Add(2 * time.Second)`), esperá 3 segundos con
   `time.Sleep`, e intentá validar el token vencido. Confirmá que
   `validarToken` devuelve un error, y fijate qué mensaje de error
   específico da la librería `golang-jwt` para el caso de "token
   expirado" (puede que necesites revisar la documentación del paquete
   `jwt` si el mensaje no es obvio).

4. Escribí (en un archivo nuevo, o agregando al mismo) una función
   `simularLogin(usuarios map[string]string, usuario, contrasena string)
   (string, error)` donde `usuarios` es un mapa de nombre de usuario a
   HASH de contraseña (ya hasheada con bcrypt de antemano). La función
   debe verificar la contraseña con `verificarContrasena` y, si es
   correcta, devolver un JWT generado con `generarToken`. Si es
   incorrecta, debe devolver un error genérico (por seguridad, un mensaje
   como "usuario o contraseña incorrectos", NUNCA "la contraseña es
   incorrecta" — eso le regalaría información a un atacante sobre si el
   usuario existe o no).

5. **Reto extra (opcional):** Investigá qué es un "refresh token" y por
   qué los sistemas de autenticación reales suelen emitir DOS tokens (uno
   de acceso, de corta duración, y uno de refresco, de larga duración) en
   vez de un único JWT de larga duración. Escribí, con tus palabras, en
   qué problema de seguridad ayuda ese diseño.

---

## Día 25 — Consumir una API externa desde Go

### Teoría

Repasá `ejemplos/dia25_cliente_http_externo.go`: el flujo (armar el
cuerpo JSON, crear el `*http.Request` con headers, mandarlo con un
`*http.Client` con timeout, revisar el código de estado, decodificar el
JSON de respuesta) es exactamente el mismo sin importar si hablás con una
API de pagos, de envío de emails, o con una API de un modelo de IA como
Claude o GPT. Dominar este patrón es, literalmente, la pieza que te separa
de poder integrar IA en tu ERP.

### Ejercicios

1. Modificá el servidor de prueba (`crearServidorDePruebaIA`) para que,
   además de `respuesta` y `modelo`, devuelva un campo `tokens_usados int`
   simulado (calculalo, por ejemplo, como la cantidad de palabras del
   prompt). Las APIs de IA reales casi siempre devuelven cuánto "costó" en
   tokens una petición — es información que vas a necesitar mostrar o
   loguear en un ERP real para controlar gastos.

2. Agregá un nuevo caso de prueba en `main` donde el servidor de prueba
   devuelva un código `500 Internal Server Error` a propósito (agregá una
   condición en el handler del servidor de prueba, por ejemplo "si el
   prompt contiene la palabra 'error', respondé 500"), y confirmá que
   `consultarAsistenteIA` maneja ese caso devolviendo un error legible sin
   hacer panic.

3. Agregá reintentos (retry) simples a `consultarAsistenteIA`: si la
   petición falla por un error de red (no por un 4xx/5xx, sino porque
   `cliente.Do` devuelve error), reintentá hasta 3 veces con una pequeña
   pausa entre intento e intento (`time.Sleep`). Este patrón es común
   cuando integrás con APIs externas que pueden tener fallos transitorios.

4. Cambiá el timeout del cliente HTTP a un valor muy corto (por ejemplo,
   `1 * time.Nanosecond`) y confirmá que `consultarAsistenteIA` devuelve
   un error de timeout en vez de colgarse esperando para siempre. Esto te
   confirma por qué un timeout es obligatorio en cualquier cliente HTTP de
   producción.

5. **Reto extra (opcional, requiere red):** Si tenés una cuenta y una
   clave de API de algún proveedor de IA real (Anthropic, OpenAI, o
   cualquier otro con una API HTTP documentada), adaptá
   `consultarAsistenteIA` para hablar con la API real en vez del servidor
   de prueba, respetando el formato de petición/respuesta específico de
   ese proveedor (leé su documentación). Guardá tu clave de API en una
   variable de entorno (día 23), NUNCA hardcodeada en el código.

---

## Día 26 — Modelado de un mini-ERP

### Teoría

Repasá `ejemplos/dia26_modelado_erp.go` con mucha atención — es el modelo
de dominio exacto que vas a reimplementar, ahora con capas y HTTP, en el
proyecto integrador. Los puntos clave a entender de memoria antes de
seguir:

- Las entidades se relacionan por ID (`ItemFactura.ProductoID`), no
  incrustando copias completas de otras entidades.
- Los cálculos de negocio (`Factura.Total()`) son métodos que operan sobre
  los datos ya cargados, sin tocar ninguna base de datos ni HTTP.
- Las validaciones de negocio (`ValidarStock`) se hacen ANTES de aplicar
  ningún cambio (`AplicarDescuentoStock`), nunca a medias.
- Los errores de negocio se declaran como variables de paquete
  (`ErrStockInsuficiente`, `ErrProductoNoExiste`) para poder distinguirlos
  con `errors.Is`, en vez de comparar strings de mensajes de error.

### Ejercicios

1. Agregá una nueva entidad, `Proveedor` (ID, Nombre, Email), y una
   relación: cada `Producto` ahora tiene un campo `ProveedorID int`.
   Escribí una función `ObtenerProveedor(producto Producto, proveedores
   map[int]Proveedor) (Proveedor, error)` que devuelva un error si el
   `ProveedorID` no existe en el mapa.

2. Escribí una función `CancelarFactura(factura *Factura, catalogo
   map[int]Producto) error` que haga lo inverso de `AplicarDescuentoStock`:
   le devuelva el stock a cada producto de los items de esa factura, y
   marque `factura.Confirmada = false`. Pensá: ¿debería poder cancelarse
   una factura que nunca se confirmó?

3. Agregá un campo `Descuento float64` (un porcentaje, ej. `0.10` para
   10%) a `Factura`, y modificá el método `Total()` para que aplique ese
   descuento sobre la suma de los items. Escribí también una validación:
   el descuento no puede ser negativo ni mayor a 1 (100%).

4. Escribí una función `ProductosConBajoStock(catalogo map[int]Producto,
   umbral int) []Producto` que devuelva todos los productos cuyo stock sea
   menor o igual al umbral dado. Esta es una funcionalidad real y muy
   común en cualquier ERP: alertas de reposición de inventario.

5. **Reto extra (opcional):** Escribí tests con `testing` (repasá la
   semana 3, día 21) para `ValidarStock` y para tu nueva
   `CancelarFactura`, cubriendo al menos: caso exitoso, producto
   inexistente, stock insuficiente, y (para `CancelarFactura`) intentar
   cancelar una factura no confirmada.

---

## Días 27-28 — Guía del proyecto integrador

Esta sección es distinta a las anteriores: en vez de una lista de
ejercicios sueltos, es una **guía paso a paso** para construir
`proyecto-final/`. No hay solución escrita en ningún lado — el andamiaje
(`go.mod`, modelos, interfaces, wiring) ya está armado y verificado de que
compila; tu trabajo es llenar cada `// TODO:` en el orden sugerido acá
abajo. Cada paso incluye cómo probarlo con `curl` ANTES de avanzar al
siguiente, para que nunca acumules más de una pieza sin probar a la vez.

Antes de empezar, leé `proyecto-final/README.md` completo — ahí está el
detalle exacto de qué archivos ya vienen resueltos y cuáles son tuyos.

### Paso 1 — Repositorio de productos (persistencia)

Andá a `proyecto-final/internal/repository/producto_memoria.go` e
implementá, en este orden, los métodos de `RepositorioProductoMemoria`:

1. `Crear` — asignar ID (`r.siguienteID`, incrementarlo), guardar en el
   mapa, devolver el producto con su ID ya puesto.
2. `Listar` — devolver todos los productos como un slice.
3. `ObtenerPorID` — devolver el producto o `ErrProductoNoEncontrado`.
4. `Actualizar` — reemplazar el producto existente, o
   `ErrProductoNoEncontrado` si el ID no existe.
5. `Eliminar` — borrar del mapa, o `ErrProductoNoEncontrado` si no existe.

No te olvides del `r.mu.Lock()` / `defer r.mu.Unlock()` en cada método que
lea o escriba el mapa.

**Cómo confirmar que este paso está bien** (sin todavía tocar service ni
handler): escribí un test rápido, o un `main` temporal, que llame
directamente a `NewRepositorioProductoMemoria()` y ejercite los 5 métodos.
No hace falta HTTP todavía para validar esta capa — esa es justamente la
ventaja de tenerla separada.

### Paso 2 — Service de productos (reglas de negocio)

Andá a `proyecto-final/internal/service/producto_servicio.go` e
implementá las validaciones antes de delegar al repositorio:

- `Crear`: nombre no vacío, precio mayor a 0, stock no negativo. Si algo
  falla, devolvé un error que envuelva `ErrProductoInvalido` con
  `fmt.Errorf("...: %w", ErrProductoInvalido)`.
- `Actualizar`: las mismas validaciones que `Crear`.
- `ObtenerPorID`, `Listar`, `Eliminar`: decidí vos qué validaciones (si
  alguna) tienen sentido acá.

### Paso 3 — Handlers de productos (HTTP)

Andá a `proyecto-final/internal/handler/producto_handler.go` e
implementá los 5 métodos, siguiendo los pasos sugeridos que ya están
escritos como comentarios en cada uno. Usá `respuesta.JSON` y
`respuesta.Error` de `pkg/respuesta` — no reinventes la codificación JSON.

**Cómo probar este paso:**

```powershell
cd proyecto-final && go run ./cmd/api
```

```powershell
curl -X POST http://localhost:8080/productos -H "Content-Type: application/json" \
  -d '{"nombre":"Teclado","precio":45,"stock":10}'
curl http://localhost:8080/productos
curl http://localhost:8080/productos/1
curl -X PUT http://localhost:8080/productos/1 -H "Content-Type: application/json" \
  -d '{"id":1,"nombre":"Teclado RGB","precio":55,"stock":8}'
curl -X DELETE http://localhost:8080/productos/1
```

No sigas al paso 4 hasta que los 5 endpoints de productos se comporten
como esperás, incluyendo los casos de error (probá pedir un producto con
ID inexistente y confirmá que te da 404, no un 500 ni un 200 vacío).

### Paso 4 — Repositorio de facturas

Mismo patrón que el paso 1, pero en
`proyecto-final/internal/repository/factura_memoria.go`: `Crear`,
`ObtenerPorID`, `Listar`. Más simple que el de productos porque no hay
`Actualizar` ni `Eliminar` para facturas en este proyecto (pensá por qué
podría tener sentido esa asimetría: ¿debería poder editarse una factura ya
confirmada?).

### Paso 5 — Service de facturas: la parte más importante

Este es el paso central de todo el proyecto integrador. En
`proyecto-final/internal/service/factura_servicio.go`, el método `Crear`
tiene que, en este orden exacto (los comentarios del archivo ya detallan
cada sub-paso):

1. Para cada item de la factura, buscar el producto correspondiente con
   `s.productoRepo.ObtenerPorID`. Si no existe, error que envuelva
   `ErrProductoNoExiste`.
2. Verificar que haya stock suficiente. Si no, error que envuelva
   `ErrStockInsuficiente`.
3. Recién si TODOS los items pasaron validación: fijar el
   `PrecioUnitario` de cada item al precio actual del producto, descontar
   stock de cada producto con `s.productoRepo.Actualizar`, y guardar la
   factura con `s.facturaRepo.Crear` marcada como `Confirmada = true`.

Repasá `ejemplos/dia26_modelado_erp.go` si te trabás con la lógica de
validación de stock — la idea es idéntica, solo que ahora coordinás dos
repositorios en vez de un solo mapa.

### Paso 6 — Handlers de facturas

En `proyecto-final/internal/handler/factura_handler.go`, implementá
`Listar`, `Crear` y `ObtenerPorID`. El más interesante es `Crear`: usá
`errors.Is` contra `service.ErrProductoNoExiste` y
`service.ErrStockInsuficiente` para decidir el código de estado HTTP de
cada caso de error (pensá cuál corresponde a cada uno — no hay un único
"500 para todo").

**Cómo probar este paso** (necesitás al menos un producto ya creado, del
paso 3):

```powershell
curl -X POST http://localhost:8080/productos -H "Content-Type: application/json" \
  -d '{"nombre":"Mouse","precio":25,"stock":3}'

curl -X POST http://localhost:8080/facturas -H "Content-Type: application/json" \
  -d '{"cliente_nombre":"Cliente Demo","items":[{"producto_id":1,"cantidad":2}]}'

# Este debería fallar por stock insuficiente (con el código que hayas elegido):
curl -i -X POST http://localhost:8080/facturas -H "Content-Type: application/json" \
  -d '{"cliente_nombre":"Cliente Demo","items":[{"producto_id":1,"cantidad":9999}]}'

curl http://localhost:8080/facturas
curl http://localhost:8080/facturas/1

# Confirmá que el stock del producto bajó después de facturar:
curl http://localhost:8080/productos/1
```

### Paso 7 — Retos extra (opcionales, para después de tener todo funcionando)

Una vez que los 6 pasos anteriores funcionan de punta a punta, estos son
caminos naturales para seguir profundizando, en dificultad creciente:

- **Tests**: escribí tests con `testing` (semana 3, día 21) para
  `factura_servicio.go`, cubriendo al menos: factura válida, producto
  inexistente, stock insuficiente.
- **Persistencia real**: escribí una implementación de
  `ProductoRepository` y `FacturaRepository` respaldada por SQLite
  (repasá semana 3, día 20, y el patrón de `go get` + `go.mod` del día 24
  de esta semana para manejar la dependencia externa
  `modernc.org/sqlite`, que es una implementación de SQLite en Go puro,
  sin necesitar un compilador de C). Cambiá SOLO `cmd/api/main.go` para
  instanciar la nueva implementación — si tuviste que tocar
  `internal/service` o `internal/handler` para que esto funcione, algo
  del diseño de interfaces no quedó bien separado.
- **Autenticación**: agregá un endpoint `POST /login` que reciba
  usuario/contraseña, valide contra un usuario hardcodeado (con su
  contraseña ya hasheada con bcrypt, día 24) y devuelva un JWT. Agregá un
  middleware nuevo que valide ese JWT en las rutas de escritura
  (`POST`/`PUT`/`DELETE`), dejando las de lectura (`GET`) sin protección.
- **Integración con IA**: agregá un endpoint, por ejemplo
  `POST /facturas/{id}/resumen`, que arme un prompt con los datos de la
  factura y llame a una API de IA (real, si tenés clave de API, o
  simulada con `httptest` como en el día 25) para generar un resumen en
  lenguaje natural de esa factura. Este último reto es, literalmente, el
  primer paso de integración de IA en un ERP — el punto exacto donde
  termina este plan de estudio y empieza tu próximo proyecto.
