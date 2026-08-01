# Conceptos básicos · Empezando desde cero absoluto

> Este archivo no tiene ejercicios. Es solo para **entender de qué se está hablando**.
> Léelo despacio. Si algo no se entiende, no importa: vuelve después.

---

## Parte 1 · El mundo antes del código

### ¿Qué es un archivo?

Un archivo es **un papel guardado en el computador**. Tiene un nombre y un apellido:

```
main.go
 │    │
 │    └── el apellido (extensión): dice QUÉ TIPO de papel es
 └─────── el nombre: se lo pones tú
```

- `.txt` → un papel con texto normal
- `.md` → un papel con texto y formato (como este que estás leyendo)
- `.go` → **un papel con instrucciones para el computador, escritas en el idioma Go**

Un archivo `.go` es texto plano. Nada mágico. Lo podrías escribir en el Bloc de notas.

### ¿Qué es una carpeta?

Una caja donde metes archivos y otras cajas. Igual que en el escritorio de Windows.

Tu proyecto se ve así:

```
Aprendiendo GO\          ← caja grande (tu proyecto)
├── go.mod               ← papel de identificación del proyecto
├── cmd\                 ← caja de "programas ejecutables"
│   └── dia01\           ← caja de UN programa
│       └── main.go      ← el papel con las instrucciones
└── semana1\
    └── semana-1-sintaxis.md
```

Escribir una **ruta** es dar la dirección de un archivo:

```
cmd\dia01\main.go
```

Se lee: *"entra a la caja `cmd`, luego a la caja `dia01`, ahí está `main.go`"*.

> ⚠️ Windows usa `\` (barra invertida) para las carpetas.
> Go usa `/` (barra normal) cuando le hablas de paquetes.
> Por eso escribes `go run ./cmd/dia01` con barras normales. No es un error.

### ¿Qué es la terminal?

Es la ventana negra (PowerShell). Es **la misma computadora de siempre**, pero en vez de hacer doble clic, le escribes órdenes.

| Con el mouse | En la terminal |
|---|---|
| Abrir una carpeta | `cd nombre-carpeta` |
| Ver qué hay dentro | `ls` |
| Volver a la carpeta de atrás | `cd ..` |
| Saber dónde estoy parado | `pwd` |

Lo que ves al inicio de la línea te dice **dónde estás parado**:

```
PS C:\Users\cg282\OneDrive\Desktop\Aprendiendo GO>
   └──────────────── estás aquí ────────────────┘ └── escribes después de esto
```

🔑 **Esto es lo más importante de la terminal:** los comandos se ejecutan *desde donde estás parado*.
Si estás parado en el lugar equivocado, el comando falla aunque esté bien escrito.
Por eso siempre hacemos primero el `cd` a la carpeta del proyecto.

### ¿Qué es un programa?

Una **lista de órdenes**, en orden, de arriba hacia abajo. El computador las lee como tú lees una receta de cocina: primero la 1, después la 2, después la 3.

El computador es rapidísimo pero **bruto**: hace exactamente lo que dice el papel. Si el papel está mal, hace la burrada sin avisar. No adivina.

### ¿Qué es Go?

Go es **un idioma** para escribirle esas órdenes al computador. Hay muchos (Python, Java, JavaScript...). Go es de los más **estrictos y aburridos**, y eso es bueno para aprender: te grita cuando algo está mal en vez de dejarte seguir.

Se pronuncia "gou". También le dicen *Golang*. Es lo mismo.

### ¿Qué es "compilar"?

El computador **no entiende** Go. Entiende solo unos y ceros.

**Compilar** = traducir tu archivo `.go` a unos y ceros.

```
main.go  ──[compilar]──>  programa que la máquina sí entiende  ──>  se ejecuta
 (tú)                              (traducción)                      (resultado)
