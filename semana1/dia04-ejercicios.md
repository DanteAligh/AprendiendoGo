# Día 4 · Structs, métodos e interfaces — explicación y ejercicios

Ayer aprendiste a juntar **muchas cosas del mismo tipo** (un slice de textos, un map de enteros).
Hoy aprendes a juntar **cosas distintas que describen un mismo objeto**: el nombre, el precio y
el stock de un producto, todo en una sola pieza.

Y después, a darle **comportamiento** a esa pieza.

---

# Parte 1 · La explicación

## Struct — una ficha con campos

Hasta ahora, para manejar un producto necesitabas tres variables sueltas:

```go
nombre := "Teclado"
precio := 1250.5
stock := 3
```

El problema: nada las une. Si pasas un producto a una función necesitas tres parámetros, y nada
impide mezclar el precio de uno con el stock de otro.

Un **struct** es una ficha con casillas. Tú defines el modelo una vez:

```go
type Producto struct {
	Nombre string
	Precio float64
	Stock  int
}
```

- **`type`** = "voy a crear un tipo nuevo". A partir de ahora `Producto` es un tipo, igual que
  `int` o `string`.
- **`struct`** = "y ese tipo es una ficha con estas casillas".
- Cada línea es un **campo**: nombre y tipo.

### Crear uno

```go
p := Producto{Nombre: "Teclado", Precio: 1250.5, Stock: 3}
```

Con los nombres delante. **Escríbelo siempre así**, aunque exista la forma corta
(`Producto{"Teclado", 1250.5, 3}`): si algún día cambias el orden de los campos, la forma corta
se rompe en silencio y la larga sigue funcionando.

### Usar los campos

```go
p.Nombre        // "Teclado"
p.Stock = 5     // cambiarlo
```

El punto significa "de dentro de", igual que en `fmt.Println`.

### Imprimir un struct entero

```go
fmt.Printf("%v\n", p)    // {Teclado 1250.5 3}
fmt.Printf("%+v\n", p)   // {Nombre:Teclado Precio:1250.5 Stock:3}
```

`%+v` añade los nombres de los campos. Para depurar es infinitamente mejor.

### Mayúscula inicial en los campos

Te lo he mencionado tres veces ya: en Go la mayúscula marca lo **público**. En los campos de un
struct la regla es la misma, y además hay una razón práctica: los campos en minúscula **no se
convierten a JSON** (semana 2). Por eso los verás casi siempre en mayúscula.

---

## Métodos — funciones pegadas a un tipo

Un **método** es una función que "pertenece" a un tipo. En vez de `valor(p)` escribes `p.Valor()`.

```go
func (p Producto) Valor() float64 {
	return p.Precio * float64(p.Stock)
}
```

La novedad es el paréntesis extra **antes** del nombre:

```
func  (p Producto)  Valor()  float64 {
      └─────┬────┘
      el RECEPTOR: sobre qué actúa el método
```

`p` es como el "yo" del método: dentro puedes usar `p.Precio`, `p.Stock`.

Se llama así:

```go
fmt.Println(p.Valor())
```

> ¿Por qué método y no función normal? Porque el comportamiento vive **junto** al dato que
> describe. Todo lo que sabe hacer un `Producto` está en un solo sitio.

---

## Receptor por valor vs por puntero — **el tema del día**

Esto es lo importante de hoy, y lo que más cuesta al principio.

### El problema

```go
func (p Producto) Vender(n int) {
	p.Stock = p.Stock - n
}
```

Parece correcto. **Y no funciona.** El stock no baja.

### Por qué

Cuando llamas a un método con receptor **por valor** (`p Producto`), Go **hace una fotocopia**
del producto y se la entrega al método. El método modifica la fotocopia. Al terminar, la
fotocopia se tira a la basura y el original sigue igual.

Es como si le dieras a alguien una copia de tu formulario para que corrija un dato: corrige la
copia, tú te quedas con el tuyo intacto.

### La solución: el puntero

```go
func (p *Producto) Vender(n int) {
	p.Stock = p.Stock - n
}
```

Un asterisco. `*Producto` significa **"la dirección del producto"**, no una copia de él.

Un **puntero** es una dirección. En vez de darle la fotocopia del formulario, le das **la
dirección del archivador donde está el original**. Va allí y corrige el de verdad.

### La regla, en una frase

> **¿El método necesita CAMBIAR algo? → receptor por puntero (`*T`).
> ¿Solo lee? → por valor (`T`).**

`Valor()` solo lee → `(p Producto)`.
`Vender()` modifica → `(p *Producto)`.

### Y una comodidad de Go

Podrías esperar tener que escribir algo raro para llamar a un método de puntero. No:

```go
p.Vender(2)
```

Igual que siempre. Go pone y quita los `&` y los `*` por ti. Solo tienes que acertar con el
asterisco **en la definición del método**, que es donde importa.

> Verás muchos programas donde **todos** los métodos de un tipo usan puntero, aunque algunos
> solo lean. Es una convención muy extendida: mezclar los dos estilos en el mismo tipo confunde.
> Hoy los mezclamos a propósito, para que veas la diferencia.

---

## Interfaces — describir lo que algo *sabe hacer*

Un struct dice **qué es** una cosa (sus campos). Una interfaz dice **qué sabe hacer**.

```go
type Describible interface {
	Describir() string
}
```

Se lee: *"cualquier tipo que tenga un método `Describir()` que devuelva un texto, es un
`Describible`"*.

### Lo particular de Go

En la mayoría de lenguajes hay que **declarar** que implementas una interfaz
(`class Producto implements Describible`). **En Go no se declara nada.** Si tu tipo tiene los
métodos que la interfaz pide, ya la cumple. Automáticamente. Se llama *implementación implícita*.

