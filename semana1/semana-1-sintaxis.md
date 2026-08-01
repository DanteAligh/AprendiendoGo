# Semana 1 · Sintaxis y programas de terminal

**Meta de la semana:** dominar la sintaxis básica de Go y terminar con un CLI que lee un CSV y calcula estadísticas.

**Carpeta de trabajo:** `cmd\stats\` e `internal\csv\`

Cada día crea el archivo indicado, ejecútalo con el comando que aparece, y luego resuelve el ejercicio.

---

## Día 1 · Variables, tipos y funciones

**Conceptos:** `var` vs `:=`, tipos básicos, funciones con múltiples retornos, `fmt.Printf`.

Crea `cmd\dia01\main.go`:

```go
package main

import "fmt"

// Go infiere tipos con :=, pero var es obligatorio a nivel de paquete.
var version = "1.0"

// Una función puede devolver varios valores. Es la base del manejo
// de errores en Go: (resultado, error).
func dividir(a, b float64) (float64, bool) {
	if b == 0 {
		return 0, false
	}
	return a / b, true
}

func main() {
	var nombre string = "Go"
	edad := 16 // int por inferencia
	pi := 3.14159

	fmt.Printf("Lenguaje: %s, edad: %d, pi: %.2f, version: %s\n", nombre, edad, pi, version)

	resultado, ok := dividir(10, 3)
	if !ok {
		fmt.Println("no se puede dividir entre cero")
		return
	}
	fmt.Printf("10/3 = %.4f\n", resultado)
}
```

Ejecuta:

```powershell
go run ./cmd/dia01
```

**Verbos de `Printf` que usarás siempre:** `%s` cadena, `%d` entero, `%f` decimal (`%.2f` con dos cifras), `%v` cualquier valor, `%+v` struct con nombres de campo, `%T` el tipo, `%q` cadena entrecomillada.

### Ejercicio

Escribe una función `celsiusAFahrenheit(c float64) float64` y otra `esFiebre(tempC float64) bool` (fiebre a partir de 38 °C). En `main`, recorre las temperaturas `36.5, 37.2, 38.4, 39.0` y por cada una imprime los grados en ambas escalas y si hay fiebre.

> Pista: la conversión es `c*9/5 + 32`. Ojo con la división entera: si escribes `9/5` con enteros da `1`. Usa `9.0/5.0` o multiplica primero.

---

## Día 2 · Control de flujo

**Conceptos:** `if` con inicializador, `for` (la única palabra de bucle en Go), `switch` sin `break`.

Crea `cmd\dia02\main.go`:

```go
package main

import (
	"fmt"
	"strings"
)

func clasificar(nota int) string {
	// switch en Go no necesita break: no hay caída entre casos.
	switch {
	case nota >= 90:
		return "excelente"
	case nota >= 70:
		return "aprobado"
	default:
		return "reprobado"
	}
}

func main() {
	// for clásico
	for i := 1; i <= 3; i++ {
		fmt.Println("vuelta", i)
	}

	// for como while
	n := 8
	pasos := 0
	for n != 1 {
		if n%2 == 0 {
			n = n / 2
		} else {
			n = 3*n + 1
		}
		pasos++
	}
	fmt.Println("pasos de Collatz:", pasos)

	// if con inicializador: la variable solo vive dentro del if
	if nota := 85; nota >= 70 {
		fmt.Println("resultado:", clasificar(nota))
	}

	fmt.Println(strings.Repeat("-", 20))
}
```

Ejecuta:

```powershell
go run ./cmd/dia02
```

### Ejercicio

Programa el clásico FizzBuzz del 1 al 30, pero devolviendo la cadena desde una función `fizzbuzz(n int) string` en vez de imprimir dentro del bucle. Luego añade un `switch` con `case` de valores múltiples (`case 1, 2, 3:`) que clasifique un mes en trimestre.

---

## Día 3 · Slices y maps

**Conceptos:** slices (crecen, a diferencia de los arrays), `append`, `len` y `cap`, recorrer con `range`, maps y la comprobación de existencia.

Crea `cmd\dia03\main.go`:

```go
package main

