# Día 5 · Errores y `defer` — explicación y ejercicios

El día 1 usaste `(float64, bool)`: un valor y un semáforo. Funcionaba, pero el `false` solo dice
**que** algo salió mal, no **qué**.

Hoy el `bool` se convierte en un `error` de verdad, que además **lleva el motivo dentro**.
Este es el sistema real de manejo de errores de Go, el que vas a usar en las cuatro semanas.

---

# Parte 1 · La explicación

## `error` es solo un tipo más

```go
func leerPrecio(txt string) (float64, error)
```

Mismo patrón del día 1, cambiando `bool` por `error`:

| Situación | Devuelve |
|---|---|
| Todo bien | el valor y **`nil`** |
| Algo falló | un valor de relleno y **un error** |

**`nil`** significa "nada", "vacío", "no hay". Cuando el error es `nil`, es que no hubo error.

Y quien llama comprueba:

```go
precio, err := leerPrecio("abc")
if err != nil {
	fmt.Println("falló:", err)
	return
}
```

`if err != nil` es **la línea más escrita en todo Go**. Se lee: *"si el error no es nada..."*, o
sea, *"si sí hubo error"*. La vas a escribir miles de veces.

> **Go no tiene excepciones.** En otros lenguajes un fallo "salta" por encima de tu código hasta
> que alguien lo atrapa, y no ves dónde puede pasar. En Go el error es **un valor normal** que
> vuelve por la puerta principal, y tú decides qué hacer con él ahí mismo. Es más verboso y
> mucho más difícil de ignorar.

## Crear un error

Dos formas:

```go
errors.New("el precio no puede ser negativo")          // texto fijo
fmt.Errorf("precio inválido: %s", txt)                 // con huecos, como Printf
```

`errors` y `fmt` son paquetes de la librería estándar; hay que importarlos.

## Errores centinela — errores con nombre

Si el error es solo texto, quien lo recibe no puede reaccionar distinto según el motivo: tendría
que comparar cadenas, que es frágil. Por eso los errores importantes se declaran **una vez, como
variables de paquete**:

```go
var (
	ErrNoEsNumero = errors.New("no es un número")
	ErrNegativo   = errors.New("no puede ser negativo")
)
```

Se llaman **errores centinela** (*sentinel*). La convención de nombre es empezar por `Err`.

Ahora quien te llama puede preguntar **cuál** de los dos fue.

## Envolver con `%w`

```go
return 0, fmt.Errorf("leyendo %q: %w", txt, ErrNoEsNumero)
```

El verbo **`%w`** (de *wrap*, envolver) es especial: mete el error original **dentro** del nuevo,
sin perderlo.

¿Para qué? Para añadir contexto sin destruir la información. El error resultante dice:

```
leyendo "abc": no es un número
```

Tienes las dos cosas: **dónde** pasó (tu contexto) y **qué** pasó (el error original).
Cada capa del programa puede añadir su parte, y se van encadenando.

> `%w` solo funciona en `fmt.Errorf`, no en `Printf`. Y solo debe haber uno por mensaje.

## Desenvolver con `errors.Is`

```go
if errors.Is(err, ErrNegativo) {
	fmt.Println("el usuario metió un número negativo")
}
```

`errors.Is` mira dentro de todas las capas del envoltorio buscando ese error concreto.
**No compares con `==`**: si el error viene envuelto, `err == ErrNegativo` da `false` aunque
esté dentro. `errors.Is` sí lo encuentra.

---

## `defer` — "haz esto al salir, pase lo que pase"

```go
func leer() {
	defer fmt.Println("fin de la lectura")
	fmt.Println("leyendo...")
}
```

```
leyendo...
fin de la lectura
```

`defer` **aplaza** una llamada hasta que la función termine. Y termine **como termine**:
por un `return` normal, por un `return` de error a mitad, o incluso por un `panic`.

### Para qué sirve de verdad

Para **limpiar**. Todo lo que abres hay que cerrarlo: un archivo (día 7), una conexión a base de
datos (semana 3), un candado (semana 2). Sin `defer` tendrías que acordarte de cerrarlo en cada
uno de los seis `return` de la función, y el día que añadas un séptimo se te olvidará.

```go
f, err := os.Open("datos.csv")
if err != nil {
	return err
}
defer f.Close()      // pase lo que pase de aquí en adelante, se cierra
```

**La línea de cerrar va justo debajo de la de abrir.** Así se ven juntas y no hay forma de
olvidarla.

### El orden es al revés (LIFO)

Si hay varios `defer`, se ejecutan **del último al primero**:

```go
defer fmt.Println("uno")
defer fmt.Println("dos")
defer fmt.Println("tres")
```

```
tres
dos
uno
```

Tiene lógica: si abriste A y luego B, hay que cerrar B antes que A. Se llama **LIFO** (*last in,
first out*: el último en entrar es el primero en salir), como una pila de platos.

### Una trampa: los argumentos se congelan

```go
i := 1
defer fmt.Println("valor:", i)
i = 99
// imprime "valor: 1", no 99
```

Los argumentos se evalúan **cuando escribes el `defer`**, no cuando se ejecuta. Lo que se aplaza
es la llamada, no el cálculo de lo que va dentro.

---

## Un paquete nuevo: `strconv`

Para convertir texto a número:

```go
n, err := strconv.ParseFloat(txt, 64)
```

`strconv` = *string conversion*. El `64` es la precisión (`float64`). Devuelve el número y un
error — otra vez el mismo patrón.

> **No confundas con la conversión de tipos del día 1.** `float64(unidades)` convierte un
> **número** a otro tipo de número, y no puede fallar. `strconv.ParseFloat` interpreta un
> **texto**, y sí puede fallar: `"abc"` no es ningún número.

---

# Parte 2 · Los ejercicios