```go
func (p Producto) Describir() string {
	return "..."
}
// Listo: Producto ya es Describible. Sin escribir nada más.
```

### Para qué sirve

Para escribir **una función que sirva para tipos distintos**:

```go
func mostrar(d Describible) {
	fmt.Println(d.Describir())
}
```

Esa función acepta un `Producto`, un `Servicio`, o cualquier cosa futura que sepa describirse.
No le importa **qué es**, solo **qué sabe hacer**.

Y un slice puede mezclar tipos distintos si todos cumplen la interfaz:

```go
cosas := []Describible{producto, servicio}
```

Ese slice de arriba es lo que hace que las interfaces valgan la pena: es el único momento en que
puedes tener cosas de tipos diferentes en una misma lista.

---

# Parte 2 · Los ejercicios

---

## Ejercicio A · La ficha de producto

**Carpeta:** `cmd\dia04a\main.go`
**Practica:** definir un struct, crear valores, un método que solo lee.

### Qué tienes que hacer

1. Define `type Producto struct` con `Nombre string`, `Precio float64`, `Stock int`.
2. Método `Valor() float64` que devuelva `Precio × Stock`.
3. Crea dos productos e imprime cada uno con `%+v` y su valor.

### Salida exacta

```
{Nombre:Teclado mecánico Precio:1250.5 Stock:3}
Valor en almacén de Teclado mecánico: 3751.50
{Nombre:Ratón Precio:499.9 Stock:10}
Valor en almacén de Ratón: 4999.00
```

### Pistas

- `Valor()` solo lee → **receptor por valor**, `(p Producto)`.
- Vuelve el choque de tipos del día 1: `p.Precio` es `float64` y `p.Stock` es `int`.
  Hace falta `float64(p.Stock)`.
- Fíjate en que `%+v` imprime `1250.5`, no `1250.50`. Es el formato por defecto; para controlar
  los decimales necesitas `%.2f` en un `Printf` aparte, que es justo lo que hace la segunda línea.

---

## Ejercicio B · Vender stock

**Carpeta:** `cmd\dia04b\main.go`
**Practica:** receptor por puntero, y **ver con tus ojos** por qué hace falta.

Copia el `Producto` del ejercicio A y añade lo siguiente.

### Qué tienes que hacer

1. Método **por valor** `VenderMal(n int)` que reste `n` al stock.
2. Método **por puntero** `VenderBien(n int)` que haga exactamente lo mismo.
3. Con un producto de stock 10: imprime el stock, llama a `VenderMal(3)`, imprime otra vez,
   llama a `VenderBien(3)`, imprime otra vez.

### Salida exacta

```
Stock inicial: 10
Tras VenderMal(3): 10
Tras VenderBien(3): 7
```

### Pistas

- **Las dos funciones tienen el cuerpo idéntico.** Lo único que cambia es un asterisco en el
  receptor. Ese asterisco es todo el ejercicio.
- El resultado de `VenderMal` no es un fallo tuyo: es exactamente lo que tiene que pasar. Estás
  fabricando el error a propósito para no volver a caer en él nunca.
- `go vet` te avisará de algunas cosas raras, pero **no** de esta. Es un error que compila, corre
  y no se queja. Por eso hay que conocerlo.

### Bonus (hazlo después de que salga la salida de arriba)

Que `VenderBien` no permita dejar el stock en negativo: si `n` es mayor que el stock, que devuelva
`bool` diciendo que no se pudo. Es el patrón `(ok)` del día 1, ahora dentro de un método.

```
Vender 100 unidades: no hay stock suficiente
```

---

## Ejercicio C · La interfaz `Describible`

**Carpeta:** `cmd\dia04c\main.go`
**Practica:** una interfaz cumplida por dos tipos distintos, y un slice que los mezcla.

### Qué tienes que hacer

1. `type Describible interface { Describir() string }`.
2. `Producto` (nombre, precio, stock) con su `Describir()`.
3. `Servicio struct` con `Nombre string` y `HorasEstimadas int`, y su `Describir()`.
4. Función `mostrar(d Describible)` que imprima la descripción.
5. Un slice `[]Describible` con dos productos y un servicio, recorrido con `range`, llamando a
   `mostrar` en cada uno.

### Salida exacta

```
Producto: Teclado mecánico — 1250.50 MXN (3 en stock)
Producto: Ratón — 499.90 MXN (10 en stock)
Servicio: Instalación de red — 8 horas estimadas
```

### Pistas

- **En ningún sitio escribes "Producto implementa Describible".** Con que tenga el método
  `Describir() string`, ya lo es. Compruébalo: borra el método `Describir` de `Servicio` y mira
  el error. Dirá algo como:

  ```
  cannot use servicio (variable of type Servicio) as Describible value:
  Servicio does not implement Describible (missing method Describir)
  ```

  Léelo entero: te dice **exactamente** qué método falta. Vuelve a ponerlo después.
- Para construir el texto dentro de `Describir()` usa `fmt.Sprintf` (el del día 2): igual que
  `Printf` pero **devuelve** el texto en vez de imprimirlo. Un método que devuelve `string` no
  debe imprimir nada por su cuenta.
- El guion largo `—` lo copias tal cual; Go maneja UTF-8 sin problema.
- Los dos tipos no se parecen en nada —uno tiene precio y stock, el otro horas— y aun así conviven
  en el mismo slice. **Eso** es lo que te da la interfaz.

---

## Ejecutar

```powershell
go run ./cmd/dia04a
go run ./cmd/dia04b
go run ./cmd/dia04c
```

## Lo que tienes que poder explicar al terminar

- Por qué `VenderMal` no cambia nada.
- Cuándo un receptor va por puntero.
- Por qué `Producto` es `Describible` sin haberlo declarado en ninguna parte.
