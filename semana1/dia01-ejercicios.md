# Día 1 · Ejercicios prácticos

Tres ejercicios con **solo lo del día 1**: variables, tipos, funciones y `Printf`.
Sin bucles ni listas — eso llega en los días 2 y 3. Van de menos a más.

Cada ejercicio va en **su propia carpeta** bajo `cmd\`. Una carpeta = un programa.

**Orden sugerido: A → B → C.** El A te enseña el error de tipos, el B te da soltura
con múltiples retornos, y el C te hace reproducir tú solo el patrón central del día.

---

## Cómo crear el archivo

### Con VS Code

```
code C:\Users\cg282\dev\aprendiendo-go
```

Clic derecho sobre `cmd` → **New Folder** → `dia01a`.
Clic derecho sobre `dia01a` → **New File** → `main.go`.

Acepta instalar la **extensión de Go** cuando te lo ofrezca: te subraya los errores
mientras escribes, en vez de que te los cuente el compilador después.

### Desde la terminal

En **cmd.exe**:

```cmd
mkdir cmd\dia01a
notepad cmd\dia01a\main.go
```

En **PowerShell**:

```powershell
New-Item -ItemType Directory cmd\dia01a
notepad cmd\dia01a\main.go
```

### El esqueleto que va dentro

Escribe esto primero y compruébalo antes de complicarte:

```go
package main

import "fmt"

