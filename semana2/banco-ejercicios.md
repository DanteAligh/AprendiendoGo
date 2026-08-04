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

# Ejercicios — Semana 2: Intermedio

Este archivo es tu profesor esta semana. Cada día tiene una explicación
conceptual (léela completa antes de saltar a los ejercicios, no es relleno)
y luego una serie de ejercicios progresivos que tienes que resolver por tu
cuenta, en archivos nuevos, sin copiar los ejemplos de `ejemplos/`.

---

## Día 8: Structs

### Teoría

Hasta ahora has trabajado con tipos que ya existen: `int`, `string`, `bool`,
slices, maps. Un **struct** es cómo defines **tu propio tipo compuesto**:
una colección de campos con nombre, cada uno con su propio tipo, agrupados
bajo un solo nombre.

```go
type Usuario struct {
    Nombre string
    Edad   int
    Activo bool
}
```

¿Por qué Go usa structs en vez de clases? Esta es una decisión de diseño
central de Go: **Go no tiene clases**. No hay herencia clásica, no hay
constructores especiales, no hay jerarquías de tipos. Go separa
deliberadamente los **datos** (el struct) del **comportamiento** (los
métodos, que ves mañana en el día 9). Esto es intencional: los diseñadores
de Go vivieron los dolores de cabeza de jerarquías de herencia profundas en
otros lenguajes (C++, Java) y decidieron que Go iba a favorecer la
**composición sobre la herencia** desde el diseño del lenguaje, no como una
buena práctica opcional.

Puedes instanciar un struct de varias formas:

```go
u1 := Usuario{Nombre: "Ana", Edad: 30, Activo: true} // literal con nombres de campo (la forma recomendada)
u2 := Usuario{"Ana", 30, true}                       // literal posicional (frágil: si cambia el orden de campos, se rompe)
u3 := new(Usuario)                                    // new() reserva memoria y devuelve un puntero a un struct "vacío" (zero value)
var u4 Usuario                                        // declaración simple: u4 ya es un struct usable con sus valores zero
```

`new(Usuario)` te devuelve un `*Usuario` (un puntero) con todos los campos en
su "zero value" (`""` para strings, `0` para números, `false` para bools).
En la práctica, la mayoría del código Go usa el literal con nombres de campo
o `var` — `new()` lo vas a ver, pero se usa poco fuera de casos específicos.

**Structs anidados**: un campo de un struct puede ser otro struct.

```go
type Direccion struct {
    Calle  string
    Ciudad string
}

type Usuario struct {
    Nombre    string
    Direccion Direccion
}
```

**Structs embebidos (composición, la "herencia" a la Go)**: si en vez de
darle nombre al campo simplemente pones el tipo, Go "promueve" sus campos y
métodos al struct contenedor:

```go
type Persona struct {
    Nombre string
    Edad   int
}

type Empleado struct {
    Persona    // embebido: sin nombre de campo, solo el tipo
    Puesto string
    Salario float64
}

e := Empleado{Persona: Persona{Nombre: "Luis", Edad: 28}, Puesto: "Backend Dev", Salario: 50000}
fmt.Println(e.Nombre) // "Luis" — accedes directo, como si Empleado tuviera ese campo
```

Esto **no es herencia real**: `Empleado` no "es un" `Persona` desde el punto
de vista del sistema de tipos (no puedes pasar un `Empleado` donde se espera
un `Persona`). Es azúcar sintáctico sobre composición: Go simplemente te
deja acceder a los campos del struct embebido sin escribir `e.Persona.Nombre`.
Esta diferencia es importante y te va a evitar confusiones cuando vengas de
lenguajes orientados a objetos clásicos.

**Por qué importa en un backend real**: en un ERP, prácticamente cada tabla
de tu base de datos (Producto, Cliente, Factura, Empleado) va a ser un
struct. La composición vía embebido es cómo vas a modelar, por ejemplo, que
`Factura` y `Cotizacion` comparten un conjunto de campos comunes
(`DatosCliente`) sin duplicar código.

### Ejercicios

1. Define un struct `Producto` con campos `Nombre` (string), `Precio`
   (float64) y `Stock` (int). Crea tres instancias usando el literal con
   nombres de campo e imprímelas con `fmt.Printf` usando el verbo `%+v` (que
   muestra nombres de campo y valores).

2. Define un struct `Punto` con campos `X` y `Y` (int). Escribe una función
   `DistanciaOrigen(p Punto) float64` que reciba un `Punto` y devuelva su
   distancia al origen (0,0) usando la fórmula de distancia euclidiana
   (`math.Sqrt(x*x + y*y)`). Pruébala con al menos 3 puntos distintos.
   Entrada esperada: `Punto{X: 3, Y: 4}` → Salida esperada: `5`.

3. Define un struct anidado: `Motor` con campos `Potencia` (int, en HP) y
   `Combustible` (string), y `Auto` con campos `Marca` (string), `Modelo`
   (string) y `Motor` (de tipo `Motor`, como campo con nombre, NO embebido).
   Crea un `Auto` completo y accede al campo `Potencia` del motor usando la
   notación `auto.Motor.Potencia`.

4. Ahora rehazlo usando **embebido**: define `Auto2` que embeba `Motor`
   directamente (sin nombre de campo) además de `Marca` y `Modelo`. Accede a
   `Potencia` directamente como `auto2.Potencia` (sin pasar por `.Motor`).
   Imprime ambos (`Auto` y `Auto2`) y compara cómo cambia la forma de acceder
   al dato aunque el struct interno sea el mismo.

