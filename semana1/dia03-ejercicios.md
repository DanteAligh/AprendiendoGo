# Día 3 · Slices y maps — explicación y ejercicios

Hasta hoy has manejado valores sueltos: un peso, una altura, un contador.
Hoy aprendes a guardar **muchos valores juntos**. Son las dos estructuras que más vas a usar
en toda tu vida como programador.

- **Slice** = una lista ordenada. Como la lista de la compra: hay un primero, un segundo, un tercero.
- **Map** = un diccionario. Buscas por una palabra y te da su significado. No hay orden.

---

# Parte 1 · La explicación

## Slices

### Crear

```go
frutas := []string{"manzana", "pera", "uva"}
```

Los corchetes vacíos `[]` delante del tipo significan "lista de". `[]string` es "lista de textos",
`[]float64` "lista de decimales", `[]int` "lista de enteros". **Todos los elementos son del mismo
tipo**: no puedes mezclar textos y números en el mismo slice.

Las llaves `{}` con los valores dentro son los elementos iniciales.

### Acceder por posición

```go
frutas[0]   // "manzana"
frutas[1]   // "pera"
```

> **Se cuenta desde CERO.** El primero es el 0, el segundo el 1. Es la convención de casi todos
> los lenguajes y al principio despista. Piensa en la posición como "cuántos me salto para llegar":
> al primero no me salto ninguno, así que es el 0.

### Cuántos hay: `len`

```go
len(frutas)   // 3
```

`len` (de *length*, longitud) es una función interna de Go, no hay que importarla.

**Consecuencia importante:** si hay 3 elementos, las posiciones válidas son 0, 1 y 2.
**La última posición siempre es `len(x)-1`**:

```go
frutas[len(frutas)-1]   // el último, sea cual sea el tamaño
```

Si pides `frutas[3]` cuando solo hay 3 elementos, el programa **muere en el acto**:

```
panic: runtime error: index out of range [3] with length 3
```

Un **panic** es la muerte súbita de un programa Go: se para y escupe el motivo. `index out of
range` = "te has salido de la lista".

### Añadir: `append`

```go
frutas = append(frutas, "kiwi")
```

Lee bien esa línea, porque es **la trampa nº 1 de Go para principiantes**:

> **`append` no modifica el slice: devuelve uno nuevo.** Hay que reasignarlo.

Si escribes solo `append(frutas, "kiwi")` sin el `frutas =` delante, el kiwi se pierde y el
compilador ni te avisa en algunos casos. La razón profunda tiene que ver con cómo Go gestiona la
memoria (cuando la lista se llena, se crea otra más grande en otro sitio), pero de momento
**apréndetelo como una regla: `append` siempre se reasigna.**

### Recorrer: `for ... range`

```go
for i, f := range frutas {
	fmt.Printf("%d: %s\n", i, f)
}
```

`range` entrega **dos** cosas en cada vuelta: la **posición** y el **valor**.
Si no te interesa la posición, la tiras con `_`:

```go
for _, f := range frutas {
```

Y si solo quieres la posición, escribes una sola variable: `for i := range frutas`.

### Cortar: `frutas[1:3]`

```go
frutas[1:3]   // desde la posición 1 hasta ANTES de la 3 → elementos 1 y 2
```

El primer número **se incluye**, el segundo **no**. Suena arbitrario, pero tiene una ventaja:
`fin - inicio` te da directamente cuántos elementos hay (3-1 = 2).

---

## Maps

### Crear

```go
edades := map[string]int{"Ana": 30, "Luis": 25}
```

```
map[TIPO_DE_LA_CLAVE]TIPO_DEL_VALOR
```

Aquí: la **clave** es texto (el nombre) y el **valor** es un entero (la edad).
Buscas por clave, te devuelve el valor. Como un diccionario de papel.

Para uno vacío:

```go
contador := map[string]int{}
```

### Leer y escribir

```go
edades["Ana"]        // 30
edades["Sara"] = 28  // añade o reemplaza
```

**Aquí hay algo que sorprende:**

```go
edades["Nadie"]      // 0  ← ¡no falla!
```

Pedir una clave que no existe **no da error**: devuelve el **valor cero** del tipo.
Para `int` es `0`, para `string` es `""`, para `bool` es `false`.

Eso permite un truco precioso para contar:

```go
contador[palabra]++
```

La primera vez la palabra no existe → vale 0 → suma 1 → queda 1. No hace falta comprobar nada.

### Distinguir "no está" de "está y vale cero"

Como leer una clave inexistente devuelve 0, ¿cómo sabes si el 0 es real? Con el patrón
`(valor, ok)` que ya conoces del día 1:

```go
n, ok := contador["rust"]
if !ok {
	fmt.Println("esa palabra no aparece")
}
```

**Es exactamente el mismo mecanismo que en `dividir` y `calcularIMC`**: un valor y un semáforo.
Verás este patrón por todo Go.

### Borrar y contar

```go
delete(edades, "Luis")
len(edades)
```

### ⚠️ El orden de un map es ALEATORIO

```go
for clave, valor := range edades {
```

Esto funciona, pero **el orden cambia en cada ejecución**. No es un fallo: Go lo hace **a
propósito**, aleatorizándolo, para que nadie escriba código que dependa de un orden que no está
garantizado.