---

## Ejercicio A · `defer` y su orden

**Carpeta:** `cmd\dia05a\main.go`
**Practica:** ver el comportamiento de `defer` antes de usarlo en serio.

### Qué tienes que hacer

1. Función `procesar()` que tenga un `defer` al principio imprimiendo `fin del proceso`, y
   después imprima `procesando...`.
2. Función `pila()` con tres `defer` seguidos que impriman `uno`, `dos`, `tres`, y luego un
   `fmt.Println("cuerpo")`.
3. Función `congelado()` que demuestre la trampa de los argumentos.
4. Llama a las tres desde `main`.

### Salida exacta

```
--- procesar ---
procesando...
fin del proceso
--- pila ---
cuerpo
tres
dos
uno
--- congelado ---
valor al final: 99
valor congelado en el defer: 1
```

### Pistas

- El `defer` se **escribe** el primero pero se **ejecuta** el último. Toda la gracia está ahí.
- En `pila`, fíjate en que `cuerpo` sale antes que los tres `defer`, y estos al revés de como los
  escribiste.
- Para el punto 3: declara `i := 1`, escribe el `defer` con `i` dentro, luego `i = 99`, luego
  imprime `i`. Verás que el `Println` normal dice 99 y el aplazado dice 1.

---

## Ejercicio B · `leerPrecio` con errores de verdad

**Carpeta:** `cmd\dia05b\main.go`
**Practica:** devolver `error`, envolver con `%w`, distinguir con `errors.Is`, `defer`.

### Qué tienes que hacer

1. Dos errores centinela: `ErrNoEsNumero` y `ErrNegativo`.
2. `leerPrecio(txt string) (float64, error)`:
   - Un `defer` al principio que imprima `fin de la lectura de "..."`.
   - Si el texto no es un número → `0` y un error que **envuelva** `ErrNoEsNumero`.
   - Si es negativo → `0` y un error que **envuelva** `ErrNegativo`.
   - Si todo bien → el número y `nil`.
3. En `main`, probar con `"1250.50"`, `"abc"` y `"-5"`, distinguiendo con `errors.Is` cuál de los
   dos motivos fue.

### Salida exacta

```
fin de la lectura de "1250.50"
OK: 1250.50

fin de la lectura de "abc"
ERROR: leyendo "abc": no es un número
  -> el texto no era numérico

fin de la lectura de "-5"
ERROR: leyendo "-5": no puede ser negativo
  -> el precio era negativo
```

### Pistas

- **Mira el orden de la salida**: el `fin de la lectura` sale **antes** que el `OK` o el `ERROR`.
  Es correcto: el `defer` corre al terminar `leerPrecio`, y `main` imprime su parte después.
  Si no te cuadra, sigue el hilo con el dedo.
- La estructura de los mensajes envueltos es
  `fmt.Errorf("leyendo %q: %w", txt, ErrNoEsNumero)`. El `%q` pone las comillas.
- Para distinguir en `main`:

  ```go
  if errors.Is(err, ErrNoEsNumero) { ... }
  if errors.Is(err, ErrNegativo)   { ... }
  ```

- **Prueba a comparar con `==` en vez de `errors.Is`** y verás que no entra en ningún caso: el
  error está envuelto, no es el mismo objeto. Ese experimento vale más que cualquier explicación.
- La línea en blanco entre bloques la imprimes tú desde `main`.
- El error de `strconv.ParseFloat` tíralo con `_` y devuelve el tuyo: hoy queremos mensajes
  nuestros, en español.

---

## Ejercicio C · Carrito con validación en cadena

**Carpeta:** `cmd\dia05c\main.go`
**Practica:** propagar errores entre varias funciones y ver crecer el envoltorio.

### Qué tienes que hacer

1. Reutiliza `leerPrecio` del ejercicio B.
2. Función `totalCarrito(textos []string) (float64, error)` que sume los precios. Si **alguno**
   falla, devuelve `0` y un error que diga **en qué posición** falló, envolviendo el error que
   vino de abajo.
3. En `main`, probar con dos carritos: uno bueno y uno con un texto malo.

### Salida exacta

```
Carrito 1: total 1900.40
Carrito 2: error -> artículo 2: leyendo "diez": no es un número
Carrito 2: la causa raíz fue que no era un número
```

### Pistas

- `totalCarrito` no imprime nada: calcula, y si falla devuelve el error hacia arriba. **Imprimir
  es trabajo de `main`.** Esta separación es la que hace que una función se pueda reutilizar (y
  testear, día 6).
- El envoltorio se apila: `leerPrecio` ya envolvió `ErrNoEsNumero`, y `totalCarrito` envuelve
  **ese** error otra vez añadiendo la posición. Por eso el mensaje final tiene tres partes.
- **Y aun con dos capas, `errors.Is` sigue encontrando `ErrNoEsNumero`.** Esa es la tercera línea
  de la salida y el punto entero del ejercicio: el contexto se acumula sin perder la causa.
- La posición es `i+1` porque los slices cuentan desde 0 y las personas desde 1.
- Para el carrito bueno usa textos que sumen `1900.40`, por ejemplo `"1250.50"`, `"499.90"`,
  `"150.00"`.
- Quítale el `defer` a `leerPrecio` para este ejercicio, o la salida se llenará de líneas de
  `fin de la lectura`.

---

## Ejecutar

```powershell
go run ./cmd/dia05a
go run ./cmd/dia05b
go run ./cmd/dia05c
```

## Lo que tienes que poder explicar al terminar

- Qué significa `nil` y por qué `if err != nil` es la comprobación estándar.
- Qué añade `%w` que no añade `%v`.
- Por qué `errors.Is` y no `==`.
- Qué hace `defer` y por qué la línea de cerrar va pegada a la de abrir.