func main() {
	fmt.Println("hola")
}
```

```powershell
go run ./cmd/dia01a
```

Si sale `hola`, la carpeta y el archivo están bien puestos. A partir de ahí solo es rellenar.

### Tres cosas que te pueden morder

1. **El nombre `main.go` es costumbre, no obligación.** Lo obligatorio es que dentro
   diga `package main` y tenga una `func main()`.
2. **Bloc de notas y la extensión.** Si en "Tipo" deja *Documentos de texto*, guardará
   `main.go.txt` y Go no lo verá. Elige "Todos los archivos" o pon el nombre entre comillas.
3. **Una carpeta bajo `cmd\` = un programa.** Dos `func main()` en la misma carpeta y Go
   se queja de que hay dos puntos de entrada.

---

# Ejercicio A · Ficha de producto

**Carpeta:** `cmd\dia01a\main.go`
**Practica:** declarar variables, los cuatro tipos básicos, verbos de `Printf`.

## Qué tienes que hacer

En `main`, declara cuatro variables:

| Variable | Valor | Tipo |
|---|---|---|
| `nombre` | `"Teclado mecánico"` | `string` (texto) |
| `precio` | `1250.5` | `float64` (decimal) |
| `unidades` | `3` | `int` (entero) |
| `disponible` | `true` | `bool` (sí/no) |

Declara **`nombre` con la forma larga** (`var nombre string = ...`) y las otras tres
**con la forma corta** (`:=`). Es a propósito: quiero que escribas las dos.

Luego calcula el total (`precio × unidades`) e imprímelo todo.

## Salida exacta que debes conseguir

```
Producto: Teclado mecánico
Precio unitario: 1250.50
Unidades: 3
Disponible: true
Total: 3751.50
El precio es de tipo float64 y las unidades de tipo int
```

## Pistas

- Verbos que necesitas: `%s` texto, `%.2f` decimal con dos cifras, `%d` entero,
  `%t` booleano (*true/false*), `%T` el nombre del tipo.
- **Aquí vas a chocar con algo.** `precio * unidades` **no compila**. Go dirá:

  ```
  invalid operation: precio * unidades (mismatched types float64 and int)
  ```

  Léelo: *"tipos que no encajan, float64 e int"*. Go **no mezcla tipos** ni siquiera
  entre números. Hay que convertir uno explícitamente:

  ```go
  total := precio * float64(unidades)
  ```

  `float64(unidades)` se lee "dame `unidades` visto como decimal". No cambia la variable
  original: crea un valor nuevo del otro tipo.
- Ese error es el objetivo real del ejercicio. **Provócalo primero**, léelo, y luego arréglalo.

---

# Ejercicio B · De segundos a reloj

**Carpeta:** `cmd\dia01b\main.go`
**Practica:** una función que devuelve **tres** valores, división entera y resto.

## Qué tienes que hacer

Escribe una función:

```go
func convertirTiempo(totalSegundos int) (int, int, int)
```

Recibe una cantidad de segundos y devuelve **horas, minutos y segundos** por separado.

En `main`, llámala con `9875` y con `59` e imprime los dos resultados.

## Salida exacta

```
9875 segundos son 2 h 44 min 35 s
59 segundos son 0 h 0 min 59 s
```

## Pistas

- Dos operadores nuevos, ambos sobre **enteros**:

  | Operador | Qué hace | Ejemplo |
  |---|---|---|
  | `/` | división **entera**: tira los decimales | `9875 / 3600` → `2` |
  | `%` | el **resto** de esa división | `9875 % 3600` → `2675` |

  > Cuidado: este `%` (operador resto, entre números) **no tiene nada que ver** con el
  > `%s` de `Printf` (hueco de plantilla). Se escriben igual y son cosas distintas.

- La receta: una hora son 3600 segundos, un minuto son 60.
  1. Horas = total dividido entre 3600.
  2. Lo que sobra = total resto 3600.
  3. Minutos = lo que sobra dividido entre 60.
  4. Segundos = lo que sobra resto 60.
- Recoger tres valores se hace igual que dos: `h, m, s := convertirTiempo(9875)`.
- Aquí `/` con decimales sería un **error**: no quieres "2.74 horas", quieres 2 horas
  y pico. La división entera es justo la herramienta correcta.

---

# Ejercicio C · Calculadora de IMC

**Carpeta:** `cmd\dia01c\main.go`
**Practica:** el patrón `(valor, ok)` del ejemplo de hoy, aplicado por ti.

## Qué tienes que hacer

Dos funciones:

```go
func calcularIMC(pesoKg, alturaM float64) (float64, bool)
func clasificarIMC(imc float64) string
```

**`calcularIMC`** devuelve el índice de masa corporal: `peso ÷ (altura × altura)`.
Si la altura es cero o negativa **no hay respuesta posible** — devuelve `0, false`, igual
que hace `dividir` en el ejemplo de hoy. Si todo está bien, devuelve el valor y `true`.

**`clasificarIMC`** devuelve un texto según el número:

| IMC | Texto |
|---|---|
| menor que 18.5 | `bajo peso` |
| menor que 25 | `normal` |
| menor que 30 | `sobrepeso` |
| 30 o más | `obesidad` |

En `main`, prueba **tres** casos: `(82.5, 1.78)`, `(60, 1.70)` y `(70, 0)`.
En el tercero tienes que detectar el `false` y avisar, **sin** imprimir un IMC falso.

## Salida exacta

```
Peso 82.5 kg, altura 1.78 m -> IMC 26.04 (sobrepeso)
Peso 60.0 kg, altura 1.70 m -> IMC 20.76 (normal)
Error: la altura debe ser mayor que cero
```

## Pistas

- La estructura de `calcularIMC` es **la misma** que la de `dividir` en el ejemplo:
  comprobar el caso imposible primero, salir con `return 0, false`, y si no, devolver
  el cálculo y `true`.
- Para encadenar las cuatro categorías se usa `if / else if / else`:

  ```go
  if imc < 18.5 {
      return "bajo peso"
  } else if imc < 25 {
      return "normal"
  }
  // ...
  ```

  El orden importa: se evalúa de arriba abajo y se queda en la **primera** que se cumple.
  Por eso `imc < 25` no necesita decir "y mayor que 18.5" — si lo fuera, ya habría salido antes.
- Una función puede devolver texto directamente: el tipo de retorno es `string` y
  devuelves `"normal"`.
- Al llamar: `imc, ok := calcularIMC(82.5, 1.78)` y luego `if !ok { ... }`.
- Para meter el texto de la categoría dentro del `Printf`, puedes llamar a la función ahí
  mismo: `..., imc, clasificarIMC(imc))`. Una llamada a función **es** un valor.
- `60` y `1.70` los escribes tal cual; como los parámetros son `float64`, Go los toma como
  decimales. Que salga `60.0` y `1.70` en pantalla es cosa de `%.1f` y `%.2f`.

---

## Ejecutar

```powershell
go run ./cmd/dia01a
go run ./cmd/dia01b
go run ./cmd/dia01c
```

Cuando tengas alguno —o cuando te atasques a mitad— pídeme que lo revise: qué falla,
**por qué**, y cómo se arregla. Si un error del compilador no lo entiendes, pégalo tal cual:
leer errores es la mitad del oficio.