```

Cuando escribes `go run`, Go hace **las dos cosas de una**: traduce y ejecuta.

Si tu archivo tiene un error de escritura, la traducción **falla** y no se ejecuta nada. Eso es bueno: el error aparece antes, no en la cara del usuario.

---

## Parte 2 · Los comandos que vas a usar

Todos se escriben **parado en la carpeta del proyecto**.

| Comando | Qué hace, en cristiano |
|---|---|
| `go run ./cmd/dia01` | Traduce y ejecuta ese programa **ahora**. No deja archivos. Es el que más usarás. |
| `go build -o bin/x.exe ./cmd/dia01` | Traduce y **guarda** el resultado como un `.exe` que puedes darle a otra persona. |
| `go fmt ./...` | Acomoda tu código bonito (espacios, sangrías). En Go **no se discute de estilo**: manda la máquina. |
| `go vet ./...` | Revisa errores tontos que sí compilan pero están mal. Si no dice nada, todo bien. |
| `go test ./...` | Corre las pruebas automáticas (esto llega en el Día 6). |
| `go mod init nombre` | Crea el proyecto. **Ya lo hicimos. No se repite nunca.** |

> Los `./...` significan *"en esta carpeta y en todas las de adentro"*.

---

## Parte 3 · Las piezas de un archivo Go

Este es el programa que ya ejecutaste. Vamos pieza por pieza.

```go
package main

import "fmt"

func main() {
	fmt.Println("hola")+
}
```

### `package main` — la primera línea, siempre

Dice: **"este archivo es un programa que se puede ejecutar"**.

La palabra `main` significa *principal*. Es una palabra reservada, mágica.
Si escribieras otra cosa (`package stats`), sería una **caja de piezas** para que otros programas la usen, no algo ejecutable.

### `import "fmt"` — pedir herramientas prestadas

Go trae un montón de herramientas listas, pero **no te las da si no las pides**.

`import "fmt"` significa: *"préstame la caja `fmt`"*. Esa caja sirve para **escribir en pantalla**. Sin esa línea, tu programa no sabría imprimir.

Si necesitas varias, se escriben en grupo:

```go
import (
	"fmt"
	"os"
	"strings"
)
```

> 🚨 Go **no te deja** importar algo que no usas. El programa ni compila. Es a propósito: obliga a mantener las cosas limpias.

Cajas que verás pronto:

| Caja | Para qué sirve |
|---|---|
| `fmt` | escribir en pantalla y dar formato |
| `os` | hablar con el sistema: archivos, salir del programa |
| `strings` | manipular texto |
| `strconv` | convertir texto ↔ número |
| `errors` | manejar errores |
| `math` | matemáticas |

### `func main()` — la puerta de entrada

`func` viene de **función**: una máquina que hace algo.

`main` es la función mágica: **cuando ejecutas el programa, Go la busca y arranca ahí**. Siempre. Todo programa tiene exactamente una.

```go
func main() {
	// aquí adentro va lo que hace el programa
}
```

Las llaves `{` y `}` marcan **dónde empieza y dónde termina**. Todo lo que esté adentro le pertenece.

### Los comentarios — notas para humanos

```go
// Esto es un comentario. El computador lo IGNORA por completo.
```

Sirven para recordarte a ti mismo qué estabas pensando. Escribe muchos mientras aprendes.

---

## Parte 4 · Variables: cajas con nombre

Una **variable** es una caja donde guardas algo y le pones etiqueta.

```go
edad := 16
```

Se lee: *"crea una caja llamada `edad` y mete adentro el número 16"*.

Después, en cualquier parte, escribir `edad` es lo mismo que escribir `16`.

### Las dos formas de crear una caja

```go
var nombre string = "Go"   // forma larga: caja, tipo, valor
nombre := "Go"             // forma corta: Go adivina el tipo
```

El `:=` significa **"crea la caja Y métele esto"**. Es lo que más vas a usar.

Cuando la caja **ya existe** y solo quieres cambiar el contenido, se usa `=` a secas:

```go
edad := 16   // la creo
edad = 17    // le cambio el valor (ya existe, sin los dos puntos)
```

> ⚠️ `:=` solo funciona **dentro** de una función. Fuera hay que usar `var`.

### 🚨 El error #1 de todo principiante

| Símbolo | Significa | Ejemplo |
|---|---|---|
| `=` | **guardar** | `edad = 17` → mete 17 en la caja |
| `==` | **comparar** | `edad == 17` → pregunta: ¿es igual a 17? Responde sí o no |

No son lo mismo. Confundirlos es un clásico.

---

## Parte 5 · Los tipos: de qué está hecha cada cosa

Go es estricto: **cada caja guarda un solo tipo de cosa y no se puede mezclar**.

| Tipo | Qué guarda | Ejemplo |
|---|---|---|
| `string` | texto | `"hola"` · siempre entre comillas dobles `"` |
| `int` | número entero (sin coma) | `16`, `-5`, `0` |
| `float64` | número con decimales | `3.14`, `48.50` |
| `bool` | verdadero o falso, nada más | `true`, `false` |

