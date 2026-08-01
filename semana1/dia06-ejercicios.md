# Día 6 · Paquetes y tests — explicación y ejercicios

Hasta hoy todo tu código ha vivido en un `main.go` que se ejecuta y ya. Hoy cambian dos cosas
y las dos son de "programador de verdad":

1. **Paquetes:** sacar código de `main` para poder reutilizarlo.
2. **Tests:** que el propio Go compruebe que tus funciones hacen lo que dices.

---

# Parte 1 · La explicación

## Paquetes

Un **paquete** es una carpeta con código reutilizable. Hasta ahora solo has usado los que vienen
con Go (`fmt`, `errors`, `strconv`). Hoy escribes el tuyo.

### Dónde va

```
aprendiendo-go\
├── cmd\dia06a\main.go        ← programa (package main)
└── internal\stats\stats.go   ← librería  (package stats)
```

- El **nombre del paquete** es el de la carpeta: `internal\stats\` → `package stats`.
- Puede haber varios archivos `.go` en la carpeta; todos llevan `package stats` y **se ven entre
  ellos** sin importar nada.
- **`internal`** es una carpeta con magia: Go impide que nadie fuera de tu módulo importe lo que
  hay dentro. Es la forma de decir "esto es mío, no es una librería pública".

### Cómo se importa

```go
import "aprendiendo-go/internal/stats"
```

La ruta empieza por el **nombre del módulo** (está en `go.mod`: `aprendiendo-go`), y sigue con las
carpetas separadas por `/`, incluso en Windows.

Y se usa poniendo el nombre del paquete delante:

```go
stats.Media(numeros)
```

### La mayúscula, otra vez — y ahora es de verdad

Te lo he repetido varios días. Aquí es donde importa:

> **Mayúscula inicial = visible desde fuera del paquete. Minúscula = solo dentro.**

```go
func Media(ns []float64) (float64, error)   // se puede llamar desde main
func sumar(ns []float64) float64            // ayudante privado, invisible fuera
```

Esto no es un convenio de estilo: **es el mecanismo del lenguaje**. No hay `public` ni `private`
en Go, hay mayúsculas.

---

## Tests

Un **test** es código que comprueba tu código. Lo escribes una vez y Go lo ejecuta cuantas veces
quieras. Sirve para dos cosas:

- Saber que **hoy** funciona.
- Saber que sigue funcionando **después** de tocarlo. Esto es lo valioso: te deja cambiar código
  sin miedo.

### Las reglas (Go es muy estricto aquí)

1. El archivo se llama `algo_test.go`. **El sufijo `_test` es obligatorio.**
2. Va **en la misma carpeta** que el código que prueba.
3. Las funciones se llaman `TestLoQueSea`, **con `Test` en mayúscula**.
4. Reciben un único parámetro: `t *testing.T`.

```go
package stats

import "testing"

