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

# Ejercicios Semana 3: Concurrencia profunda + Backend

Este archivo es tu profesor. Cada día tiene: una explicación de fondo, ejercicios
progresivos (sin código de solución, para que pienses) y un reto extra opcional.
Los archivos en `ejemplos/` son referencia de sintaxis — úsalos para consultar "¿cómo se
escribe esto?", no para copiar la solución del ejercicio.

---

## Día 15: sync.WaitGroup, sync.Mutex, race conditions, select, sync.Once

### La explicación

En la semana 2 viste goroutines y channels básicos: lanzar `go algo()` y comunicarte con
`chan`. Pero en un backend real casi nunca lanzas *una* goroutine y esperas — lanzas
*muchas* (una por petición HTTP, una por tarea en segundo plano) y necesitas dos cosas que
Go resuelve con herramientas distintas y muy deliberadas:

1. **"¿Ya terminaron todas?"** → `sync.WaitGroup`. Es un contador: `Add(n)` suma trabajo
   pendiente, cada goroutine llama `Done()` al terminar (típicamente con `defer`), y
   `Wait()` bloquea hasta que el contador llega a cero. No hay canal de por medio: es la
   herramienta correcta cuando lo único que te importa es "espera a que todos terminen",
   no intercambiar datos.

2. **"¿Cómo evito que dos goroutines pisen el mismo dato al mismo tiempo?"** → `sync.Mutex`.
   Una **race condition** ocurre cuando dos o más goroutines leen y escriben la misma
   variable sin coordinación, y el resultado final depende del orden impredecible en que el
   sistema operativo las ejecuta. Por ejemplo, `contador++` no es una operación atómica:
   internamente es "leer, sumar 1, escribir", y si dos goroutines lo hacen "al mismo
   tiempo" pueden pisarse y perder incrementos. Un `Mutex` (mutual exclusion) garantiza que
   solo una goroutine a la vez ejecuta el bloque protegido entre `Lock()` y `Unlock()`.
   Go además trae una herramienta para *detectar* estos bugs automáticamente: el
   **race detector** (`go run -race archivo.go`), que instrumenta tu programa y avisa en
   tiempo de ejecución si detectó un acceso no sincronizado.

   ¿Por qué Go no evita esto "por diseño" como otros lenguajes intentan hacer con
   inmutabilidad forzada? Porque la filosofía de Go es dar herramientas simples y explícitas
   (goroutines baratas, channels, mutex) y confiar en que el desarrollador declare
   intención clara de sincronización — y te da el race detector para cuando te equivocas.

3. **"Tengo varios channels, ¿de cuál leo primero?"** → `select`. Es como un `switch` pero
   para operaciones de channel: espera a que *cualquiera* de varios `case` (recibir o
   enviar en un channel) esté listo, y ejecuta ese. Si ninguno está listo y hay un
   `default`, lo ejecuta sin bloquear. Esto es fundamental para patrones como "timeout de
   una operación" o "atender la primera de varias fuentes que responda".

4. **"Necesito que algo se ejecute UNA sola vez, sin importar cuántas goroutines lo pidan"**
   → `sync.Once`. Típico para inicializar una conexión compartida (a una base de datos, a
   una caché) de forma perezosa y seguridad ante condiciones de carrera: sin importar
   cuántas goroutines llamen `once.Do(inicializar)` concurrentemente, `inicializar` corre
   una sola vez.

### En un backend real (piensa en tu ERP)

Imagina un endpoint que actualiza el stock de un producto. Si dos peticiones HTTP llegan
"al mismo tiempo" (dos goroutines, porque `net/http` atiende cada petición en su propia
goroutine) y ambas hacen "leer stock, restar 1, guardar stock" sin un mutex, puedes vender
más unidades de las que existen. Esto no es un caso hipotético — es el bug de concurrencia
más común en sistemas de inventario/ERP reales. `sync.WaitGroup` lo vas a usar para tareas
como "lanza 5 reportes en paralelo y espera a que todos terminen antes de responder".

### Ejercicios

1. Escribe un programa que lance 10 goroutines, cada una imprimiendo su número de goroutine
   (0 a 9), y usa `sync.WaitGroup` para que el programa principal espere a que las 10
   terminen antes de imprimir "listo". Pista: el número de goroutine debe pasarse como
   parámetro a la función lanzada, no capturado directamente de la variable del `for`.

