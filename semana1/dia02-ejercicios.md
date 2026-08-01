# Día 2 · Control de flujo — explicación y ejercicios

Ayer el programa iba en línea recta: de arriba abajo, una vez cada línea.
Hoy aprendes a **repetir** y a **decidir**. Con esto ya puedes escribir casi cualquier programa.

Tres herramientas: `for`, `if` y `switch`.

---

# Parte 1 · La explicación

## `for` — la única palabra de bucle en Go

Un **bucle** es "repite esto varias veces". Otros lenguajes tienen tres o cuatro palabras
distintas (`while`, `do`, `foreach`…). **Go tiene una sola: `for`.** Cambia de forma según
lo que le pongas, pero la palabra siempre es la misma. Menos que memorizar.

### Forma 1 · El `for` clásico (contar)

```go
for i := 1; i <= 3; i++ {
	fmt.Println("vuelta", i)
}
```

```
vuelta 1
vuelta 2
vuelta 3
```

Dentro del `for` hay **tres partes separadas por punto y coma**, y cada una se ejecuta en
un momento distinto:

```
for  i := 1  ;  i <= 3  ;  i++  {
     └──┬──┘     └──┬─┘    └┬┘
        │           │       └── 3. DESPUÉS de cada vuelta
        │           └────────── 2. ANTES de cada vuelta: ¿sigo?
        └────────────────────── 1. UNA vez, al empezar
```

El orden real es: **1 → 2 → cuerpo → 3 → 2 → cuerpo → 3 → 2 → …** hasta que la 2 sea falsa.

- **`i := 1`** crea la variable contador. Vive **solo dentro del bucle**; fuera no existe.
- **`i <= 3`** es la condición para seguir. Si es falsa, el bucle termina y el programa continúa
  en la línea de después.
- **`i++`** es abreviatura de `i = i + 1`: súmale uno. Existe también `i--` (réstale uno).

> `i` es el nombre tradicional del contador (de *índice*). Aquí sí se acepta un nombre de una
> letra, porque vive tres líneas y todo el mundo sabe qué significa.

> **Cuidado con `<` y `<=`.** `i <= 3` da tres vueltas (1, 2, 3). `i < 3` da dos (1, 2).
> Equivocarse aquí se llama *error de uno* y es de los más frecuentes que existen.

### Forma 2 · El `for` con condición (repetir mientras)

```go
n := 8
for n != 1 {
	// ...
}
```

Sin contador y sin incremento: solo la condición. Se lee *"mientras `n` no sea 1"*.
Esto es lo que en otros lenguajes se llama `while`.

> **Peligro:** si dentro del cuerpo nunca cambias `n`, la condición nunca se hace falsa y el
> programa se queda dando vueltas para siempre. Es un **bucle infinito**. Si te pasa, corta
> con **Ctrl+C** en la terminal.

### Forma 3 · El `for` infinito

```go
for {
	// ...
}
```

Sin nada. Repite eternamente hasta que algo lo corte desde dentro. Se usa en servidores
(semana 2), que están siempre a la espera.

---

## Dos frenos: `break` y `continue`

Dentro de cualquier bucle:

| Palabra | Qué hace |
|---|---|
| `break` | **Sale del bucle** entero. Se acabó, a la línea de después. |
| `continue` | **Salta el resto de esta vuelta** y va directo a la siguiente. |

```go
for i := 1; i <= 5; i++ {
	if i == 3 {
		continue      // la vuelta del 3 no imprime nada
	}
	fmt.Println(i)
}
```

```
1
2
4
5
```

La diferencia en una frase: **`break` apaga la máquina, `continue` se salta una pieza.**

---

## `if` — decidir

Ya lo usaste ayer. Dos cosas nuevas hoy.

### `else if` para encadenar

```go
if nota < 50 {
	// ...
} else if nota < 70 {
	// ...
} else {
	// ...
}
```

Se evalúa de arriba abajo y **se queda en la primera que se cumple**. Las demás ni se miran.
Por eso `nota < 70` no necesita decir "y además mayor o igual que 50": si lo fuera, habría
salido en la anterior.

### `if` con inicializador

```go
if nota := 85; nota >= 70 {
	fmt.Println("aprobado")
}
```