func TestMedia(t *testing.T) {
	resultado, err := Media([]float64{2, 4, 6})
	if err != nil {
		t.Fatalf("no esperaba error: %v", err)
	}
	if resultado != 4 {
		t.Errorf("Media = %v; quería 4", resultado)
	}
}
```

Fíjate en que **no hay `main`**. Los tests no se ejecutan con `go run`, sino con:

```powershell
go test ./...
```

`./...` significa "esta carpeta y todas las de dentro".

### `t.Errorf` vs `t.Fatalf`

| | Qué hace |
|---|---|
| `t.Errorf` | marca el test como fallido **y sigue** |
| `t.Fatalf` | marca fallido y **para ese test en seco** |

Usa `Fatalf` cuando seguir no tenga sentido (por ejemplo, si hubo un error inesperado y el
resultado es basura). Usa `Errorf` cuando quieras ver todos los fallos de golpe.

### El mensaje de fallo importa

```go
t.Errorf("Media(%v) = %v; quería %v", entrada, obtenido, esperado)
```

Cuando un test falle dentro de seis meses, ese mensaje es **todo lo que vas a tener**. Que diga
qué entró, qué salió y qué esperabas. La fórmula "obtenido; quería" es la costumbre en Go.

---

## Tests de tabla — la forma idiomática

Si quieres probar cinco casos, no escribas cinco funciones casi iguales. Pon los casos en un
slice y recórrelos:

```go
func TestMedia(t *testing.T) {
	casos := []struct {
		nombre   string
		entrada  []float64
		esperado float64
		conError bool
	}{
		{"tres números", []float64{2, 4, 6}, 4, false},
		{"un elemento", []float64{7}, 7, false},
		{"vacío", []float64{}, 0, true},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			// ...
		})
	}
}
```

Dos cosas nuevas:

- **El struct anónimo.** `[]struct{...}{...}` define el tipo y crea el slice de una vez, sin
  ponerle nombre. Se usa así porque ese tipo no sirve para nada fuera del test.
- **`t.Run(nombre, ...)`** crea un **subtest**. Cada caso aparece con su nombre en la salida, y
  si falla sabes cuál sin contar posiciones.

**Ventaja real:** añadir un caso nuevo es añadir **una línea** a la tabla. Por eso en Go casi
todos los tests son así.

### Qué casos hay que probar

No los que sabes que funcionan. **Los bordes:**

- La lista **vacía** (¿divide entre cero?).
- **Un solo** elemento.
- Números **negativos**.
- Cantidad **par** e **impar** (crítico para la mediana).

Ahí es donde viven los bugs.

---

## Cómo leer la salida de `go test`

Bien:

```
ok      aprendiendo-go/internal/stats   0.234s
```

Mal:

```
--- FAIL: TestMedia (0.00s)
    --- FAIL: TestMedia/vacío (0.00s)
        stats_test.go:31: Media([]) = 0; quería un error
FAIL
```

Te da el test, el subtest, el **archivo y la línea** exactos, y tu mensaje. Con `-v` (*verbose*)
ves también los que pasan:

```powershell
go test -v ./...
```

---

# Parte 2 · Los ejercicios

Los tres construyen **el mismo paquete**, que además es la pieza que usarás en el proyecto de
mañana. No los saltes.

---

## Ejercicio A · Tu primer paquete

**Carpeta:** `internal\stats\stats.go` y `cmd\dia06a\main.go`
**Practica:** crear un paquete, importarlo, mayúscula vs minúscula.

### Qué tienes que hacer

1. Crea `internal\stats\stats.go` con `package stats`.
2. `func Media(ns []float64) (float64, error)` — error si el slice está vacío.
   Error centinela `ErrVacio` (día 5).
3. Un ayudante **privado** `func sumar(ns []float64) float64`.
4. En `cmd\dia06a\main.go`, importa el paquete y prueba con una lista normal y con una vacía.

### Salida exacta

```
Media de [2 4 6 8]: 5.00
Media de []: error -> el slice está vacío
```

### Pistas

- El nombre del módulo es `aprendiendo-go` (mira `go.mod`), así que el import es
  `"aprendiendo-go/internal/stats"`. Con `/`, no con `\`.
- **Prueba a llamar a `stats.sumar(...)` desde `main`** y lee el error:

  ```
  undefined: stats.sumar
  ```

  No es que no exista: es que está en minúscula y no sale de su paquete. Ese experimento te
  ahorrará dudas para siempre.
- Recuerda `float64(len(ns))` para dividir.
- La lista vacía se escribe `[]float64{}`.

---

## Ejercicio B · El test de tabla

**Carpeta:** `internal\stats\stats_test.go`
**Practica:** escribir tests, subtests, casos límite.

### Qué tienes que hacer

Un `TestMedia` de tabla con al menos estos casos:

| Nombre | Entrada | Esperado |
|---|---|---|
| `varios` | `{2, 4, 6, 8}` | `5` |
| `un elemento` | `{7}` | `7` |
| `negativos` | `{-2, -4}` | `-3` |
| `mezclados` | `{-10, 10}` | `0` |
| `vacío` | `{}` | error |

### Salida esperada

```powershell
go test ./internal/stats
```

```
ok      aprendiendo-go/internal/stats   0.003s
```

(El tiempo variará; lo que importa es el `ok`.)

Con `-v`:

```
=== RUN   TestMedia
=== RUN   TestMedia/varios
=== RUN   TestMedia/un_elemento
...
--- PASS: TestMedia (0.00s)
PASS
```

> Fíjate en que Go convierte los espacios del nombre en `_`.

### Pistas

- El campo `conError bool` en la tabla te dice si ese caso debe fallar. Dentro del subtest:

  ```go
  if c.conError {
      // aquí err DEBE ser distinto de nil
      return
  }
  // aquí err DEBE ser nil
  ```

- **Comprueba siempre el error antes que el valor.** Si hubo error, el valor es relleno y
  compararlo no significa nada.
- **Rompe tu código a propósito.** Cambia el `+` de `sumar` por un `-`, corre el test y mira cómo
  falla. Un test que nunca has visto fallar no te consta que funcione.

> ⚠️ **Sobre comparar decimales.** Los casos de arriba están elegidos para dar números exactos.
> En general, comparar `float64` con `==` es peligroso: `0.1 + 0.2` no es exactamente `0.3` en
> ningún lenguaje, por cómo se guardan los decimales en binario. Lo correcto es comprobar que la
> diferencia sea diminuta. Hoy no hace falta, pero que lo sepas.

---

## Ejercicio C · Completar el paquete

**Carpeta:** `internal\stats\` (los dos archivos)
**Practica:** más funciones, el caso par/impar de la mediana, ampliar la tabla.

### Qué tienes que hacer

Añade a `stats`, cada una con su error para slice vacío:

- `Max(ns []float64) (float64, error)`
- `Min(ns []float64) (float64, error)`
- `Mediana(ns []float64) (float64, error)`

Y amplía los tests para cubrirlas.

### Salida exacta de `cmd\dia06c\main.go`

```
Datos: [5 3 9 1 7]
Media:    5.00
Mediana:  5.00
Máximo:   9.00
Mínimo:   1.00