2. Crea un `contador` compartido (un `int`) y lanza 100 goroutines que cada una lo
   incremente 1000 veces, **sin** mutex. Corre el programa varias veces y observa que el
   resultado final rara vez es 100000. Luego corre el mismo programa con
   `go run -race archivo.go` y observa el reporte de race condition.

3. Ahora arregla el ejercicio anterior agregando un `sync.Mutex` que proteja el
   incremento. Verifica que ahora el resultado siempre sea exactamente 100000 y que
   `-race` ya no reporte nada.

4. Escribe una función que reciba dos channels (`chan string`) que simulan dos fuentes de
   datos, y un `time.After` como tercer caso. Usa `select` para imprimir el primer mensaje
   que llegue de cualquiera de los dos channels, o imprimir "timeout" si pasan más de 2
   segundos sin que llegue nada.

5. Simula una "conexión compartida" con una struct que tiene un campo `conectado bool` y un
   método `Conectar()` que imprime "conectando..." y pone `conectado = true`. Usa
   `sync.Once` para que, aunque 5 goroutines llamen `Conectar()` concurrentemente a través
   de `once.Do(...)`, el mensaje "conectando..." se imprima una sola vez.

💡 **Reto extra:** combina `select` con un `context.Context` (revísalo brevemente en la
documentación estándar) para cancelar una goroutine que "trabaja" en un loop infinito
cuando el contexto se cancela o vence su timeout. Es el patrón real que vas a usar para
cancelar operaciones HTTP de larga duración.

---

## Día 16: Manejo de archivos con os, io, bufio y encoding/csv

### La explicación

Go no tiene "excepciones" para errores de I/O (abrir un archivo que no existe, escribir en
un disco lleno) por la misma razón que no las tiene para nada más: los errores son
**valores de retorno explícitos**, no un flujo de control invisible. Cuando haces
`os.Open("archivo.txt")`, obtienes `(*os.File, error)`, y **tú decides** qué hacer si el
error no es `nil`. Esto es intencional: en un backend, un archivo que no existe no es algo
"excepcional" — es un caso esperado que tu código debe manejar explícitamente (¿reintentar?
¿crear el archivo? ¿responder 404?).

Las piezas clave:

- **`os`**: abrir (`os.Open`), crear (`os.Create`), y operaciones a bajo nivel sobre
  archivos y el sistema operativo. `os.Open` es solo lectura; `os.Create` trunca/crea para
  escritura.
- **`io`**: las interfaces fundamentales, sobre todo `io.Reader` y `io.Writer`. Esto es
  *el* diseño elegante de Go: un archivo, una conexión de red, un buffer en memoria, la
  entrada estándar... todos implementan `io.Reader`/`io.Writer`. Si tu función acepta un
  `io.Reader`, funciona igual con un archivo que con una respuesta HTTP.
- **`bufio`**: envuelve un `io.Reader`/`io.Writer` con un buffer, para no hacer una llamada
  al sistema operativo por cada byte. `bufio.Scanner` es la forma idiomática de leer un
  archivo línea por línea; `bufio.Writer` acumula escrituras y las manda en bloque (por
  eso necesitas `Flush()`).
- **`encoding/csv`**: CSV es *el* formato con el que cualquier ERP intercambia datos con
  Excel, otros sistemas o cargas masivas de clientes. El paquete te da un `Reader` que
  parsea línea por línea respetando comillas y comas dentro de campos, y un `Writer`
  simétrico.

### En un backend real

Un ERP real necesita: exportar un reporte de ventas a CSV para que el cliente lo abra en
Excel, leer un archivo de carga masiva de productos que te manda un proveedor, escribir
logs de auditoría a disco. Todo esto pasa por `os`/`io`/`bufio`/`encoding/csv` antes de que
existiera cualquier base de datos.

### Ejercicios

1. Escribe un programa que cree un archivo `notas.txt` y escriba 5 líneas de texto (usa
   `os.Create` y luego `bufio.NewWriter` o escritura directa con `file.WriteString`).
   Verifica con `defer file.Close()`.

2. Ahora lee ese mismo archivo línea por línea con `bufio.NewScanner` y cuenta cuántas
   líneas tiene. Imprime cada línea con su número (1, 2, 3...).

3. Crea un archivo `productos.csv` con encabezado `nombre,precio,stock` y al menos 4 filas
   de datos, usando `encoding/csv.Writer`. Pista: `Write` recibe `[]string`, y no olvides
   `w.Flush()` al final.

