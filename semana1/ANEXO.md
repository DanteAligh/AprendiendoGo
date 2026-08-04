# Anexo — Cómo investigar por tu cuenta (semana 1)

Este anexo no enseña sintaxis de Go. Enseña algo más importante a largo
plazo: **cómo averiguar tú mismo** qué hace una función y cómo se usa,
sin depender de que alguien (o una IA) te lo explique cada vez. Esta es
la habilidad que más vas a usar el resto de tu carrera con Go.

## El problema: usaste `strconv.Atoi` sin saber qué devolvía

Si viste código de ejemplo con `strconv.Atoi(texto)` y no tenías forma de
saber que devuelve **dos valores** (`int` y `error`) hasta que alguien te
lo dijo, este anexo es para ti.

## Herramienta 1: `go doc` (la más rápida, funciona sin internet)

Desde cualquier terminal con Go en el `PATH`:

```powershell
go doc strconv.Atoi
```

Esto imprime:

```
func Atoi(s string) (int, error)
    Atoi is equivalent to ParseInt(s, 10, 0), converted to type int.
```

Con solo esa firma ya sabes tres cosas sin ejecutar nada:
1. Recibe un `string`.
2. Devuelve un `int`.
3. Devuelve además un `error` — y en Go, **si una firma termina en
   `, error)`, esa función puede fallar y tienes que revisarlo con
   `if err != nil`.**

También puedes pedir la documentación de un paquete completo:

```powershell
go doc strconv
```

Te lista todas las funciones y tipos del paquete con una línea de
descripción cada uno. Sirve para "explorar" qué hay disponible cuando no
sabes el nombre exacto de la función que necesitas.

## Herramienta 2: pkg.go.dev

[pkg.go.dev](https://pkg.go.dev) es la documentación oficial de los
paquetes de Go, en un sitio web navegable. Para cualquier paquete de la
librería estándar (por ejemplo `strconv`, `math`, `fmt`), la URL sigue el
patrón `pkg.go.dev/<paquete>`. Ahí ves lo mismo que con `go doc`, pero con
mejor formato y ejemplos de uso ("Example" en cada función).

Úsalo cuando quieras ver ejemplos completos de una función, no solo su
firma.

## Herramienta 3: el editor (si usas VSCode, GoLand, etc.)

Si instalas la extensión oficial de Go en VSCode (o usas GoLand), al
escribir el nombre de una función y pasar el cursor encima, te muestra la
firma sin salir del archivo. Es la forma más cómoda en el día a día, pero
vale la pena que sepas usar `go doc` también porque siempre está
disponible, incluso sin editor gráfico.

## La regla que resuelve tu pregunta de fondo

**Cualquier función cuya firma termina en `, error)` puede fallar, y el
patrón casi siempre es:**

```go
resultado, err := funcion(argumentos)
if err != nil {
    // decide qué hacer: mostrar un mensaje, terminar el programa, reintentar, etc.
}
```

No es algo que tengas que memorizar función por función — es una
convención que se repite en *toda* la librería estándar de Go y en
prácticamente todo el código Go que vas a leer. En cuanto veas esa forma
en una firma, ya sabes qué hacer.

En el día 5 vas a profundizar en el manejo de errores como concepto (`errors.New`, `%w`, `errors.Is`/`errors.As`). Este anexo es la base práctica de "cómo detectarlo"; ese día es
"qué hacer con él" a fondo.

---

# Segunda parte — investigar tipos, métodos e interfaces (días 4 a 6)

Arriba viste `go doc` para funciones sueltas. A partir del día 4 trabajas con
structs, métodos e interfaces, así que las preguntas cambian de forma: ya no es
solo "¿qué devuelve esta función?", sino "¿qué métodos tiene este tipo?" y
"¿qué interfaz necesito implementar?".

## Cómo ver todos los métodos de un tipo

```powershell
go doc time.Time
```

Esto no solo te muestra el struct `Time`, sino la lista completa de
métodos que tiene (`Year()`, `Month()`, `Before()`, `After()`, etc.), cada
uno con su firma. Cuando trabajes con un tipo de la librería estándar y no
sepas qué puedes hacer con él, empieza siempre por aquí.

Si quieres ver un método específico:

```powershell
go doc time.Time.Before
```

## Cómo saber si tu struct "cumple" una interfaz

En Go, un tipo implementa una interfaz **implícitamente**: si tiene todos
los métodos que la interfaz pide, ya la implementa, sin escribir `implements`
en ningún lado. Esto es distinto a Java/C# donde lo declaras explícito, y
es justo el tema del día 4.

Para verificar que tu tipo cumple una interfaz *en tiempo de compilación*
(antes de que lo descubras en producción), existe un truco idiomático que
vas a ver mucho en código Go real:

```go
var _ NombreInterfaz = (*TuStruct)(nil)
```

Esta línea no hace nada en tiempo de ejecución — es una comprobación que
el compilador hace por ti. Si `TuStruct` no implementa `NombreInterfaz`,
esta línea no compila y te dice exactamente qué método te falta.

## Errores: `go doc` también aplica

```powershell
go doc errors
go doc fmt.Errorf
```

`fmt.Errorf` con el verbo `%w` es la forma idiomática de "envolver" un
error con más contexto sin perder el error original. Vale la pena que
corras `go doc errors.Is` y `go doc errors.As` cuando llegues al día 5 —
sus descripciones cortas explican mejor que cualquier tutorial cuándo usar
cada uno.

## `go vet`: tu segundo par de ojos

Además de compilar, corre esto sobre tu código de vez en cuando:

```powershell
go vet ./...
```

`go vet` detecta errores sutiles que compilan pero probablemente están
mal (por ejemplo, un `Printf` con el verbo de formato equivocado para el
tipo que le pasaste, o una comparación de structs que no se puede hacer).
No reemplaza pensar, pero atrapa errores tontos antes de que los notes tú.

## Regla para los días 4 a 6

Cuando te encuentres con un tipo o interfaz que no conoces:
1. `go doc <paquete>.<Tipo>` para ver toda su superficie (campos y
   métodos).
2. Si necesitas implementar una interfaz, busca su definición con
   `go doc <paquete>.<Interfaz>` para ver exactamente qué métodos exige.
3. Corre `go vet` antes de asumir que tu código está bien solo porque
   compiló.