import (
	"fmt"
	"sort"
)

func main() {
	// Slice: tamaño dinámico. append puede reasignar memoria,
	// por eso SIEMPRE hay que reasignar el resultado.
	numeros := []int{5, 2, 9, 1}
	numeros = append(numeros, 7)
	fmt.Println(numeros, "len:", len(numeros), "cap:", cap(numeros))

	// range da índice y valor. Usa _ para descartar el que no necesites.
	suma := 0
	for _, v := range numeros {
		suma += v
	}
	fmt.Println("suma:", suma)

	sort.Ints(numeros)
	fmt.Println("ordenado:", numeros)

	// Map: hay que inicializarlo con make o con literal.
	conteo := map[string]int{}
	palabras := []string{"go", "rust", "go", "python", "go"}
	for _, p := range palabras {
		conteo[p]++ // el valor cero de int es 0, no hace falta inicializar
	}

	// La forma de dos valores distingue "existe con valor cero"
	// de "no existe".
	if n, existe := conteo["java"]; existe {
		fmt.Println("java aparece", n)
	} else {
		fmt.Println("java no aparece")
	}

	// El orden de recorrido de un map es ALEATORIO por diseño.
	// Si necesitas orden, extrae las claves y ordénalas.
	claves := make([]string, 0, len(conteo))
	for k := range conteo {
		claves = append(claves, k)
	}
	sort.Strings(claves)
	for _, k := range claves {
		fmt.Printf("%-8s %d\n", k, conteo[k])
	}
}
```

Ejecuta:

```powershell
go run ./cmd/dia03
```

### Ejercicio

Escribe `estadisticas(datos []float64) (min, max, media float64)`. Debe manejar el slice vacío devolviendo ceros sin entrar en pánico. Después, escribe `moda(datos []int) int` usando un map de frecuencias.

> Trampa a evitar: `var s []int` crea un slice **nil**, pero `append` funciona igual sobre él. En cambio un map nil revienta al escribir — los maps hay que crearlos con `make` o literal.

---

## Día 4 · Structs, métodos e interfaces

**Conceptos:** `struct`, receptores por valor vs por puntero, interfaces implícitas.

Crea `cmd\dia04\main.go`:

```go
package main

import (
	"fmt"
	"math"
)

type Figura interface {
	Area() float64
	Nombre() string
}

type Rectangulo struct {
	Ancho, Alto float64
}

type Circulo struct {
	Radio float64
}

// Receptor por VALOR: recibe una copia, no puede modificar el original.
func (r Rectangulo) Area() float64  { return r.Ancho * r.Alto }
func (r Rectangulo) Nombre() string { return "rectángulo" }

func (c Circulo) Area() float64  { return math.Pi * c.Radio * c.Radio }
func (c Circulo) Nombre() string { return "círculo" }

// Receptor por PUNTERO: puede modificar el struct original.
func (r *Rectangulo) Escalar(factor float64) {
	r.Ancho *= factor
	r.Alto *= factor
}

func main() {
	// Nadie declara "implements". Si tiene los métodos, cumple la interfaz.
	figuras := []Figura{
		Rectangulo{Ancho: 3, Alto: 4},
		Circulo{Radio: 2},
	}

	for _, f := range figuras {
		fmt.Printf("%-12s área %.2f\n", f.Nombre(), f.Area())
	}

	r := Rectangulo{Ancho: 2, Alto: 5}
	r.Escalar(3) // Go hace &r automáticamente
	fmt.Printf("%+v área %.2f\n", r, r.Area())
}
```

Ejecuta:

```powershell
go run ./cmd/dia04
```

**Regla práctica sobre receptores:** usa puntero si el método modifica el struct o si el struct es grande. Y sé consistente: si un método usa puntero, que todos lo usen.

### Ejercicio

Define `type Empleado struct { Nombre string; Salario float64; Antiguedad int }` y un método `Aumentar(pct float64)` con receptor por puntero. Crea una interfaz `Describible` con `Describir() string` e impleméntala. Luego prueba a cambiar `Aumentar` a receptor por valor y observa que el salario ya no cambia — ese error lo cometerás alguna vez, mejor verlo hoy.

---

## Día 5 · Errores y `defer`

**Conceptos:** `error` como valor, `errors.New`, `fmt.Errorf` con `%w`, `errors.Is`, `defer`.

Crea `cmd\dia05\main.go`:

```go
package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