Punto y coma en medio: **primero crea la variable, luego pregunta**. `nota` vive **solo dentro
del `if`**; fuera no existe. Sirve para no dejar variables sueltas que ya no hacen falta.
Lo verás muchísimo en Go a partir del día 5, con los errores.

---

## `switch` — decidir entre muchos

Cuando hay cuatro o cinco caminos, la escalera de `else if` se hace ilegible. Para eso está
`switch`. En Go tiene **dos formas**.

### Forma A · `switch` sobre un valor

```go
switch mes {
case 1, 2, 3:
	return 1
case 4, 5, 6:
	return 2
default:
	return 0
}
```

Compara `mes` contra cada `case`. **Un `case` puede llevar varios valores separados por comas**
— eso ahorra un montón de líneas. `default` es "si no coincidió ninguno".

### Forma B · `switch` sin valor (para rangos)

```go
switch {
case nota >= 90:
	return "sobresaliente"
case nota >= 70:
	return "notable"
default:
	return "suspenso"
}
```

Sin nada detrás de `switch`, cada `case` es una **pregunta completa** de verdadero/falso.
Es un `else if` mejor peinado, y es la forma que se usa para rangos de números.

> **En Go no se escribe `break` dentro del `switch`.** Si vienes de ver código de C, Java o
> JavaScript te sonará raro: allí, si olvidas el `break`, la ejecución **se cae** al caso
> siguiente y produce errores absurdos. Go corta solo al terminar cada `case`. Un problema
> menos que existe.

---

## Un operador que ya conoces, con un uso nuevo

`%` (el resto de una división entera) es la herramienta estándar para preguntar
**"¿es múltiplo de?"**:

```go
n % 2 == 0    // ¿es par?      (el resto de dividir entre 2 es cero)
n % 3 == 0    // ¿múltiplo de 3?
n % 30 == 0   // ¿múltiplo de 30?
```

Si algo se divide exacto, no sobra nada. Ese "no sobra nada" es `== 0`.

---

# Parte 2 · Los ejercicios

Tres, de menos a más. Uno por carpeta, como siempre.

---

## Ejercicio A · Tabla de multiplicar

**Carpeta:** `cmd\dia02a\main.go`
**Practica:** el `for` clásico y una variable acumuladora.

### Qué tienes que hacer

Imprime la tabla del 7, del 1 al 10. Al final, la suma de todos los resultados.

### Salida exacta

```
7 x 1 = 7
7 x 2 = 14
7 x 3 = 21
7 x 4 = 28
7 x 5 = 35
7 x 6 = 42
7 x 7 = 49
7 x 8 = 56
7 x 9 = 63
7 x 10 = 70
Suma total: 385
```

### Pistas

- Un solo `for` con contador de 1 a 10. Dentro, un `Printf` con tres huecos `%d`.
- Para la suma necesitas un **acumulador**: una variable creada **antes** del bucle, a la que
  vas sumando en cada vuelta.

  ```go
  suma := 0                 // fuera del bucle, si no se reinicia cada vez
  for i := 1; i <= 10; i++ {
      suma = suma + 7*i     // o su abreviatura: suma += 7*i
  }
  ```

  **Tiene que declararse fuera.** Si pones `suma := 0` dentro, se crea de nuevo en cada vuelta
  y siempre acabará valiendo lo último que sumaste. Es un error clásico: pruébalo a propósito
  para ver qué sale.
- `suma += x` es abreviatura de `suma = suma + x`. Existen también `-=`, `*=`, `/=`.

### Extra (opcional)

Que los números queden alineados a la derecha, así:

```
7 x  1 =   7
7 x 10 =  70
```

Se consigue con `%2d` y `%3d`: el número indica **el ancho mínimo** en caracteres, rellenando
con espacios por delante.

---

## Ejercicio B · Clasificador de notas

**Carpeta:** `cmd\dia02b\main.go`
**Practica:** `for` con salto de 10, `switch` por rangos, `continue`.

### Qué tienes que hacer

1. Un `for` que recorra de **0 a 100, de 10 en 10**.
2. Una función `clasificarNota(n int) string` con un **`switch` sin valor** que devuelva:

   | Nota | Texto |
   |---|---|
   | menor que 50 | `suspenso` |
   | menor que 70 | `aprobado` |
   | menor que 90 | `notable` |
   | 90 o más | `sobresaliente` |