5. Modela un mini-inventario: define `Categoria` (string) y `Articulo` con
   campos `Nombre`, `PrecioUnitario` (float64), `Cantidad` (int) y
   `Categoria` (usa el tipo `Categoria`). Crea un slice de al menos 4
   `Articulo` distintos (de al menos 2 categorías diferentes) y escribe un
   programa que recorra el slice y calcule el valor total del inventario
   (suma de `PrecioUnitario * Cantidad` de todos los artículos), imprimiendo
   el resultado con 2 decimales.

💡 **Reto extra opcional**: agrega un campo `Proveedor` que sea a su vez un
struct embebido con `NombreProveedor` y `Telefono`. Escribe una función que
reciba un slice de `Articulo` y devuelva solo los artículos cuyo stock
(`Cantidad`) sea menor a un umbral dado (simulando una alerta de "reabastecer
inventario" — algo que en un ERP real verías en un dashboard de compras).

---

## Día 9: Métodos

### Teoría

Un **método** en Go es una función que tiene un **receptor**: un parámetro
especial declarado antes del nombre de la función que "ata" esa función a un
tipo específico.

```go
type Rectangulo struct {
    Base, Altura float64
}

func (r Rectangulo) Area() float64 {
    return r.Base * r.Altura
}
```

Aquí `r Rectangulo` es el receptor. Se lee: "el método `Area` está definido
sobre el tipo `Rectangulo`". Se llama así: `rect.Area()`.

La decisión más importante que vas a tomar con cada método es: **¿receptor
por valor o receptor por puntero?**

```go
func (r Rectangulo) DuplicarPorValor() {
    r.Base *= 2 // esto modifica una COPIA de r, el original no cambia
}

func (r *Rectangulo) DuplicarPorPuntero() {
    r.Base *= 2 // esto modifica el struct ORIGINAL a través del puntero
}
```

Cuando el receptor es por **valor** (`r Rectangulo`), Go copia todo el
struct al invocar el método. Cualquier cambio que hagas dentro del método
solo afecta a esa copia; el struct original, fuera del método, queda
intacto. Cuando el receptor es por **puntero** (`r *Rectangulo`), el método
recibe la dirección de memoria del struct original, así que los cambios que
hagas sí se reflejan afuera.

**Reglas prácticas para decidir**:

- Si el método necesita **modificar** el struct (o alguno de sus campos),
  el receptor debe ser puntero. Si usas receptor por valor para "modificar"
  algo, tu código va a compilar sin error pero el cambio va a desaparecer en
  silencio — este es uno de los bugs más comunes y confusos para quien
  empieza en Go.
- Si el struct es **grande** (muchos campos, o campos pesados como slices
  grandes), usar puntero evita copiar toda esa memoria cada vez que llamas
  al método, aunque no necesites modificar nada. Es una decisión de
  rendimiento.
- Por convención y consistencia, si **alguno** de los métodos de un tipo
  necesita puntero, es común (y recomendado) que **todos** los métodos de
  ese tipo usen puntero, para no mezclar comportamientos y confundir a quien
  lea o use tu código.
- Go te deja llamar métodos de puntero sobre una variable normal (no
  puntero) y viceversa, siempre que la variable sea "direccionable" (una
  variable normal sí lo es; un valor literal temporal no). Go inserta el `&`
  o el `*` automáticamente por ti en esos casos. Esto es conveniente pero
  vale la pena saber que pasa "por debajo".

**Por qué esto importa en un backend real**: imagina un struct
`CuentaBancaria` con un campo `Saldo`. Un método `Depositar(monto float64)`
**tiene** que usar receptor por puntero — si no, cada depósito se perdería
en cuanto el método terminara, porque estarías modificando una copia
descartable. Este es exactamente el tipo de bug que, en un ERP financiero
real, causaría que los saldos nunca se actualicen aunque tu código "no
marque ningún error".

### Ejercicios

1. Define un struct `Contador` con un campo `Valor int`. Escribe un método
   `Incrementar()` con receptor por puntero que sume 1 a `Valor`. Escribe
   también un método `Valor_Actual() int` (o `Obtener() int`) que devuelva el
   valor actual. Crea un `Contador`, llama `Incrementar()` tres veces y
   verifica que el valor final sea 3.

2. Repite el ejercicio anterior pero define un método `IncrementarMal()` con
   receptor **por valor** que también intente sumar 1 a `Valor`. Llama
   `IncrementarMal()` tres veces y observa (imprime el resultado) que
   `Valor` **no** cambió. Escribe como comentario en tu código una
   explicación de una línea de por qué no cambió.

3. Define un struct `CuentaBancaria` con campos `Titular string` y
   `Saldo float64`. Escribe métodos con receptor por puntero:
   `Depositar(monto float64)` y `Retirar(monto float64) error` — este último
   debe devolver un error (puedes usar `errors.New` o simplemente
   `fmt.Errorf`, lo vemos a fondo el día 12, por ahora solo úsalo) si se
   intenta retirar más dinero del que hay en `Saldo`. Prueba con una cuenta
   inicial de 100: deposita 50, retira 30, e intenta retirar 200 (debe
   fallar).

4. Define un struct `Rectangulo` con `Base` y `Altura` (float64). Escribe
   métodos `Area() float64` y `Perimetro() float64` (ambos con receptor por
   valor, ya que no modifican nada) y un método `Escalar(factor float64)` con
   receptor por puntero que multiplique `Base` y `Altura` por `factor`.
   Verifica que después de `Escalar(2)`, el área se haya multiplicado por 4
   (no por 2 — piensa por qué).

5. Define un struct `ListaTareas` que contenga un slice de strings
   (`Tareas []string`). Escribe métodos con receptor por puntero:
   `Agregar(tarea string)` (agrega al slice), `Completar(indice int) error`
   (elimina la tarea en esa posición, o devuelve error si el índice es
   inválido) y un método con receptor por valor `Cantidad() int` que
   devuelva cuántas tareas hay. Prueba agregando 4 tareas, completando la
   tarea en el índice 1, y verificando que `Cantidad()` devuelva 3.

💡 **Reto extra opcional**: en el ejercicio 5, después de resolverlo, cambia
`Agregar` para que use receptor por **valor** en vez de puntero (a propósito)
y observa qué pasa con el slice fuera del método. Slices son un caso
especial interesante: aunque el struct se copie, el slice interno comparte
el mismo arreglo subyacente. Investiga y escribe en un comentario qué
diferencia notas entre "agregar un elemento que cabe" vs "agregar un
elemento que fuerza al slice a crecer" (relación con `cap()` que viste en la
semana 1).