4. Lee `productos.csv` con `encoding/csv.Reader`, salta el encabezado, y calcula el valor
   total del inventario (suma de `precio * stock` de todas las filas). Vas a necesitar
   convertir strings a números con `strconv`.

5. Escribe una función que reciba un `io.Reader` (no un `*os.File` directamente) y cuente
   cuántas palabras tiene el texto. Pruébala primero pasándole un archivo abierto, y luego
   pasándole un `strings.NewReader("texto de prueba")` — sin cambiar la función. Esto
   demuestra por qué programar contra la interfaz `io.Reader` es más flexible que programar
   contra `*os.File`.

💡 **Reto extra:** escribe un programa que lea `productos.csv` y genere un segundo archivo
`productos_bajo_stock.csv` solo con las filas donde `stock < 10`, reutilizando
`encoding/csv.Writer`.

---

## Día 17: JSON con encoding/json

### La explicación

JSON es el idioma universal de las APIs backend modernas. Go lo maneja con el paquete
estándar `encoding/json` (sin librerías externas) usando **reflection** para mapear
structs a JSON y viceversa:

- **`json.Marshal(valor)`** convierte un valor Go a `[]byte` de JSON.
- **`json.Unmarshal(datos, &valor)`** hace lo inverso: parsea JSON hacia una struct Go (por
  eso siempre pasas un puntero — Unmarshal necesita escribir en tu variable).
- **Struct tags** (`` `json:"nombre"` ``) le dicen a `encoding/json` cómo se llama el campo
  en el JSON, porque en Go los campos exportados empiezan con mayúscula (`Nombre`) pero en
  JSON casi siempre usamos minúsculas/snake_case (`"nombre"`). Sin el tag, Go usaría el
  nombre del campo tal cual.
- **`omitempty`** en el tag (`` `json:"telefono,omitempty"` ``) le dice a Marshal que omita
  el campo del JSON de salida si tiene el valor cero de su tipo (string vacío, 0, nil,
  slice vacío...). Es clave para APIs reales donde no todos los campos son obligatorios.
- **JSON anidado**: un struct puede tener otro struct (o slice de structs) como campo, y
  `encoding/json` lo serializa/deserializa recursivamente sin que tengas que hacer nada
  especial — solo refleja la jerarquía en tus structs Go.

¿Por qué Go no necesita un framework para esto (a diferencia de otros lenguajes que a veces
requieren librerías de terceros para JSON básico)? Porque el diseño de Go prioriza que la
biblioteca estándar cubra lo esencial de un backend moderno: JSON y HTTP vienen "de fábrica"
precisamente para que puedas escribir una API sin instalar nada.

### En un backend real

Cada vez que tu API reciba un `POST` con un body JSON (crear un producto, una factura, un
usuario) vas a hacer `Unmarshal` hacia una struct. Cada vez que respondas, vas a hacer
`Marshal` de una struct hacia JSON. `omitempty` te va a servir constantemente para
respuestas de API donde campos opcionales no deben aparecer si están vacíos.

### Ejercicios

1. Define una struct `Producto` con campos `Nombre` (string), `Precio` (float64) y `Stock`
   (int), con struct tags en minúsculas (`json:"nombre"`, etc). Crea una instancia y
   conviértela a JSON con `Marshal`, imprimiendo el resultado como string.

2. Toma un string JSON como `` `{"nombre":"Teclado","precio":450.50,"stock":12}` `` y
   conviértelo de vuelta a una struct `Producto` con `Unmarshal`. Imprime los campos por
   separado para confirmar que se parsearon bien.

3. Agrega un campo `Descripcion string` con tag `` `json:"descripcion,omitempty"` `` a
   `Producto`. Crea dos instancias: una con descripción y otra sin ella (string vacío).
   Convierte ambas a JSON y observa que en la segunda el campo `descripcion` no aparece en
   absoluto en el resultado.

4. Define una struct `Pedido` que tenga un campo `Cliente string` y un campo
   `Productos []Producto` (slice anidado). Crea un pedido con 2-3 productos y conviértelo a
   JSON. Verifica visualmente que el JSON resultante tiene un arreglo anidado de objetos.