Consecuencia práctica: si necesitas imprimir en un orden concreto, **guarda las claves en un
slice** (que sí tiene orden) y recorre el slice.

Por eso los ejercicios de hoy están diseñados para no depender del orden del map. Si tu salida
sale bailando entre ejecuciones, ya sabes por qué.

---

# Parte 2 · Los ejercicios

---

## Ejercicio A · Manejo de una lista

**Carpeta:** `cmd\dia03a\main.go`
**Practica:** crear, `len`, `append`, recorrer, acceder al último, cortar.

### Qué tienes que hacer

Partiendo de `frutas := []string{"manzana", "pera", "uva"}`:

1. Imprime cuántas hay.
2. Añade `"kiwi"` e imprime cuántas hay ahora.
3. Recórrelas imprimiendo posición y nombre.
4. Imprime la primera y la última (**la última sin escribir el número 3 a mano**).
5. Imprime el trozo de la posición 1 a la 2.

### Salida exacta

```
Longitud inicial: 3
Longitud tras append: 4
0: manzana
1: pera
2: uva
3: kiwi
Primera: manzana
Ultima: kiwi
Del 1 al 2: [pera uva]
```

### Pistas

- El punto 4 es el importante: `frutas[3]` funciona hoy, pero se rompe en cuanto la lista cambie
  de tamaño. Usa `len(frutas)-1`. **El código no debe conocer números mágicos que puedan cambiar.**
- Para la última línea, `fmt.Println` sabe imprimir un slice entero: sale entre corchetes y con
  espacios, `[pera uva]`. No tienes que recorrerlo.
- Prueba a propósito `fmt.Println(frutas[10])` para ver tu primer `panic`. Luego bórralo.

---

## Ejercicio B · Contador de palabras

**Carpeta:** `cmd\dia03b\main.go`
**Practica:** map como contador, recorrer para buscar el máximo, el patrón `(valor, ok)`.

### Qué tienes que hacer

Con esta lista:

```go
palabras := []string{"go", "java", "go", "python", "go", "java"}
```

1. Función `contar(palabras []string) map[string]int` que devuelva cuántas veces aparece cada una.
2. Imprime el recuento de `go`, `java` y `python` **en ese orden**.
3. Función `masRepetida(c map[string]int) (string, int)` que devuelva la palabra más repetida y
   su número.
4. Comprueba que `"rust"` no está, y dilo.

### Salida exacta

```
go aparece 3 veces
java aparece 2 veces
python aparece 1 veces
La mas repetida es "go" con 3 apariciones
rust no aparece en la lista
```

### Pistas

- Todo el punto 1 cabe en tres líneas gracias a `contador[p]++`. Relee la explicación si no ves
  por qué no hace falta comprobar si la clave existe.
- **Punto 2, atención:** no recorras el map con `range` — el orden saldría aleatorio. Pide las
  tres claves explícitamente, o recórrelas desde un slice `[]string{"go","java","python"}`.
- Para el máximo: una variable con el mejor hasta ahora, empezando en 0, y la vas sustituyendo
  cuando encuentres algo mayor. Es el mismo patrón del acumulador de ayer, pero comparando en
  vez de sumando.
- Punto 4: `_, ok := c["rust"]` — el `_` porque el valor no te importa, solo si existe.
- `%q` imprime un texto **entre comillas**: `"go"`. Es distinto de `%s`, que lo imprime pelado.

---

## Ejercicio C · Notas por alumno

**Carpeta:** `cmd\dia03c\main.go`
**Practica:** un map cuyo valor es un slice — combinar las dos estructuras.

### Qué tienes que hacer

```go
notas := map[string][]float64{
	"Ana":  {8.5, 9.0, 7.5},
	"Luis": {6.0, 5.5},
	"Sara": {10.0, 9.5, 9.0, 8.5},
}
```

1. Función `promedio(ns []float64) float64`.
2. Por cada alumno —**en el orden Ana, Luis, Sara**— imprime cuántas notas tiene y su promedio
   con dos decimales.

### Salida exacta

```
Ana: 3 notas, promedio 8.33
Luis: 2 notas, promedio 5.75
Sara: 4 notas, promedio 9.25
```

### Pistas

- Lee el tipo despacio: `map[string][]float64` es *"un map cuya clave es texto y cuyo valor es
  una lista de decimales"*. Se lee de izquierda a derecha.
- Dentro de las llaves, los slices de cada alumno se escriben sin repetir `[]float64` delante:
  Go ya sabe de qué tipo son. Basta `{8.5, 9.0, 7.5}`.
- **El orden fijo otra vez:** define `alumnos := []string{"Ana", "Luis", "Sara"}` y recorre ese
  slice, usando cada nombre como clave del map. Si recorres el map directamente, el orden bailará
  en cada ejecución.
- En `promedio`, ojo con los tipos: `len(ns)` es un `int` y la suma es `float64`. **No se pueden
  dividir directamente.** Te hará falta `float64(len(ns))`, como en el ejercicio A del día 1.
- ¿Qué pasa si la lista está vacía? Dividirías entre cero y saldría `NaN` (*Not a Number*).
  Hoy no lo arreglamos —mañana y el día 5 verás cómo—, pero piénsalo: ya sabes detectarlo.

---

## Ejecutar

```powershell
go run ./cmd/dia03a
go run ./cmd/dia03b
go run ./cmd/dia03c
```