---

## Día 10: Punteros

### Teoría

Un **puntero** es una variable cuyo valor es la **dirección de memoria** de
otra variable. En Go:

- `&x` obtiene la dirección de memoria de `x` (el operador "dirección de").
- `*p` obtiene el valor almacenado en la dirección que guarda `p` (el
  operador de "desreferencia").
- El tipo `*T` significa "puntero a un valor de tipo `T`".

```go
x := 10
p := &x        // p es de tipo *int, guarda la dirección de x
fmt.Println(*p) // 10 — "lo que hay en la dirección que apunta p"
*p = 20         // modifica x indirectamente, a través del puntero
fmt.Println(x)  // 20
```

Ya viste en el día 9 por qué esto importa con métodos. Pero el concepto es
más general: en Go, **todo se pasa por valor** cuando llamas una función —
siempre, sin excepción. Cuando pasas un struct grande a una función sin
puntero, Go copia el struct completo. Si quieres que la función modifique el
original, o si quieres evitar la copia por rendimiento, pasas un puntero
explícitamente.

```go
func duplicar(n int) {
    n = n * 2 // modifica la copia local, no afecta al original
}

func duplicarPtr(n *int) {
    *n = *n * 2 // modifica el valor en la dirección original
}
```

**¿Por qué Go usa punteros explícitos en vez de "paso por referencia" como
otros lenguajes?** Esta es una decisión de diseño deliberada: en lenguajes
como Python o JavaScript, los objetos se pasan "por referencia" de forma
implícita y automática — nunca ves el mecanismo, simplemente pasa. Go
prefiere que el mecanismo sea **explícito y visible en la firma de la
función**: cuando ves `func Depositar(c *CuentaBancaria, monto float64)`,
sabes de inmediato, solo leyendo la firma, que esta función puede modificar
`c`. Esto hace que el código Go sea más fácil de razonar sin tener que
"adivinar" si una función tiene efectos secundarios sobre lo que le pasaste.
Es una decisión que prioriza la legibilidad y la previsibilidad sobre la
comodidad de escritura.

Otra diferencia importante frente a C/C++: Go **no tiene aritmética de
punteros** (no puedes hacer `p + 1` para "avanzar" a la siguiente posición
de memoria). Esto elimina una categoría entera de bugs de memoria
peligrosos. Go también tiene **garbage collector**, así que no necesitas
liberar memoria manualmente (`free()`): el runtime de Go decide cuándo un
valen apuntado por un puntero ya no se usa y lo libera automáticamente.

`new(T)` es una forma de crear un puntero a un valor zero de tipo `T`:
`p := new(int)` es equivalente a `var x int; p := &x`. Lo vas a ver poco en
la práctica comparado con `&MiStruct{...}`, pero es bueno reconocerlo cuando
aparezca en código de otros.

**Un puntero puede ser `nil`**: el zero value de cualquier tipo `*T` es
`nil` (ningún struct válido). Intentar desreferenciar (`*p`) un puntero
`nil` provoca un **panic** en tiempo de ejecución. Este es un error tan
común que ya debes empezar a acostumbrarte a preguntarte, cada vez que
recibes un puntero como parámetro: "¿podría ser `nil` aquí? ¿Debería
validarlo?".

**Por qué importa en un backend real**: vas a recibir punteros a structs
constantemente — de funciones que buscan un registro en base de datos
(`func BuscarUsuario(id int) (*Usuario, error)`, que devuelve `nil` si no lo
encuentra), de funciones que necesitan modificar el estado de un objeto
compartido, etc. Entender cuándo algo puede ser `nil` y protegerte contra
eso es una habilidad central de un backend developer en Go.

### Ejercicios

1. Declara una variable `edad int` con valor 25. Crea un puntero `p` que
   apunte a `edad`. Imprime: el valor de `edad`, la dirección que guarda
   `p` (con `%p` o simplemente imprimiendo `p`), y el valor desreferenciado
   `*p`. Luego, usando `*p = 30`, cambia el valor y confirma que `edad`
   también cambió a 30.

2. Escribe una función `intercambiar(a, b *int)` que reciba dos punteros a
   `int` e intercambie los valores a los que apuntan (sin usar una tercera
   variable auxiliar fuera de la función, usa una dentro de la función si
   la necesitas). Pruébala con dos variables `x := 5` e `y := 10`, pasando
   `&x` y `&y`, y verifica que después de la llamada `x` valga 10 e `y`
   valga 5.