5. Usa `json.MarshalIndent` en vez de `json.Marshal` para el ejercicio anterior, y compara
   la salida. ¿Cuándo usarías uno u otro en un backend real (piensa en logs vs. respuestas
   de API de producción)?

💡 **Reto extra:** define una struct con un campo que sea un `map[string]interface{}` para
representar datos "flexibles" cuya forma no conoces de antemano (por ejemplo, metadata
arbitraria de un producto), y experimenta serializándola y deserializándola.

---

## Día 18: net/http básico — tu primer servidor

### La explicación

Un servidor HTTP es, en el fondo, un programa que escucha un puerto de red, recibe
peticiones con un método (GET, POST...) y una ruta (`/productos`), y devuelve una
respuesta. Go trae esto en la biblioteca estándar porque su caso de uso original (Google)
es infraestructura de servidores — por eso `net/http` es sorprendentemente completo sin
necesitar ningún framework externo para empezar (a diferencia de otros lenguajes donde
casi siempre instalas algo tipo Express/Flask desde el día uno).

Piezas clave:

- **`http.HandleFunc(ruta, funcion)`**: registra una función que se ejecuta cuando llega
  una petición a esa ruta. La función tiene la firma
  `func(w http.ResponseWriter, r *http.Request)`.
- **`http.ResponseWriter`**: es por donde *escribes* la respuesta (con `w.Write(...)`,
  `fmt.Fprintf(w, ...)`, o `w.WriteHeader(codigo)` para el código de estado HTTP como 200,
  404, 500).
- **`*http.Request`**: contiene todo sobre la petición entrante: `r.Method` (GET/POST...),
  `r.URL.Query()` para leer query params (`?nombre=valor`), y `r.Body` (un `io.Reader`) para
  leer el cuerpo de la petición, típicamente con `json.NewDecoder(r.Body).Decode(&struct)`.
- **`http.ListenAndServe(":8080", nil)`**: arranca el servidor y bloquea escuchando en el
  puerto indicado.

Un detalle importante: **cada petición entrante se atiende en su propia goroutine**,
automáticamente, sin que tú escribas `go`. Esto es exactamente por qué el día 15
(mutex, race conditions) importa tanto para HTTP: si dos peticiones tocan el mismo dato
compartido en memoria, tienes una race condition real.

### En un backend real

Todo lo que vas a construir de aquí en adelante — tu ERP incluido — empieza como un
servidor HTTP escuchando peticiones. Frameworks como Gin o Echo (que verás más adelante en
tu camino) son solo capas de conveniencia sobre exactamente estos mismos conceptos de
`net/http`.

### Ejercicios

1. Crea un servidor que responda `"Hola, backend"` (texto plano) en la ruta `/` para
   peticiones GET. Arráncalo en el puerto 8080 y pruébalo con
   `curl http://localhost:8080/`.

2. Agrega una ruta `/saludo` que lea un query param `nombre`
   (`GET /saludo?nombre=Carlos`) y responda `"Hola, Carlos"`. Si el query param viene
   vacío, responde `"Hola, desconocido"`. Prueba con
   `curl "http://localhost:8080/saludo?nombre=Carlos"`.

3. Agrega una ruta `/eco` que solo acepte `POST`, lea el body completo de la petición
   (texto plano, no JSON todavía) y lo devuelva tal cual en la respuesta. Si el método no
   es POST, responde con código 405 (`http.StatusMethodNotAllowed`) y un mensaje de error.
   Prueba con `curl -X POST -d "hola mundo" http://localhost:8080/eco`.

4. Crea una struct `Saludo` con un campo `Mensaje string` (con su tag JSON). Agrega una
   ruta `/saludo-json` que responda ese struct como JSON usando
   `json.NewEncoder(w).Encode(...)`, y asegúrate de poner el header
   `w.Header().Set("Content-Type", "application/json")` antes de escribir. Verifica con
   `curl http://localhost:8080/saludo-json`.

5. Agrega una ruta `/estado` que reciba un `POST` con body JSON
   `{"nombre": "Ana", "edad": 30}`, lo decodifique a una struct, e imprima en la consola del
   servidor (no en la respuesta) los valores recibidos, respondiendo simplemente
   `"recibido"` con código 200. Prueba con
   `curl -X POST -d '{"nombre":"Ana","edad":30}' http://localhost:8080/estado`.