Datos: [5 3 9 1]
Mediana:  4.00
```

### Pistas — la mediana tiene truco

La mediana es "el valor de en medio **una vez ordenados**". Dos avisos:

**1. Hay que ordenar primero.** Con `slices.Sort(copia)` (paquete `slices`, ya viene con Go).

**2. Y aquí está la trampa gorda:** ordenar **modifica el slice original**. Un slice es una
ventana a los datos, no una copia. Si `Mediana` ordena la lista que le pasaron, el `main` se
encuentra sus datos cambiados sin haber hecho nada. Eso es un efecto secundario y es de los bugs
más desagradables que hay.

Haz una copia antes:

```go
copia := make([]float64, len(ns))
copy(copia, ns)
slices.Sort(copia)
```

`make` reserva un slice nuevo del tamaño indicado; `copy` vuelca los valores. **Compruébalo:**
imprime el slice original después de llamar a `Mediana` con y sin copia.

**3. Par vs impar:**

- **Impar** (5 elementos): el de en medio, posición `len/2`. Con 5 → posición 2 → el tercero.
- **Par** (4 elementos): no hay uno en medio; es el **promedio de los dos centrales**,
  posiciones `len/2 - 1` y `len/2`.

Por eso la tabla de tests **tiene que llevar los dos casos**. Con `[5 3 9 1 7]` ordenado
(`1 3 5 7 9`) la mediana es 5. Con `[5 3 9 1]` ordenado (`1 3 5 9`) es (3+5)/2 = 4.

### Pistas para `Max` y `Min`

- Es el patrón del acumulador del día 2, pero comparando. **No empieces en 0**: si todos los
  números son negativos, `Max` te devolvería 0, que no está en la lista. **Empieza en el primer
  elemento** (`ns[0]`) y compara desde el segundo.
- Ese caso —"todos negativos"— es exactamente el que debe estar en tu tabla de tests.

---

## Comandos

```powershell
go test ./...                    # todos los tests del módulo
go test -v ./internal/stats      # con detalle de cada subtest
go test -run TestMediana ./...   # solo los que coincidan con ese nombre
go run ./cmd/dia06a
go run ./cmd/dia06c
```

## Lo que tienes que poder explicar al terminar

- Por qué `sumar` no se ve desde `main` y `Media` sí.
- Por qué un test de tabla es mejor que cinco funciones de test.
- Por qué `Mediana` tiene que copiar el slice antes de ordenarlo.
- Por qué `Max` no puede empezar en 0.