### 🚨 La trampa de la división entera

```go
9 / 5      // da 1   ← ¡son enteros! Go TIRA los decimales a la basura
9.0 / 5.0  // da 1.8 ← con el .0 le dices "estos tienen decimales"
```

Go **no redondea**, corta. Esta trampa te va a morder alguna vez. Ya estás avisado.

### Mezclar tipos está prohibido

```go
edad := 16         // int
altura := 1.75      // float64
edad + altura       // ❌ ERROR: no compila
```

Otros lenguajes lo dejan pasar y luego dan resultados raros. Go prefiere gritarte de una.

### El valor cero

Si creas una caja sin meterle nada, Go **no la deja vacía**, le mete un valor por defecto:

| Tipo | Valor cero |
|---|---|
| `int` | `0` |
| `float64` | `0` |
| `string` | `""` (texto vacío) |
| `bool` | `false` |

Esto es una de las cosas buenas de Go: nunca hay cajas con basura adentro.

---

## Parte 6 · Imprimir en pantalla

### `Println` — lo simple

```go
fmt.Println("hola")        // imprime: hola
fmt.Println("edad:", 16)   // imprime: edad: 16
```

Imprime lo que le des y **salta de línea** solo (`ln` = *line*).

### `Printf` — con formato

```go
fmt.Printf("Tengo %d años\n", 16)
```

Sale: `Tengo 16 años`

Los `%algo` son **huecos**. Go los va rellenando en orden con lo que pongas después de la coma.

| Hueco | Rellena con |
|---|---|
| `%s` | texto (**s**tring) |
| `%d` | número entero (**d**ígito) |
| `%f` | decimal · `%.2f` = con 2 decimales |
| `%v` | lo que sea (el comodín) |
| `%T` | te dice **de qué tipo** es (para investigar) |
| `%%` | un símbolo `%` de verdad |

Y `\n` significa **"salta de línea"**. `Printf` **no** salta solo — si se te olvida, todo sale pegado.

```go
fmt.Printf("%s tiene %d años y mide %.2f\n", "Ana", 30, 1.6543)
//          └─texto  └─entero      └─decimal con 2 cifras
// Sale: Ana tiene 30 años y mide 1.65
```

---

## Parte 7 · Funciones: máquinas que hacen algo

```go
func sumar(a int, b int) int {
	return a + b
}
```

Desarmémosla:

```
func  sumar  (a int, b int)  int  {  return a + b  }
 │      │          │          │             │
 │      │          │          │             └── qué devuelve
 │      │          │          └── el TIPO de lo que sale
 │      │          └── lo que ENTRA (con su tipo cada uno)
 │      └── el nombre que le pones
 └── "voy a crear una función"
```

Para **usarla** (se dice *llamarla*):

```go
resultado := sumar(3, 4)   // resultado ahora vale 7
```

### 🔑 Lo más particular de Go: devolver VARIAS cosas

En casi todos los lenguajes una función devuelve **una** cosa. En Go puede devolver varias:

```go
func dividir(a, b float64) (float64, bool) {
	if b == 0 {
		return 0, false      // "no pude" 
	}
	return a / b, true       // "aquí está, y salió bien"
}
```

Y al usarla recoges **las dos** en dos cajas:

```go
resultado, ok := dividir(10, 3)
```

**¿Para qué sirve esto?** Para responder dos preguntas a la vez: *¿cuál es el resultado?* y *¿salió bien?*

Esta es **la base de todo Go**. Más adelante, en vez de `bool`, la segunda cosa será un `error` que además explica *qué* falló. Es el patrón que vas a ver mil veces:

```go
valor, err := hacerAlgo()
if err != nil {
	// algo salió mal, manéjalo
}
```

### El basurero `_`

Si la función devuelve dos cosas y solo te interesa una:

```go
resultado, _ := dividir(10, 3)   // el _ tira la segunda a la basura
```

🔑 **¿Por qué existe esto?** Porque Go **prohíbe crear una variable y no usarla** — el programa ni compila. El `_ ` es la forma oficial de decir *"esto no lo necesito, no me regañes"*.