💡 **Reto extra:** agrega manejo de rutas no encontradas: si alguien pide una ruta que no
registraste, ¿qué responde tu servidor por defecto? Investiga cómo usar un
`http.NewServeMux()` explícito en vez de depender del mux por defecto (`http.DefaultServeMux`),
y explica en un comentario por qué en un backend real conviene ser explícito.

---

## Día 19: API REST en memoria — tu primer CRUD

### La explicación

Hoy juntas todo: JSON (día 17) + HTTP (día 18) + una estructura de datos en memoria
(slice o map, semana 1) para construir un **CRUD** completo (Create, Read, Update, Delete)
sobre un recurso, sin base de datos todavía. Esto es intencional: separar "cómo expongo un
recurso por HTTP" de "dónde lo guardo" es exactamente el mismo patrón que vas a usar cuando
el día 20 cambies el almacenamiento en memoria por SQLite — la parte HTTP casi no cambia.

El patrón REST estándar para un recurso (usemos `tareas`) es:

| Método | Ruta            | Acción                          |
|--------|-----------------|----------------------------------|
| GET    | `/tareas`       | listar todas las tareas          |
| GET    | `/tareas/{id}`  | obtener una tarea por id         |
| POST   | `/tareas`       | crear una tarea (body JSON)      |
| PUT    | `/tareas/{id}`  | actualizar una tarea (body JSON) |
| DELETE | `/tareas/{id}`  | eliminar una tarea               |

Como el mux estándar de Go (antes de Go 1.22) no distingue fácilmente `/tareas` de
`/tareas/{id}` con método, es común revisar `r.Method` y parsear el id manualmente del
`r.URL.Path` (o usar los patrones con método incluidos desde Go 1.22, tipo
`"GET /tareas/{id}"`, que revisarás en el ejemplo del día 18/19).

Como este almacenamiento en memoria se comparte entre peticiones (goroutines distintas),
**necesitas el `sync.Mutex` del día 15** para proteger el map/slice contra escrituras
concurrentes. Este es el primer ejercicio donde de verdad ves *todo* conectado.

### En un backend real

Este es, literalmente, el esqueleto de cada módulo de un ERP: el recurso "producto" o
"factura" o "empleado" expuesto como API REST. La única diferencia entre este ejercicio y
un módulo real de producción es dónde vive el dato (memoria vs. base de datos) — que es
justo lo que resuelves mañana.

### Ejercicios

1. Define una struct `Tarea` con `ID int`, `Titulo string`, `Completada bool` (con tags
   JSON). Crea un almacenamiento en memoria (`map[int]Tarea` protegido por un
   `sync.Mutex`, o un slice) con 2-3 tareas iniciales.

2. Implementa `GET /tareas` que devuelva **todas** las tareas como un arreglo JSON.
   Respuesta esperada (ejemplo):
   ```json
   [{"id":1,"titulo":"Comprar leche","completada":false}]
   ```

3. Implementa `POST /tareas` que reciba un body JSON como
   `{"titulo":"Nueva tarea"}`, le asigne un `id` nuevo automáticamente (el siguiente
   disponible), la guarde, y responda con código 201 (`http.StatusCreated`) y la tarea
   creada completa en JSON (incluyendo su nuevo `id`).

4. Implementa `GET /tareas/{id}` (obtener una) y `DELETE /tareas/{id}` (eliminar una).
   Para `GET`, si el id no existe, responde 404 con un mensaje de error en JSON, ej.
   `{"error":"tarea no encontrada"}`. Para `DELETE`, responde 200 con
   `{"mensaje":"tarea eliminada"}` o 404 si no existía.

5. Implementa `PUT /tareas/{id}` que reciba un body JSON parcial o completo (ej.
   `{"titulo":"Actualizado","completada":true}`) y actualice la tarea existente. Si el id
   no existe, responde 404.

💡 **Reto extra:** agrega validación básica: si `POST /tareas` llega sin `titulo` o con
`titulo` vacío, responde 400 (`http.StatusBadRequest`) con un mensaje de error claro en
JSON, en vez de crear una tarea inválida.

---

## Día 20: Bases de datos con database/sql + SQLite

### La explicación

`database/sql` es la interfaz **estándar** de Go para hablar con bases de datos
relacionales — MySQL, PostgreSQL, SQLite, SQL Server, todos usan la misma API
(`sql.DB`, `Query`, `Exec`, `Scan`...). Lo que cambia entre motores es el **driver**, una
librería que implementa esa interfaz para un motor específico y se registra con
`sql.Open("nombre_driver", cadena_de_conexion)`.

