# Día 7 · Proyecto — CLI de estadísticas

Hoy no hay conceptos sueltos: hay **un programa de verdad**. Uno que un compañero podría usar sin
saber que lo escribiste tú.

Junta todo lo de la semana: funciones y tipos (día 1), bucles y `switch` (día 2), slices (día 3),
structs (día 4), errores y `defer` (día 5), y el paquete `stats` con sus tests (día 6).

**Objetivo:**

```powershell
go run ./cmd/stats -archivo ventas.csv -columna 2
```

```
Archivo:  ventas.csv
Columna:  2 (importe)
Filas:    5

Media:    286.00
Mediana:  250.00
Máximo:   500.00
Mínimo:   120.00
```

---

# Parte 1 · La explicación

## Qué es una CLI

**CLI** = *Command Line Interface*, interfaz de línea de comandos. Un programa sin ventanas: se
lanza desde la terminal, recibe sus datos como **argumentos** y responde escribiendo texto.

`go`, `git` y `docker` son CLIs. Hoy escribes la tuya.

## Flags — los argumentos con nombre

Un **flag** (bandera) es una opción con nombre: `-archivo ventas.csv`. Se llaman así porque
"levantas una banderita" para activar algo.

Podrías leer los argumentos a mano de `os.Args`, pero Go trae el paquete `flag`:

```go
import "flag"

func main() {
	archivo := flag.String("archivo", "", "ruta del CSV a analizar")
	columna := flag.Int("columna", 1, "número de columna (empezando en 1)")
	flag.Parse()

	fmt.Println(*archivo, *columna)
}
```

Tres cosas:

**1. Los tres parámetros son:** nombre del flag, **valor por defecto**, y descripción de la ayuda.

**2. `flag.Parse()` es obligatorio.** Hasta que no lo llamas, las variables tienen el valor por
defecto. Se llama una sola vez, al principio de `main`.

**3. Devuelven punteros — de ahí el asterisco.** `flag.String` no devuelve un `string`, devuelve
un `*string`: **la dirección** donde estará el valor cuando se llame a `Parse()`. Para leer el
valor que hay en esa dirección se pone `*` delante: `*archivo`.

> Es el mismo puntero del día 4, visto por el otro lado. Allí el `*` estaba en el receptor de un
> método para poder **modificar** el original. Aquí `flag` necesita darte la dirección **antes**
> de conocer el valor, porque el valor solo se sabe al llamar a `Parse()`.
>
> Si te resulta incómodo, existe la variante que escribe sobre una variable tuya:
> ```go
> var archivo string
> flag.StringVar(&archivo, "archivo", "", "ruta del CSV")
> ```
> El `&` significa "la dirección de". Así luego usas `archivo` sin asterisco.

**De regalo, `-h` funciona solo:**

```
Usage of stats:
  -archivo string
        ruta del CSV a analizar
  -columna int
        número de columna (empezando en 1) (default 1)
```

## Leer un archivo

```go
f, err := os.Open(ruta)
if err != nil {
	return err
}
defer f.Close()
```

`os.Open` devuelve el archivo abierto y un error (no existe, sin permisos...). El `defer f.Close()`
va **pegado debajo**: día 5.

## Leer un CSV

Un **CSV** (*comma-separated values*) es un archivo de texto donde cada línea es una fila y las
columnas van separadas por comas:

```
fecha,producto,importe
2026-01-05,teclado,250.00
2026-01-06,ratón,120.50
```

Se podría partir por comas a mano, pero es más difícil de lo que parece (comillas, comas dentro
de un campo, saltos de línea). Go trae `encoding/csv`:

```go
lector := csv.NewReader(f)
filas, err := lector.ReadAll()
```

`filas` es `[][]string`: **una lista de filas, y cada fila una lista de textos**. Es el mismo tipo
"lista de listas" que viste en el día 1. `filas[0]` es la cabecera; `filas[1][2]` es la tercera
columna de la segunda fila.

> Todo llega como **texto**, aunque parezca un número. Convertir es cosa tuya:
> `strconv.ParseFloat` (día 5).

## Códigos de salida

Cuando un programa termina, deja un número: **0 = todo bien, cualquier otro = falló**. Nadie lo
ve en pantalla, pero es como los programas se avisan entre sí (y como un servidor de CI sabe si
tu build pasó, semana 4).

```go
fmt.Fprintln(os.Stderr, "error:", err)
os.Exit(1)
```

Dos detalles importantes:

- **`os.Stderr`.** Hay dos canales de salida: `Stdout` (el resultado) y `Stderr` (los errores).
  Existen separados para poder redirigir uno sin el otro. **Los errores van siempre por `Stderr`.**
  `fmt.Fprintln` es como `Println` pero eligiendo el canal.
- **`os.Exit(1)` mata el programa en el acto** y **no ejecuta los `defer` pendientes**. Por eso el
  patrón habitual es: `main` no hace casi nada, llama a una función `ejecutar() error`, y solo
  `main` decide salir.

```go
func main() {
	if err := ejecutar(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
```

Esa es la forma canónica de una CLI en Go, y de paso usa el `if` con inicializador del día 2.

---

# Parte 2 · El proyecto por fases