---

## Parte 8 · Vocabulario de rescate

Cuando alguien (o yo) diga estas palabras, esto es lo que significan:

| Palabra | En cristiano |
|---|---|
| **código / código fuente** | el texto que tú escribes en el archivo |
| **compilar** | traducir tu texto a algo que la máquina entiende |
| **ejecutar / correr** | ponerlo a funcionar |
| **variable** | una caja con nombre que guarda algo |
| **tipo** | de qué está hecha la caja (texto, número...) |
| **función** | una máquina: le metes cosas, te devuelve cosas |
| **llamar una función** | usarla, ponerla a trabajar |
| **parámetro / argumento** | lo que le metes a la función |
| **retornar / devolver** | lo que la función te entrega de vuelta |
| **paquete** | una caja de herramientas (`fmt`, `os`...) |
| **importar** | pedir prestada una caja de herramientas |
| **módulo** | tu proyecto completo (lo define `go.mod`) |
| **bug** | un error de lógica: compila pero hace algo mal |
| **error de compilación** | está mal escrito, ni siquiera arranca |
| **stack trace / panic** | el programa explotó y te muestra dónde |
| **idiomático** | "la forma en que los que saben Go lo escriben" |
| **CLI** | programa que se usa desde la terminal, sin ventanas |
| **debuggear** | buscar dónde está el error |

---

## Parte 9 · Qué hacer cuando algo falla

**Esto va a pasar todo el tiempo. Es normal, no es que seas malo.** Hasta el que lleva 20 años programando ve errores cada media hora.

### Lee el error de abajo hacia arriba

Go te dice **el archivo, la línea y qué pasa**:

```
./main.go:12:5: undefined: fmt.Printl
   │      │  │              │
   │      │  │              └── el problema: eso no existe
   │      │  └── columna 5
   │      └── LÍNEA 12  ← anda ahí
   └── en este archivo
```

Ese ejemplo es un `Println` mal escrito (le falta la `n`).

### Los 5 errores que más vas a ver

| Mensaje | Qué significa | Cómo se arregla |
|---|---|---|
| `undefined: X` | escribiste un nombre que no existe | revisa mayúsculas y ortografía |
| `declared and not used` | creaste una caja y nunca la usaste | úsala o bórrala |
| `imported and not used` | pediste una herramienta y no la usaste | borra ese `import` |
| `cannot use X (type A) as type B` | mezclaste tipos | convierte, o revisa los `.0` |
| `syntax error` | falta una llave `}`, un paréntesis o una coma | revisa la línea anterior también |

### Reglas de oro para no volverte loco

1. **Corre `go run` seguido.** Cada vez que cambies 2 o 3 líneas. Así el error es fácil de encontrar. Si escribes 100 líneas y luego ejecutas, no vas a saber dónde está el problema.
2. **Un error a la vez.** Si salen 5, arregla el primero y vuelve a correr. Muchas veces los otros 4 eran consecuencia del primero.
3. **Mayúsculas y minúsculas importan.** `Nombre` y `nombre` son cajas **distintas** para Go.
4. **Cada `{` necesita su `}`.** Cuenta si algo se ve raro.
5. **Copiar el error y buscarlo en Google funciona.** En serio. Todo el mundo lo hace.

---

## Checklist · ¿Ya puedo seguir?

No necesitas saber escribir nada todavía. Solo **reconocer**:

- [ ] Sé abrir PowerShell y pararme en mi carpeta con `cd`
- [ ] Sé que `go run ./cmd/algo` ejecuta un programa
- [ ] Reconozco `package main`, `import` y `func main()` cuando los veo
- [ ] Entiendo que una variable es una caja con nombre
- [ ] Sé la diferencia entre `=` (guardar) y `==` (comparar)
- [ ] Sé qué son `string`, `int`, `float64` y `bool`
- [ ] Sé que `%s`, `%d` y `%.2f` son huecos que se rellenan
- [ ] Entiendo que una función recibe cosas y devuelve cosas
- [ ] Sé que en Go una función puede devolver **varias** cosas
- [ ] Sé que si algo falla, el error me dice el archivo y la línea

Si marcaste la mitad, vas bien. Lo demás se acomoda con la práctica.