// Los errores centinela se declaran como variables de paquete
// para poder compararlos con errors.Is.
var ErrNegativo = errors.New("el valor no puede ser negativo")

func raiz(x float64) (float64, error) {
	if x < 0 {
		// %w "envuelve" el error: conserva la cadena para errors.Is.
		return 0, fmt.Errorf("raiz(%v): %w", x, ErrNegativo)
	}
	r := x
	for i := 0; i < 20; i++ {
		r = (r + x/r) / 2
	}
	return r, nil
}

func leerNumero(texto string) (int, error) {
	n, err := strconv.Atoi(texto)
	if err != nil {
		return 0, fmt.Errorf("no pude convertir %q: %w", texto, err)
	}
	return n, nil
}

func main() {
	// defer se ejecuta al salir de la función, pase lo que pase.
	// Es lo que garantiza que los recursos se liberen.
	defer fmt.Println("-- fin del programa --")

	if r, err := raiz(9); err == nil {
		fmt.Printf("raíz de 9 ≈ %.4f\n", r)
	}

	_, err := raiz(-4)
	if errors.Is(err, ErrNegativo) {
		fmt.Println("error esperado:", err)
	}

	if n, err := leerNumero("42x"); err != nil {
		fmt.Fprintln(os.Stderr, "fallo:", err)
	} else {
		fmt.Println("número:", n)
	}
}
```

Ejecuta:

```powershell
go run ./cmd/dia05
```

**El patrón que verás mil veces:** llama, comprueba `err != nil`, devuelve envolviendo con contexto. No lo escondas en un helper genérico; la verbosidad aquí es intencional.

### Ejercicio

Escribe `abrirYContar(ruta string) (int, error)` que abra un archivo con `os.Open`, ponga `defer f.Close()` inmediatamente después de comprobar el error, y cuente las líneas con `bufio.Scanner`. Prueba con un archivo que no exista y comprueba el error con `errors.Is(err, os.ErrNotExist)`.

---

## Día 6 · Tests de tabla

**Conceptos:** `go test`, tests de tabla, subtests con `t.Run`, cobertura.

Crea `internal\stats\stats.go`:

```go
package stats

import "errors"

var ErrVacio = errors.New("conjunto de datos vacío")

func Media(datos []float64) (float64, error) {
	if len(datos) == 0 {
		return 0, ErrVacio
	}
	var suma float64
	for _, v := range datos {
		suma += v
	}
	return suma / float64(len(datos)), nil
}
```

Y junto a él `internal\stats\stats_test.go` (el sufijo `_test.go` es obligatorio):

```go
package stats

import (
	"errors"
	"math"
	"testing"
)