Constrúyelo en tres pasos y **comprueba cada uno antes de seguir**. Es el mismo consejo del día 1:
un caso funcionando antes de generalizar.

Todo va en `cmd\stats\main.go`.

## Fase 1 · Solo los flags

Que lea `-archivo` y `-columna` y los imprima. Nada más.

```powershell
go run ./cmd/stats -archivo ventas.csv -columna 2
```

```
archivo=ventas.csv columna=2
```

Prueba también `go run ./cmd/stats -h` y sin argumentos.

## Fase 2 · Leer el CSV

Crea `ventas.csv` en la raíz del proyecto:

```
fecha,producto,importe
2026-01-05,teclado,250.00
2026-01-06,ratón,120.00
2026-01-07,monitor,500.00
2026-01-08,cable,180.00
2026-01-09,webcam,380.00
```

Que el programa lo abra, lo lea, **se salte la cabecera** y imprima los valores de la columna
pedida:

```
[250 120 500 180 380]
```

## Fase 3 · Las estadísticas

Conecta el paquete `stats` del día 6 y da el formato final.

### Salida exacta

```
Archivo:  ventas.csv
Columna:  2 (importe)
Filas:    5

Media:    286.00
Mediana:  250.00
Máximo:   500.00
Mínimo:   120.00
```

> Comprueba la mediana a mano: ordenados son `120 180 250 380 500`; cinco valores, impar, el de
> en medio es 250. Y la media: 1430 / 5 = 286.

---

# Parte 3 · Los errores (la mitad del ejercicio)

Un programa que funciona con datos buenos está a medio hacer. **Estos cuatro casos tienen que dar
un mensaje claro y salir con código 1.**

### 1 · Falta el archivo

```powershell
go run ./cmd/stats
```

```
error: falta el flag -archivo
```

### 2 · El archivo no existe

```powershell
go run ./cmd/stats -archivo noexiste.csv -columna 2
```

```
error: abriendo noexiste.csv: open noexiste.csv: The system cannot find the file specified.
```

> La parte final la escribe el sistema operativo, en inglés y con el texto de Windows. **No la
> escribas tú a mano**: sale sola al envolver con `%w` el error de `os.Open`. Tu aportación es el
> `abriendo noexiste.csv:` de delante.

### 3 · La columna no existe

```powershell
go run ./cmd/stats -archivo ventas.csv -columna 9
```

```
error: la columna 9 no existe (el archivo tiene 3 columnas)
```

### 4 · La columna no es numérica

```powershell
go run ./cmd/stats -archivo ventas.csv -columna 1
```

```
error: fila 2, columna 1: "2026-01-05" no es un número
```

> **Fila 2** porque las personas cuentan desde 1 y la 1 es la cabecera. Que el mensaje diga la
> fila exacta es la diferencia entre un error útil y uno inútil: con 3000 filas, "no es un número"
> a secas no vale para nada.

### Comprobar el código de salida

En PowerShell, después de ejecutar:

```powershell
echo $LASTEXITCODE
```

`0` si fue bien, `1` si falló.

---

## Pistas generales

- **`-columna 2` es la segunda columna para una persona, pero la posición `1` para el slice.**
  Resta 1 al usar el índice. Este desajuste es fuente inagotable de errores; escríbelo con un
  comentario donde lo hagas.
- Valida la columna **una vez, antes del bucle**, mirando `len(filas[0])`. No dentro del bucle.
- `filas[1:]` te da todas las filas **menos la cabecera**, sin necesidad de un `if` dentro del
  bucle. Es el corte de slices del día 3.
- Estructura recomendada: `main` mínimo + `ejecutar() error` + funciones pequeñas
  (`leerColumna(ruta string, col int) ([]float64, error)`). Funciones pequeñas = fáciles de testear.
- El nombre de la columna del encabezado (`importe`) sale de `filas[0][col-1]`. No lo escribas a mano.
- Todos los errores hacia arriba con `%w`; imprimir, solo `main`.

## Extras si te sobra tiempo

- Flag `-formato` que acepte `texto` o `json`, con un `switch` (día 2).
- Que `-columna` acepte también el **nombre** (`-columna importe`) buscándolo en la cabecera.
- Tests para `leerColumna` con un CSV de prueba, incluyendo los cuatro casos de error.

---

## Checklist de fin de semana 1

Antes de pasar a la semana 2, tienes que poder responder **sin mirar**:

- [ ] ¿Cuándo `var` y cuándo `:=`?
- [ ] ¿Por qué `append` siempre se reasigna?
- [ ] ¿Qué diferencia hay entre un receptor por valor y uno por puntero?
- [ ] ¿Por qué `Producto` cumple una interfaz sin declararlo?
- [ ] ¿Qué añade `%w` y por qué `errors.Is` en vez de `==`?
- [ ] ¿Qué hace `defer` y por qué va pegado a la línea de abrir?
- [ ] ¿Por qué el orden de un `map` no se puede dar por bueno?
- [ ] ¿Por qué `Mediana` debe copiar el slice antes de ordenarlo?

Si alguna no sale, vuelve a ese día. **No se avanza con huecos**: la semana 2 los da todos por
sabidos.