Para este ejercicio usamos **`modernc.org/sqlite`**, un driver de SQLite escrito 100% en
Go (sin necesitar un compilador de C ni CGO), lo que lo hace mucho más simple de instalar
y portar entre máquinas — ideal para aprender y para desplegar sin dolores de compilación.

Conceptos clave:

- **`sql.Open`** no conecta de inmediato; solo prepara el objeto `*sql.DB` (que es en
  realidad un *pool* de conexiones, seguro para usar concurrentemente desde varias
  goroutines — otra razón por la que el día 15 importa).
- **`db.Exec(sql, args...)`** para sentencias que no devuelven filas: `CREATE TABLE`,
  `INSERT`, `UPDATE`, `DELETE`.
- **`db.Query(sql, args...)`** para `SELECT` que devuelve múltiples filas — iteras con
  `rows.Next()` y extraes con `rows.Scan(&campo1, &campo2, ...)`.
- **`db.QueryRow(sql, args...)`** para cuando esperas una sola fila (ej. buscar por id).
- **Placeholders (`?`)**: en vez de concatenar strings SQL con datos del usuario (lo que
  abre la puerta a **SQL injection**), usas `?` como marcador de posición y pasas los
  valores reales como argumentos aparte. El driver se encarga de escaparlos de forma
  segura. Esto es un **prepared statement**: la sentencia SQL se compila una vez y se
  reutiliza con distintos valores, lo cual es tanto más seguro como más eficiente.

### En un backend real

Todo dato que tu ERP necesita recordar después de reiniciar el proceso vive en una base de
datos. El patrón que practicas hoy (abrir conexión, crear tabla, `Exec`/`Query` con
placeholders) es exactamente el mismo que usarás con PostgreSQL en producción — solo
cambia el driver y algunos detalles de sintaxis SQL.

### Cómo está organizado `ejemplos/dia20_base_datos/`

Es un módulo Go independiente (tiene su propio `go.mod`) porque depende de una librería
externa. La primera vez que lo uses necesitas conexión a internet para descargar el driver:

```powershell
cd ejemplos/dia20_base_datos
go mod tidy   # descarga modernc.org/sqlite si no está en go.mod/go.sum
go run .
```

Si no tienes internet en este momento, lee el archivo `main.go` de todos modos: el
**concepto** (conexión, tabla, CRUD, placeholders) es idéntico sin importar el driver, y el
archivo está comentado para que entiendas cada paso aunque no lo ejecutes todavía.

### Ejercicios

1. Antes de tocar SQL: dibuja (en un comentario o en papel) la tabla `productos` que vas a
   crear, con columnas `id` (entero, autoincremental, llave primaria), `nombre` (texto),
   `precio` (real) y `stock` (entero). Escribe el `CREATE TABLE ... IF NOT EXISTS` en SQL.

2. Escribe la función que inserta un producto nuevo usando `db.Exec` con placeholders
   (`INSERT INTO productos (nombre, precio, stock) VALUES (?, ?, ?)`), y que devuelva el
   id generado (pista: `result.LastInsertId()`).

3. Escribe una función que liste todos los productos (`SELECT * FROM productos`) e
   imprima cada fila, iterando con `rows.Next()`/`rows.Scan`. No olvides revisar
   `rows.Err()` después del loop y hacer `defer rows.Close()`.

4. Escribe una función que busque un producto por `id` usando `db.QueryRow` y actualice su
   `stock` con `UPDATE productos SET stock = ? WHERE id = ?`. Verifica cuántas filas se
   afectaron con `result.RowsAffected()` para saber si el id existía.

5. Escribe una función que elimine un producto por id (`DELETE FROM productos WHERE id = ?`)
   y confirma, igual que en el ejercicio 4, cuántas filas fueron afectadas.

💡 **Reto extra:** envuelve un conjunto de operaciones (ej. "insertar 3 productos") en una
**transacción** (`db.Begin()`, `tx.Exec(...)`, `tx.Commit()` o `tx.Rollback()` si algo
falla). Investiga por qué las transacciones importan en un ERP: piensa en una operación de
"transferir stock entre dos almacenes" que necesita restar en uno y sumar en otro — si el
proceso se cae a la mitad, no puedes dejar el sistema en un estado inconsistente.

---