func TestMedia(t *testing.T) {
	casos := []struct {
		nombre  string
		entrada []float64
		quiero  float64
		wantErr error
	}{
		{"caso normal", []float64{1, 2, 3, 4}, 2.5, nil},
		{"un elemento", []float64{7}, 7, nil},
		{"negativos", []float64{-2, 2}, 0, nil},
		{"vacío", nil, 0, ErrVacio},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			got, err := Media(c.entrada)
			if !errors.Is(err, c.wantErr) {
				t.Fatalf("error = %v, quiero %v", err, c.wantErr)
			}
			// Nunca compares floats con == : usa una tolerancia.
			if math.Abs(got-c.quiero) > 1e-9 {
				t.Errorf("Media(%v) = %v, quiero %v", c.entrada, got, c.quiero)
			}
		})
	}
}
```

Ejecuta:

```powershell
go test ./...
go test -v ./internal/stats
go test -cover ./...
```

### Ejercicio

Añade `Mediana` y `DesviacionEstandar` al paquete, cada una con su test de tabla. Incluye a propósito un caso que falle, mira cómo se lee la salida de `go test -v`, y luego arréglalo.

---

## Día 7 · Proyecto: CLI de estadísticas

Junta todo. Crea `cmd\stats\main.go`:

```go
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"

	"aprendiendo-go/internal/stats"
)

func main() {
	ruta := flag.String("archivo", "", "ruta del CSV a analizar")
	columna := flag.Int("col", 0, "índice de la columna numérica (base 0)")
	cabecera := flag.Bool("cabecera", true, "saltar la primera fila")
	flag.Parse()

	if *ruta == "" {
		flag.Usage()
		os.Exit(2)
	}

	if err := ejecutar(*ruta, *columna, *cabecera); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func ejecutar(ruta string, col int, saltarCabecera bool) error {
	f, err := os.Open(ruta)
	if err != nil {
		return fmt.Errorf("abriendo archivo: %w", err)
	}
	defer f.Close()

	filas, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return fmt.Errorf("leyendo csv: %w", err)
	}
	if saltarCabecera && len(filas) > 0 {
		filas = filas[1:]
	}

	datos := make([]float64, 0, len(filas))
	for i, fila := range filas {
		if col >= len(fila) {
			return fmt.Errorf("fila %d: no existe la columna %d", i+1, col)
		}
		v, err := strconv.ParseFloat(fila[col], 64)
		if err != nil {
			return fmt.Errorf("fila %d: %w", i+1, err)
		}
		datos = append(datos, v)
	}

	media, err := stats.Media(datos)
	if err != nil {
		return err
	}
	fmt.Printf("registros: %d\nmedia:     %.4f\n", len(datos), media)
	return nil
}
```

Crea un `ventas.csv` de prueba en la raíz:

```
producto,unidades
teclado,120
raton,85
monitor,42
```

Ejecuta y compila:

```powershell
go run ./cmd/stats -archivo ventas.csv -col 1
go build -o bin/stats.exe ./cmd/stats
.\bin\stats.exe -archivo ventas.csv -col 1
```

### Ejercicio final de la semana

1. Añade los flags `-min`, `-max` y `-mediana` para elegir qué estadísticas mostrar.
2. Permite seleccionar la columna por nombre además de por índice (`-col nombre`).
3. Si no se pasa `-archivo`, lee de la entrada estándar (`os.Stdin`) para poder hacer `type ventas.csv | .\bin\stats.exe -col 1`.
4. Ejecuta `go vet ./...` y `gofmt -l .` y deja ambos sin salida.

---

## Comandos de referencia

| Comando | Qué hace |
|---|---|
| `go run ./cmd/x` | Compila y ejecuta sin dejar binario |
| `go build -o bin/x.exe ./cmd/x` | Genera el ejecutable |
| `go test ./...` | Corre todos los tests del módulo |
| `go test -run TestMedia ./...` | Solo los tests que coincidan |
| `go fmt ./...` | Formatea (no hay debate de estilo en Go) |
| `go vet ./...` | Detecta errores comunes que compilan |
| `go doc fmt.Printf` | Documentación desde la terminal |

## Checklist antes de pasar a la semana 2

- [ ] Distingo cuándo usar receptor por valor y por puntero
- [ ] Sé por qué `append` debe reasignarse siempre
- [ ] Escribo un test de tabla sin copiar la plantilla
- [ ] Envuelvo errores con `%w` y los compruebo con `errors.Is`
- [ ] Entiendo qué hace `defer` y dónde colocarlo