3. **Sáltate los múltiplos de 30** con `continue`: no deben imprimirse.

### Salida exacta

```
Nota 10 -> suspenso
Nota 20 -> suspenso
Nota 40 -> suspenso
Nota 50 -> aprobado
Nota 70 -> notable
Nota 80 -> notable
Nota 100 -> sobresaliente
```

### Pistas

- Para avanzar de 10 en 10, la tercera parte del `for` no es `n++` sino `n += 10`.
- **Mira bien la salida esperada: el 0 no está.** ¿Por qué? Porque `0 % 30` es `0`, así que
  **cero es múltiplo de 30** y tu `continue` lo salta. Es correcto, pero conviene que entiendas
  por qué pasa en vez de descubrirlo por sorpresa. Faltan también el 30, el 60 y el 90.
- El `continue` va **al principio del cuerpo**, antes de imprimir. Si lo pones después ya no
  sirve de nada: el trabajo ya se hizo.
- Fíjate en el orden de los `case` del `switch`. Si pusieras `case n < 90` el primero, un 10
  también lo cumpliría y saldría "notable". **De más restrictivo a menos.**

---

## Ejercicio C · FizzBuzz y trimestres

**Carpeta:** `cmd\dia02c\main.go`
**Practica:** separar cálculo de impresión, `%` como test de múltiplos, `switch` con varios
valores por `case`.

Este es el ejercicio de entrevista de trabajo más famoso del mundo. Se usa para descartar a
quien no sabe programar nada. Hoy lo resuelves tú.

### Parte 1 · FizzBuzz

Función `fizzbuzz(n int) string` que devuelva:

| Condición | Devuelve |
|---|---|
| múltiplo de 3 **y** de 5 | `FizzBuzz` |
| múltiplo de 3 | `Fizz` |
| múltiplo de 5 | `Buzz` |
| ninguna | el número, como texto |

En `main`, imprímelo del 1 al 30.

### Parte 2 · Trimestres

Función `trimestre(mes int) int` que, con un **`switch` sobre el valor** y `case` de valores
múltiples, devuelva a qué trimestre pertenece un mes (1-3 → 1, 4-6 → 2, 7-9 → 3, 10-12 → 4).

Pruébala con los meses 2, 5 y 11.

### Salida exacta

```
1
2
Fizz
4
Buzz
Fizz
7
8
Fizz
Buzz
11
Fizz
13
14
FizzBuzz
16
17
Fizz
19
Buzz
Fizz
22
23
Fizz
Buzz
26
Fizz
28
29
FizzBuzz
Mes 2 -> trimestre 1
Mes 5 -> trimestre 2
Mes 11 -> trimestre 4
```

### Pistas

- **El orden de las comprobaciones es todo el ejercicio.** Si preguntas primero por el 3, el
  número 15 devolverá `Fizz` y nunca llegará a `FizzBuzz`. **El caso más exigente va primero.**
- "Múltiplo de 3 y de 5" se puede escribir de dos formas: `n%3 == 0 && n%5 == 0`, o más corto,
  `n%15 == 0` (lo que es múltiplo de ambos es múltiplo de 15). `&&` significa "y".
- **La función devuelve `string`, así que el número también hay que devolverlo como texto.**
  Aquí vas a chocar con algo: `return n` no compila, `n` es `int`. Necesitas
  `fmt.Sprintf("%d", n)` — es igual que `Printf`, pero en vez de imprimir en pantalla,
  **te devuelve el texto**. La `S` es de *string*. Te va a resultar utilísimo.
- Que la función **devuelva** en vez de imprimir es a propósito: así se podrá probar con un
  test (día 6). Una función que imprime por dentro no se puede comprobar automáticamente.
- Para `trimestre` usa la **forma A** del `switch` (sobre el valor), con `case 1, 2, 3:`.
  Es el caso donde luce.

---

## Ejecutar

```powershell
go run ./cmd/dia02a
go run ./cmd/dia02b
go run ./cmd/dia02c
```

Cuando tengas alguno, pídeme que lo revise. Y si un bucle se te queda colgado imprimiendo
sin parar, **Ctrl+C**: has hecho tu primer bucle infinito, que es un rito de paso.