3. Escribe una función `func BuscarPrimero(nums []int, objetivo int) *int`
   que reciba un slice de enteros y devuelva un puntero al primer elemento
   que sea igual a `objetivo`, o `nil` si no lo encuentra. En quien llama a
   la función, verifica **siempre** si el resultado es `nil` antes de
   desreferenciarlo, e imprime un mensaje distinto según el caso ("encontrado
   en la posición de memoria X" vs "no encontrado"). Prueba con un slice que
   sí contiene el objetivo y otro que no.

4. Define un struct `Config` con campos `Host string` y `Puerto int`.
   Escribe una función `func ActualizarPuerto(c *Config, nuevoPuerto int)`
   que modifique `Puerto` a través del puntero. Luego escribe una función
   `func ActualizarPuertoMal(c Config, nuevoPuerto int)` (sin puntero) que
   intente lo mismo. Llama ambas sobre la misma variable `Config` original
   y demuestra con `fmt.Println` la diferencia de comportamiento.

5. Escribe una función `func ContarNil(punteros []*int) int` que reciba un
   slice de punteros a `int` (algunos pueden ser `nil`, simulando "campos
   opcionales no proporcionados", algo muy común al recibir JSON en un
   backend) y devuelva cuántos de ellos son `nil`. Constrúyete un slice de
   prueba mezclando punteros válidos (usando `&variable`) y `nil` explícitos.

💡 **Reto extra opcional**: investiga por qué en Go es seguro devolver un
puntero a una variable local de una función (`func Crear() *Config { c :=
Config{...}; return &c }`), a diferencia de C, donde devolver la dirección
de una variable local es un bug clásico (la memoria del stack se libera al
salir de la función). Pista: busca el término "escape analysis" en el
contexto de Go. Escribe en un comentario, con tus propias palabras, qué
entendiste.

---

## Día 11: Interfaces

### Teoría

Una **interfaz** en Go es un conjunto de métodos que un tipo debe tener para
"cumplir" con esa interfaz. Se ve así:

```go
type Figura interface {
    Area() float64
    Perimetro() float64
}
```

Cualquier tipo que tenga métodos `Area() float64` y `Perimetro() float64`
**automáticamente** cumple con la interfaz `Figura` — no hay que declarar
"esto implementa aquello" en ningún lado. Esto se llama **implementación
implícita** o "duck typing" ("si camina como pato y grazna como pato, es un
pato"): a Go no le importa de qué tipo concreto es un valor, solo le importa
si tiene los métodos que la interfaz exige.

```go
type Circulo struct{ Radio float64 }
func (c Circulo) Area() float64      { return math.Pi * c.Radio * c.Radio }
func (c Circulo) Perimetro() float64 { return 2 * math.Pi * c.Radio }

// Circulo cumple con Figura sin escribir "implements" en ningún lado.
var f Figura = Circulo{Radio: 5}
```

**¿Por qué Go hizo esto así?** En lenguajes como Java o C#, cuando defines
una clase tienes que declarar explícitamente qué interfaces implementa
(`class Circulo implements Figura`). Eso significa que quien **escribe** el
tipo tiene que saber, de antemano, todas las interfaces que le van a
interesar a otros. Go invierte esto por completo: **las interfaces se
definen del lado de quien consume**, no del lado de quien implementa. Esto
quiere decir que tú, como consumidor de un tipo (incluso un tipo de una
librería externa que no puedes modificar), puedes definir tu propia
interfaz pequeña con exactamente los métodos que necesitas, y ese tipo la va
a cumplir automáticamente si ya tiene esos métodos. Esto habilita un patrón
muy idiomático en Go: **interfaces pequeñas, definidas justo donde se usan**
(la comunidad Go tiene un dicho: "acepta interfaces, devuelve tipos
concretos").

`any` (alias de `interface{}` desde Go 1.18) es la interfaz vacía: no exige
ningún método, así que **cualquier** valor la cumple. Es útil cuando
necesitas aceptar literalmente cualquier tipo (por ejemplo, para escribir
algo similar a un `fmt.Println` genérico), pero úsala con moderación: cuando
todo es `any`, pierdes toda la verificación de tipos en tiempo de
compilación, que es una de las mayores ventajas de Go frente a lenguajes
dinámicos.

Cuando tienes un valor de tipo interfaz y necesitas el tipo concreto de
adentro, usas **type assertion**:

```go
var f Figura = Circulo{Radio: 5}
c, ok := f.(Circulo) // ok es true si f efectivamente contiene un Circulo
if ok {
    fmt.Println("Es un círculo de radio", c.Radio)
}
```

La forma con `, ok` es la segura: si el tipo no coincide, `ok` es `false` y
`c` queda con el zero value de `Circulo`, sin panic. Si usas la forma sin
`ok` (`c := f.(Circulo)`) y el tipo no coincide, tu programa hace panic —
úsala solo cuando estés absolutamente seguro del tipo.

Cuando necesitas comparar contra **varios** tipos posibles, usas
**type switch**:

```go
switch v := f.(type) {
case Circulo:
    fmt.Println("círculo de radio", v.Radio)
case Rectangulo:
    fmt.Println("rectángulo de", v.Base, "x", v.Altura)
default:
    fmt.Println("figura desconocida")
}
```

**Por qué importa en un backend real**: las interfaces son la base de cómo
se escribe código desacoplado y testeable en Go. Por ejemplo, defines una
interfaz `Repositorio` con un método `Guardar(u Usuario) error` — tu lógica
de negocio depende de esa interfaz, no de una base de datos concreta. En
producción le pasas una implementación que habla con PostgreSQL; en tus
tests le pasas una implementación falsa en memoria. Ninguna de las dos
"declara" que implementa `Repositorio` — simplemente tienen los métodos
correctos, y Go las acepta. Este patrón es exactamente el que vas a usar
para conectar tu futuro ERP con distintos proveedores de IA (Claude, GPT,
etc.) detrás de una misma interfaz.

### Ejercicios

1. Define una interfaz `Sonido` con un método `Emitir() string`. Define dos
   structs, `Perro` y `Gato`, cada uno con su propio método `Emitir() string`
   que devuelva algo distinto (ej. "Guau" y "Miau"). Escribe una función
   `func Presentar(s Sonido)` que reciba cualquier valor que cumpla `Sonido`
   e imprima el resultado de `Emitir()`. Llámala pasando un `Perro` y un
   `Gato`.

2. Crea un slice `[]Sonido` que contenga tanto `Perro` como `Gato`
   mezclados, y recórrelo llamando `Emitir()` en cada uno. Esto demuestra
   que un slice de interfaz puede contener distintos tipos concretos
   siempre que todos cumplan la interfaz.

3. Define una interfaz `Empleado` con un método `CalcularSalario() float64`.
   Define dos structs: `EmpleadoFijo` (con `SalarioBase float64`) y
   `EmpleadoPorHoras` (con `HorasTrabajadas int` y `TarifaHora float64`).
   Implementa `CalcularSalario()` para cada uno con la lógica correspondiente
   (el fijo simplemente devuelve su salario base; el de por horas multiplica
   horas por tarifa). Escribe una función que reciba un slice de `Empleado`
   y devuelva la suma total de la nómina.

4. Usando los tipos del ejercicio 3 (o similares), escribe una función
   `func Describir(e Empleado)` que use **type switch** para imprimir un
   mensaje distinto según si el `Empleado` es un `EmpleadoFijo` o un
   `EmpleadoPorHoras` (por ejemplo, mencionando el detalle específico de
   cada tipo: el salario base en un caso, las horas y tarifa en el otro).

5. Escribe una función `func Procesar(dato any) string` que reciba un valor
   de tipo `any` y, usando type switch, devuelva un string distinto según si
   el valor recibido es un `int`, un `string`, un `bool`, o "tipo no
   soportado" para cualquier otro caso. Pruébala pasando al menos 4 valores
   de tipos distintos, incluyendo uno que caiga en el `default`.

💡 **Reto extra opcional**: define una interfaz `Validador` con un método
`Validar() error`. Haz que dos structs de ejercicios anteriores (por ejemplo
`EmpleadoFijo` y `Producto` del día 8) implementen `Validar()` con alguna
regla de negocio simple (ej. que el salario o el precio no sean negativos).
Escribe una función genérica `func ValidarTodos(items []Validador) []error`
que devuelva todos los errores encontrados. Esto es el patrón real detrás de
la validación de datos de entrada en cualquier API backend.

---

## Día 12: Manejo de errores idiomático

### Teoría

Go **no tiene excepciones** (`try`/`catch`/`throw`). Esta es, otra vez, una
decisión de diseño deliberada y muy discutida en la comunidad de
programación: los diseñadores de Go decidieron que los errores son
**valores normales**, del mismo modo que un `int` o un `string`, y que deben
manejarse con el mismo control de flujo explícito que cualquier otro dato —
no con un mecanismo especial de "salto" oculto como una excepción.

El tipo `error` es una interfaz muy simple:

```go
type error interface {
    Error() string
}
```

Cualquier tipo con un método `Error() string` cumple la interfaz `error`
(sí, es el mismo concepto de duck typing del día 11 aplicado a errores).

El patrón que vas a ver en **absolutamente todo** el código Go es:

```go
resultado, err := hacerAlgo()
if err != nil {
    // manejar el error: devolverlo, loggearlo, envolver con más contexto, etc.
    return err
}
// aquí puedes usar "resultado" con la certeza de que no hubo error
```

Este patrón se repite tanto que al principio se siente repetitivo si vienes
de `try/catch`. Pero la ventaja es que **el flujo de errores es visible y
explícito en cada línea** — no hay "saltos invisibles" a un bloque `catch`
que puede estar a cientos de líneas de distancia, ni excepciones que se
"escapan" silenciosamente de una función sin que su firma lo indique. Si una
función puede fallar, su firma lo dice (`func Dividir(a, b float64) (float64,
error)`), y quien la llama está obligado (por convención, no por el
compilador) a decidir qué hacer con ese `error`.

**Crear errores**:

```go
import "errors"

err := errors.New("el saldo es insuficiente") // error simple, un mensaje fijo

err2 := fmt.Errorf("no se pudo procesar el pedido %d", id) // error con datos formateados
```

**Wrapping (envolver) errores con `%w`**: cuando una función de bajo nivel
falla y quieres agregar contexto sin perder el error original, usas
`fmt.Errorf` con el verbo `%w`:

```go
func procesarPedido(id int) error {
    err := guardarEnBaseDeDatos(id)
    if err != nil {
        return fmt.Errorf("procesarPedido: fallo al guardar pedido %d: %w", id, err)
    }
    return nil
}
```

`%w` (a diferencia de `%v` o `%s`) preserva el error original **envuelto**
dentro del nuevo, formando una especie de "cadena de contexto". Esto te deja
inspeccionar después si, en el fondo de esa cadena, hay un error específico
que te interesa, usando:

- **`errors.Is(err, errOriginal)`**: responde "¿en algún punto de esta
  cadena de errores envueltos, está exactamente este error centinela
  (`errOriginal`)?". Se usa para errores que son valores fijos y
  reconocibles, como `sql.ErrNoRows` (muy común: "¿este error significa que
  el registro no se encontró?").
- **`errors.As(err, &miTipoDeError)`**: responde "¿en algún punto de esta
  cadena, hay un error de este tipo concreto?", y si lo hay, te lo asigna
  para que puedas acceder a sus campos específicos (útil cuando defines tus
  propios tipos de error con datos extra, como un código HTTP).

**Por qué esto importa en un backend real**: en un ERP, vas a tener capas
(handler HTTP → servicio → repositorio → base de datos). Un error de "fila
no encontrada" en la base de datos debe poder viajar hacia arriba con
contexto agregado en cada capa ("no se pudo cargar el producto 42: fila no
encontrada"), pero sin perder la posibilidad de que la capa HTTP, arriba de
todo, pregunte con `errors.Is` si el error de fondo fue "no encontrado" para
devolver un 404 en vez de un 500. Ese es exactamente el propósito de
`%w` + `errors.Is`/`errors.As`.

### Ejercicios

1. Escribe una función `func Dividir(a, b float64) (float64, error)` que
   devuelva un error (con `errors.New`) si `b` es 0, y el resultado de la
   división en caso contrario. En `main`, llama la función con un caso
   válido y un caso de división entre cero, manejando el `if err != nil` en
   ambos casos e imprimiendo el mensaje de error cuando corresponda.

2. Escribe una función `func ValidarEdad(edad int) error` que devuelva un
   error con `fmt.Errorf` (con el valor recibido incluido en el mensaje) si
   la edad es negativa o mayor a 120, y `nil` si es válida. Prueba con al
   menos 4 valores distintos, incluyendo casos válidos e inválidos.

3. Declara un error centinela a nivel de paquete:
   `var ErrStockInsuficiente = errors.New("stock insuficiente")`. Escribe una
   función `func Vender(stock, cantidad int) (int, error)` que devuelva
   `ErrStockInsuficiente` (exactamente esa variable, no un error nuevo) si
   `cantidad > stock`, o el stock restante si la venta es válida. En
   `main`, después de llamar `Vender` con un caso que falla, usa
   `errors.Is(err, ErrStockInsuficiente)` para imprimir un mensaje
   específico distinto al de un error genérico.

4. Escribe una función `func CargarPedido(id int) error` que simule un fallo
   interno devolviendo `errors.New("pedido no existe en la base de datos")`,
   y una función `func ProcesarPedido(id int) error` que llame a
   `CargarPedido` y, si falla, devuelva el error **envuelto** con
   `fmt.Errorf` y `%w`, agregando el `id` del pedido al mensaje. En `main`,
   llama `ProcesarPedido` con un id que falla e imprime el error final
   (deberías ver ambos mensajes concatenados en la cadena de contexto).

5. Define un tipo de error propio:

   ```go
   type ErrorValidacion struct {
       Campo string
       Motivo string
   }
   func (e *ErrorValidacion) Error() string {
       return fmt.Sprintf("campo '%s' inválido: %s", e.Campo, e.Motivo)
   }
   ```

   Escribe una función `func ValidarProducto(nombre string, precio float64)
   error` que devuelva un `*ErrorValidacion` si `nombre` está vacío o si
   `precio` es negativo. En `main`, captura el error y usa
   `errors.As(err, &miErrorValidacion)` (con `miErrorValidacion` de tipo
   `*ErrorValidacion`) para, si el `as` tiene éxito, imprimir por separado el
   campo `Campo` y el campo `Motivo` del error.

💡 **Reto extra opcional**: investiga la diferencia entre devolver un error
"plano" desde el inicio de una cadena de llamadas vs envolverlo en cada
capa con `%w` y contexto propio. Escribe una función de 3 niveles
(`nivelUno` llama a `nivelDos` llama a `nivelTres`, y `nivelTres` es la que
falla) donde cada nivel agregue su propio contexto con `%w`. Imprime el
error final en `main` y observa la cadena completa de mensajes.

---

## Día 13: Paquetes y módulos

### Teoría

Hasta ahora todos tus programas han sido un solo archivo `.go` con
`package main`. Cualquier proyecto real (incluido tu futuro ERP) va a estar
organizado en **múltiples archivos y múltiples paquetes**, porque poner todo
en un archivo gigante se vuelve imposible de mantener.

**Módulo vs paquete**: un **módulo** es la unidad de "proyecto" en Go —
básicamente, tu repositorio, definido por un archivo `go.mod` en la raíz,
que declara el nombre del módulo (el "module path") y la versión de Go que
usas. Un **paquete** es una unidad más pequeña: una carpeta con uno o más
archivos `.go` que comparten el mismo `package nombre` en su primera línea.
Un módulo puede (y normalmente sí) contener muchos paquetes.

`go mod init nombreDelModulo` crea el `go.mod`. Cuando trabajas sin conexión
a un repositorio remoto real (como en este ejercicio), el nombre puede ser
simple, ej. `go mod init ejemplodia13`. En un proyecto real que vas a subir a
GitHub, normalmente usarías algo como `go mod init
github.com/tuusuario/tuproyecto`, porque ese path es lo que otros
proyectos usarían para importar el tuyo — pero para un proyecto local que
no vas a publicar, un nombre simple funciona perfectamente.

**Cómo se organiza un proyecto multi-paquete**:

```
miproyecto/
  go.mod                <- module miproyecto
  main.go               <- package main, puede importar los paquetes de abajo
  figuras/
    figuras.go           <- package figuras
  clientes/
    clientes.go          <- package clientes
```

Desde `main.go`, importas el paquete local con el module path + la ruta de
la carpeta:

```go
import "miproyecto/figuras"
```

**Exportado vs no exportado**: esta es una de las convenciones más
importantes de Go, y no es solo estilo — es una **regla del lenguaje**.
Cualquier identificador (función, tipo, campo de struct, variable, constante)
cuyo nombre **empiece con mayúscula** es **exportado**: visible y usable
desde otros paquetes que importen el tuyo. Cualquier identificador que
empiece con **minúscula** es **privado al paquete**: solo se puede usar
dentro del mismo paquete, ni siquiera otros paquetes del mismo módulo lo
pueden ver.

```go
package figuras

func AreaCirculo(radio float64) float64 { ... } // exportada: main.go SÍ puede llamarla
func validarRadio(radio float64) bool { ... }   // privada: solo figuras.go puede usarla
```

Esto reemplaza el concepto de `public`/`private`/`protected` que existe en
otros lenguajes con palabras clave — en Go, la visibilidad **es** la
convención de mayúscula/minúscula del nombre. Es una decisión de diseño que
prioriza la simplicidad: no hay una palabra clave más que aprender, la
visibilidad se lee directamente del nombre.

**Por qué importa en un backend real**: cuando construyas tu ERP, vas a
tener paquetes como `productos`, `clientes`, `facturacion`, `ia`
(integración con modelos de lenguaje), etc. Cada paquete va a exponer
únicamente las funciones y tipos que otros paquetes necesitan usar
(exportados), y va a esconder los detalles internos de implementación (no
exportados) — esto es la base de un buen diseño modular: cada paquete tiene
una "API pública" pequeña y controlada, y el resto es un detalle interno que
puedes cambiar libremente sin romper a nadie más.

### El mini-proyecto de este día

En vez de ejercicios sueltos como los otros días, hoy vas a construir (por
tu cuenta, en una carpeta nueva, no en `ejemplos/dia13_paquetes` que es solo
referencia) tu propio mini-proyecto multi-paquete:

1. Crea una carpeta nueva para tu ejercicio (por ejemplo
   `mi-ejercicio-dia13/`, fuera de `ejemplos/`) y corre
   `go mod init ejerciciodia13` dentro de ella.

2. Dentro de esa carpeta, crea un subdirectorio para tu propio paquete —
   elige un dominio distinto al del ejemplo de referencia (que usa
   `figuras`). Puedes hacer, por ejemplo, un paquete `texto` con funciones
   para invertir un string, contar vocales, y verificar si es un palíndromo;
   o un paquete `conversion` con funciones para convertir entre unidades
   (celsius a fahrenheit, kilómetros a millas, etc.). Debe tener al menos 3
   funciones **exportadas**.

3. Agrega también al menos una función o variable **no exportada** dentro de
   ese paquete, que sea usada internamente por una de las funciones
   exportadas (por ejemplo, una función auxiliar de validación). Esto
   demuestra que entiendes la diferencia entre lo que expones y lo que
   escondes.

4. Escribe un `main.go` en la raíz de tu carpeta de ejercicio que importe tu
   paquete (usando el module path correcto) y llame a las 3+ funciones
   exportadas con al menos 2 casos de prueba cada una, imprimiendo los
   resultados.

5. Verifica que compile y corra con `go run main.go` desde la raíz de tu
   carpeta de ejercicio.

**Pista conceptual**: si intentas llamar desde `main.go` a una función no
exportada de tu paquete (nombre en minúscula), el compilador te va a dar un
error. Pruébalo a propósito una vez para ver el mensaje de error exacto que
da Go — te vas a topar con ese mismo error muchas veces en proyectos reales
cuando accidentalmente uses algo que no debería ser visible desde afuera.

💡 **Reto extra opcional**: agrega un segundo paquete a tu mini-proyecto
(por ejemplo, si hiciste `texto`, agrega también `numeros` con funciones
sobre números primos o Fibonacci) e impórtalo también desde `main.go`,
demostrando que un mismo módulo puede tener varios paquetes independientes
que conviven sin conflicto.

---

## Día 14: Goroutines, channels y reto integrador

### Teoría

Todo lo que has escrito hasta ahora corre de forma **secuencial**: una
línea después de la otra, una función espera a que la anterior termine. Una
**goroutine** es una función que Go ejecuta de forma **concurrente**
(potencialmente al mismo tiempo que otras), con una sintaxis increíblemente
simple: solo agregas la palabra clave `go` antes de una llamada a función.

```go
go hacerAlgo() // esto arranca "hacerAlgo" en una goroutine y continúa de inmediato,
               // sin esperar a que termine
```

Las goroutines son la razón por la que se dice que "la concurrencia está en
el corazón de Go": son extremadamente livianas (puedes tener decenas de
miles corriendo al mismo tiempo sin problema, a diferencia de los threads
del sistema operativo tradicionales, que son mucho más pesados), y el
lenguaje las trata como ciudadanos de primera clase, no como una librería
externa.

El problema inmediato que surge es: si arrancas una goroutine y tu función
`main` sigue corriendo sin esperarla, ¿cómo te comunicas con ella? ¿Cómo
sabes cuándo terminó? Ahí entran los **channels**.

Un **channel** es un tubo tipado por el que puedes enviar y recibir valores
entre goroutines, de forma segura (sin las condiciones de carrera que
tendrías si varias goroutines tocaran la misma variable directamente).

```go
canal := make(chan int) // channel sin buffer, de enteros

go func() {
    canal <- 42 // enviar un valor al channel
}()

valor := <-canal // recibir un valor del channel (esto BLOQUEA hasta que alguien envíe algo)
fmt.Println(valor) // 42
```

El operador `<-` es el mismo símbolo tanto para enviar (`canal <- valor`)
como para recibir (`valor := <-canal`) — la dirección de la flecha respecto
al nombre del canal te dice cuál es cuál.

**Channel sin buffer vs con buffer**: `make(chan int)` crea un channel sin
buffer: **cada envío bloquea hasta que alguien esté listo para recibir**, y
viceversa — es una sincronización directa entre dos goroutines, como pasarse
un objeto de mano en mano; ninguna de las dos partes puede "adelantarse".
`make(chan int, 3)` crea un channel con buffer de tamaño 3: puedes enviar
hasta 3 valores sin que nadie los reciba todavía (se quedan "en cola"), y
solo cuando el buffer se llena, el siguiente envío bloquea. El buffer no
elimina la necesidad de coordinación, solo te da un margen de holgura.

**Por qué esto es la base de la concurrencia en Go**: Go fue diseñado desde
el día uno pensando en programas concurrentes (viene de un contexto de
Google, donde los sistemas backend manejan miles de operaciones
simultáneas: requests HTTP, conexiones a bases de datos, llamadas a APIs
externas). La filosofía de Go, resumida en una frase famosa de su
documentación, es: *"no te comuniques compartiendo memoria; en cambio,
comparte memoria comunicándote"*. En vez de que varias goroutines lean y
escriban la misma variable protegida con locks (el modelo tradicional de
hilos con mutex), Go te anima a que las goroutines se pasen datos por
channels, de forma explícita y ordenada. Esto hace que la concurrencia sea
más fácil de razonar correctamente.

**Por qué importa en un backend real**: en un ERP con integración de IA,
vas a necesitar, por ejemplo, procesar varios documentos en paralelo, o
hacer varias llamadas a una API de IA al mismo tiempo sin bloquear el resto
del servidor mientras esperas la respuesta. Goroutines y channels son
exactamente la herramienta para eso. Esta semana solo ves lo básico — la
semana 4, cuando construyas backends reales, vas a profundizar con
`sync.WaitGroup`, `select`, timeouts con `context`, y patrones de
concurrencia más avanzados.

### Ejercicios

1. Escribe un programa con una función `func saludar(nombre string)` que
   simplemente imprima un saludo. Llámala normal (sin `go`) tres veces con
   nombres distintos, y luego llámala tres veces más pero anteponiendo `go`.
   Agrega un `time.Sleep(100 * time.Millisecond)` al final de `main` (vas a
   necesitar importar `"time"`) para darle tiempo a las goroutines de
   correr antes de que el programa termine. Observa (y comenta en el
   código) la diferencia de orden y comportamiento entre ambas formas de
   llamar.

2. Crea un channel sin buffer de tipo `string`. Arranca una goroutine con
   `go func() { ... }()` que envíe el mensaje `"listo"` al channel después
   de simular trabajo con `time.Sleep(1 * time.Second)`. En `main`, recibe
   del channel con `<-canal` e imprime un mensaje antes y después de recibir,
   para que quede claro que `main` efectivamente esperó (se bloqueó) hasta
   que la goroutine envió el valor.

3. Crea un channel de enteros con buffer de tamaño 5. Sin usar ninguna
   goroutine todavía, envía 5 valores al channel y luego recíbelos los 5 uno
   por uno, imprimiéndolos. Esto demuestra que un channel con buffer no
   necesita que haya alguien recibiendo "al mismo tiempo" que se envía,
   hasta que el buffer se llena.

4. Escribe una función `func calcularCuadrados(nums []int, resultados
   chan<- int)` (nota el tipo `chan<- int`: un channel de "solo envío" desde
   la perspectiva de esta función) que recorra `nums` y envíe el cuadrado de
   cada número al channel `resultados`. En `main`, crea un channel con
   buffer del tamaño de tu slice de números, lánzala con `go`, y recibe (en
   un ciclo `for`) tantos valores como números enviaste, imprimiéndolos.

5. **Reto integrador de la semana**: construye un mini "procesador de
   pedidos" que combine structs, métodos, interfaces y manejo de errores:
   - Define un struct `Pedido` con `ID int`, `Cliente string` y
     `Monto float64`.
   - Define una interfaz `Procesador` con un método
     `Procesar(p Pedido) error`.
   - Crea al menos dos implementaciones distintas de `Procesador` (por
     ejemplo, `ProcesadorEfectivo` y `ProcesadorTarjeta`), cada una con su
     propia lógica de validación que pueda fallar (ej. devolver error si
     `Monto <= 0`, o si el `Cliente` está vacío).
   - Escribe una función `func ProcesarLote(pedidos []Pedido, p Procesador)
     []error` que procese todos los pedidos con el `Procesador` recibido y
     devuelva un slice con todos los errores encontrados (puede tener
     elementos `nil` en las posiciones sin error, o puedes optar por
     devolver solo los errores reales — decide tú y sé consistente).
   - Como parte final, lanza el procesamiento de un lote de pedidos dentro
     de una goroutine, y usa un channel para recibir de vuelta el slice de
     errores en `main` una vez que termine, en vez de llamar la función
     directamente. Esto conecta todo lo de la semana en un solo programa.

💡 **Reto extra opcional**: en el ejercicio 4, cambia el channel sin buffer
por uno con buffer del tamaño exacto del slice, y luego prueba a
propósito con un buffer más pequeño que el slice sin cambiar la lógica de
recepción, para observar el comportamiento (deadlock si nadie recibe a
tiempo). Escribe en un comentario qué pasó y por qué — entender un deadlock
simple ahora te va a ahorrar mucha frustración cuando la concurrencia se
vuelva más compleja en la semana 4.