## Día 21: Testing en Go + reto integrador de la semana

### La explicación

Go trae testing en la biblioteca estándar (paquete `testing`), sin frameworks externos
necesarios para empezar. La convención es estricta y automática: un archivo `x_test.go`
junto a `x.go`, con funciones `func TestAlgo(t *testing.T)`, ejecutadas con `go test`.

El patrón idiomático de Go para probar múltiples casos es el **table-driven test**: en vez
de escribir una función `TestX` por cada caso, defines un slice (o map) de **casos de
prueba** (entrada esperada + salida esperada) y un solo loop que los recorre todos,
llamando `t.Run(nombre, func(t *testing.T) {...})` para cada uno. Esto tiene ventajas
reales en un backend: agregar un caso nuevo es agregar una línea a la tabla, no escribir
una función nueva; y cuando un caso falla, `t.Run` te dice exactamente cuál (por nombre),
no solo "algo en TestX falló".

¿Por qué Go es tan insistente con esto? Porque en un backend que va a crecer (tu ERP con
decenas de módulos), las pruebas automatizadas son lo que te permite cambiar código con
confianza — corres `go test ./...` y sabes en segundos si rompiste algo, en vez de probar
manualmente cada endpoint cada vez.

### Cómo está organizado `ejemplos/dia21_testing/`

Es un módulo Go aparte (con su propio `go.mod`) con `operaciones.go` (funciones simples de
ejemplo) y `operaciones_test.go` (sus tests table-driven). Corre las pruebas con:

```powershell
cd ejemplos/dia21_testing
go test ./...        # corre todos los tests
go test -v ./...      # con detalle de cada caso (verbose)
```

### Ejercicios

1. Escribe una función `Dividir(a, b float64) (float64, error)` en un archivo tuyo (puedes
   crear tu propia carpeta de práctica, o reutilizar `ejemplos/dia21_testing/` como
   referencia) que devuelva un error si `b == 0`. Escribe un test *simple* (no
   table-driven todavía) que verifique el caso normal y el caso de división por cero.

2. Convierte ese test a **table-driven**: define un slice de structs con campos como
   `nombre string`, `a, b float64`, `esperado float64`, `esperaError bool`, y recórrelo con
   un `for` + `t.Run`. Agrega al menos 5 casos, incluyendo negativos y decimales.

3. Escribe una función `EsPalindromo(s string) bool` y sus tests table-driven con al menos
   6 casos (incluye string vacío, un solo caracter, con mayúsculas mezcladas, con espacios).
   Decide y documenta en un comentario: ¿tu función es sensible a mayúsculas/espacios o no?
   Tus tests deben reflejar esa decisión.

4. Corre `go test -v ./...` en tu carpeta de pruebas y confirma que ves cada caso de la
   tabla listado por su nombre en la salida. Rompe una función a propósito (ej. cambia un
   `<` por un `<=`) y observa cómo `go test` te señala exactamente qué caso de la tabla
   falló.

5. Investiga y usa `go test -cover ./...` para ver el porcentaje de cobertura de tus
   pruebas. No necesitas 100% de cobertura, pero identifica si hay alguna rama de tu código
   (ej. el caso de error) que no esté cubierta por ningún test.

### 💡 Reto integrador de la semana 3

Ya tienes todas las piezas. Construye un mini-servicio que combine:

- **HTTP** (día 18/19): un endpoint `POST /procesar` que reciba un JSON con una lista de
  números, ej. `{"numeros": [4, 15, 8, 23, 42]}`.
- **Concurrencia** (día 15): procesa cada número en su propia goroutine (por ejemplo,
  calculando si es primo, o su factorial, o cualquier operación que tome "algo" de tiempo),
  usando `sync.WaitGroup` para esperar a que todas terminen antes de responder, y un
  `sync.Mutex` (o un channel) para recolectar los resultados de forma segura.
- **JSON** (día 17): responde con un JSON como
  `{"resultados": [{"numero":4,"esPrimo":false}, {"numero":15,"esPrimo":false}, ...]}`.

Este reto no tiene solución en los ejemplos — es tuyo para armar combinando lo que ya
practicaste. Si te atoras, vuelve a `ejemplos/dia18_http_servidor.go` para la sintaxis de
HTTP y a `ejemplos/dia15_sync_select.go` para la sintaxis de WaitGroup/Mutex, pero arma la
lógica de este servicio desde cero.
