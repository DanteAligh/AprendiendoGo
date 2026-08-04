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

# Ejercicios — Semana 1: Fundamentos de Go

Este documento es tu profesor. Para cada día vas a encontrar una explicación
del concepto (qué es, por qué existe en Go, cuándo se usa) y luego una lista
de ejercicios de dificultad progresiva. **No hay soluciones aquí a propósito.**
La idea es que pienses, pruebes, falles, y vuelvas a intentar — así es como se
aprende un lenguaje de verdad.

Antes de cada ejercicio, recuerda activar Go en tu terminal:

```powershell
go run tu_archivo.go
```

---

## Día 1: Instalación, estructura de un programa, variables y tipos

### Teoría

Todo programa Go ejecutable empieza igual: un archivo que declara
`package main` y una función `func main()`. El paquete `main` le dice al
compilador "este no es una librería, es un programa que se puede ejecutar", y
`func main()` es el punto de entrada — literalmente lo primero que corre.
Esto es distinto a lenguajes como Python o JavaScript donde el código "suelto"
en el archivo se ejecuta de arriba hacia abajo sin necesidad de un punto de
entrada explícito. Go es más parecido a C o Java en este sentido: todo vive
dentro de funciones, y el runtime busca específicamente `main()`.

Para las variables, Go te da dos caminos:

- `var nombre tipo = valor` — la forma explícita, útil cuando quieres dejar
  claro el tipo o cuando declaras la variable sin inicializarla todavía
  (`var contador int` la crea con el "valor cero" de ese tipo: `0` para
  números, `""` para strings, `false` para bools).
- `nombre := valor` — el operador de declaración corta, que **infiere** el
  tipo a partir del valor. Es la forma más usada dentro de funciones porque es
  más corta y Go es lo suficientemente inteligente para deducir el tipo. Ojo:
  `:=` solo funciona dentro de funciones, no a nivel de paquete.

Go es un lenguaje de **tipado estático y fuerte**: una vez que una variable es
`int`, no puedes asignarle un `string` después, y ni siquiera puedes sumar un
`int` con un `float64` directamente sin convertir uno de los dos. Esto es
diferente a lenguajes como Python o JavaScript donde el tipo "fluye" más
libremente. Al principio se siente rígido, pero en un proyecto grande (como el
ERP que es tu meta) esta rigidez es la que evita que un bug de tipos se cuele
hasta producción.

Los tipos básicos que vas a usar constantemente: `int` (enteros), `float64`
(decimales), `string` (texto), `bool` (verdadero/falso). Existen variantes
como `int8`, `int32`, `int64`, `float32`, etc., pero por ahora quédate con los
cuatro básicos.

Las **constantes** (`const`) se declaran como variables pero su valor no puede
cambiar después de declarado — el compilador te lo va a impedir. Se usan para
valores que conceptualmente nunca cambian durante la ejecución (un IVA fijo,
un límite máximo, un nombre de configuración).

Finalmente, como Go no convierte tipos automáticamente, necesitas el paquete
`strconv` para convertir, por ejemplo, un `string` leído del teclado a un
`int` (`strconv.Atoi`) o un `int` a `string` (`strconv.Itoa`). Esto lo vas a
usar todo el tiempo cuando leas entrada del usuario, que es justo el tema del
día 2.

### Ejercicios

1. Crea un programa que declare tu nombre (`string`), tu edad (`int`) y si
   estás actualmente estudiando (`bool`) como variables usando `:=`, y que
   los imprima todos en una sola línea con `fmt.Println`.
   Ejemplo de salida esperada: `Christian 30 true`

2. Declara una constante `PI` con el valor `3.1416` y una variable `radio`
   de tipo `float64`. Calcula el área de un círculo (área = PI * radio²) y
   guarda el resultado en otra variable antes de imprimirlo.
   Entrada: `radio = 4`. Salida esperada aproximada: `50.2656`

3. Declara una variable de tipo `string` que contenga un número escrito como
   texto, por ejemplo `"42"`. Usa `strconv.Atoi` para convertirla a `int`,
   súmale 8, y muestra el resultado. (Pista: `strconv.Atoi` devuelve **dos**
   valores — el número y un posible error. Por ahora puedes ignorar el error
   guardándolo en una variable, ya lo vamos a manejar bien en semanas
   futuras).
   Entrada: `"42"`. Salida esperada: `50`

4. Declara tres variables sin inicializar usando `var` (una `int`, una
   `string`, una `bool`) e imprime sus "valores cero" usando `fmt.Println`
   para comprobar qué valor les asigna Go por defecto.
   Salida esperada: `0`, `` (cadena vacía, se ve como nada), `false`

5. (Un poco más retador) Declara una variable `precioDolares` de tipo
   `float64`, una constante `TASA_CAMBIO` (por ejemplo 4000, simulando pesos
   colombianos por dólar) y calcula cuánto sería el precio en la moneda
   local. Luego convierte ese resultado a `string` usando el paquete
   `strconv` (busca la función que convierte `float64` a `string`, no es
   `Itoa`) y concaténalo con un mensaje de texto para imprimir todo junto en
   una sola línea con `fmt.Println`.
   Entrada: `precioDolares = 19.99`. Salida esperada aproximada:
   `El precio en pesos es: 79960.00...`

💡 **Reto extra:** Investiga qué pasa si intentas sumar directamente un `int`
y un `float64` sin convertir ninguno (`var a int = 5; var b float64 = 2.5;
a + b`). Escribe ese código a propósito, mira el error que te da el
compilador, y trata de explicarte a ti mismo por qué Go es tan estricto con
esto.

---

## Día 2: Operadores, entrada por consola y formateo de salida

### Teoría

Go tiene los operadores que ya conoces de otros lenguajes: aritméticos
(`+ - * / %`), de comparación (`== != < > <= >=`) y lógicos (`&& || !`). Nada
exótico aquí, salvo un detalle: la división entre dos enteros (`int / int`)
en Go **trunca** el resultado (es división entera), igual que en muchos
lenguajes de bajo nivel. Si quieres una división con decimales, al menos uno
de los operandos debe ser `float64`.

Para leer datos desde la consola tienes dos herramientas principales:

- `fmt.Scan(&variable)` / `fmt.Scanln(&variable)` — simples, leen "hasta el
  siguiente espacio en blanco". Nota el `&` antes de la variable: eso es el
  operador de **dirección de memoria** (puntero). `fmt.Scan` necesita saber
  *dónde* en memoria escribir el valor que lee, no solo el valor actual, por
  eso le pasas la dirección. No te preocupes por entender punteros a fondo
  todavía — eso llega en semanas posteriores — por ahora solo recuerda: "para
  leer con Scan, se me olvida el `&` y nada funciona".
- `bufio.NewScanner(os.Stdin)` — más robusto, lee **líneas completas**
  incluyendo espacios (útil si quieres leer un nombre completo como
  "Christian Peña"). Se usa creando un scanner, llamando `scanner.Scan()`
  para avanzar una línea, y `scanner.Text()` para obtener el texto leído.

Para mostrar salida con formato, Go usa `fmt.Printf` con **verbos** (los
`%algo`):

- `%d` — enteros
- `%f` — decimales (puedes controlar precisión con `%.2f` para 2 decimales)
- `%s` — strings
- `%v` — el valor "por defecto", funciona con casi cualquier tipo (útil
  cuando no quieres pensar en el verbo exacto)
- `%T` — el **tipo** de la variable (muy útil para depurar y entender qué te
  está devolviendo una función)
- `%t` — booleanos

Esto es distinto a, por ejemplo, el f-string de Python o los template
literals de JavaScript, pero cumple el mismo propósito: construir strings de
salida con variables incrustadas de forma controlada. En backend, vas a usar
`Printf`/`Sprintf` constantemente para logs y mensajes de depuración.

### Ejercicios

1. Pide al usuario dos números enteros por consola (usando `fmt.Scan`),
   súmalos, réstalos, multiplícalos y divídelos, mostrando los 4 resultados
   con `fmt.Printf` usando el verbo `%d`.
   Entrada: `10` y `3`. Salida esperada:
   `Suma: 13, Resta: 7, Multiplicación: 30, División: 3`

2. Repite el ejercicio anterior pero ahora también calcula el **residuo**
   (módulo) de la división. (Pista: necesitarás el operador `%`).
   Entrada: `10` y `3`. Salida esperada adicional: `Residuo: 1`

3. Pide al usuario su nombre completo usando `bufio.NewScanner` (debe
   soportar espacios) y su edad usando `fmt.Scan`. Imprime un mensaje usando
   `%s` y `%d` que diga algo como "Hola, [nombre], tienes [edad] años".
   Entrada: `Christian Peña` y `30`. Salida esperada:
   `Hola, Christian Peña, tienes 30 años`

4. Declara dos variables numéricas, una `int` y una `float64`, e imprime
   ambas junto con su tipo usando el verbo `%T` para comprobar qué reporta
   Go. Luego escribe una expresión lógica que combine `&&` y `||` con dos
   condiciones de comparación sobre esas variables (por ejemplo: "el entero
   es mayor que 5 Y el decimal es menor que 10") e imprime el resultado
   booleano con `%t`.

5. (Un poco más retador) Pide al usuario un precio (float64) y una cantidad
   de artículos (int). Calcula el total, y usa `fmt.Printf` con `%.2f` para
   mostrar el total con exactamente 2 decimales, sin importar cuántos
   decimales tenga el precio original.
   Entrada: `precio = 9.999`, `cantidad = 3`. Salida esperada aproximada:
   `Total: 29.997` truncado/redondeado a `29.99` o `30.00` según cómo lo
   calcules — experimenta y observa la diferencia.

💡 **Reto extra:** Investiga la diferencia entre `fmt.Scan` y `fmt.Scanln`
leyendo múltiples valores separados por espacios en la misma línea contra
valores en líneas distintas. ¿Cambia el comportamiento?

---

## Día 3: Condicionales — if/else y switch

### Teoría

El `if` de Go se ve casi igual al de C, Java o JavaScript, con una diferencia
notable: **no lleva paréntesis alrededor de la condición**, y las llaves `{}`
son **obligatorias** incluso para una sola línea. Esto es una decisión de
diseño deliberada: Go prefiere quitarte opciones de "estilo" para que todo el
código Go del mundo se vea parecido, algo muy valioso cuando trabajas en
equipo en un proyecto grande.

```go
if edad >= 18 {
    // ...
} else if edad >= 13 {
    // ...
} else {
    // ...
}
```

Una característica que te va a encantar: el `if` puede llevar una
**declaración corta previa**, separada por `;`, cuyo alcance (scope) es solo
ese bloque `if`/`else`. Es muy común ver esto en Go cuando una función
devuelve un valor y un error a la vez:

```go
if valor, err := algo(); err == nil {
    // valor solo existe aquí dentro
}
```

Sobre el `switch`: en Go **no necesitas `break`** al final de cada `case`
(a diferencia de C/Java/JS) — Go no hace "fall-through" por defecto, cada
case termina automáticamente al ejecutarse. Si de verdad quieres que un case
continúe al siguiente, existe la palabra `fallthrough`, pero rara vez se usa.

Lo más interesante es el **switch sin expresión** (a veces llamado
"switch verdadero"): en vez de comparar una variable contra varios valores,
cada `case` lleva su propia condición booleana, funcionando como una cadena
de `if/else if` pero más legible:

```go
switch {
case edad < 13:
    // ...
case edad < 18:
    // ...
default:
    // ...
}
```

Por último: **Go no tiene operador ternario** (`condición ? a : b`). Esto es
intencional — los diseñadores del lenguaje consideraron que el ternario hace
el código más difícil de leer a simple vista, especialmente cuando se anida.
La forma idiomática de resolver "elegir entre dos valores según una
condición" en Go es simplemente un `if/else` normal que asigna a una
variable, o encapsular esa lógica en una función pequeña. Al principio se
siente como escribir "más líneas", pero con el tiempo vas a notar que el
código es más fácil de seguir.

### Ejercicios

1. Pide un número entero al usuario e imprime si es par o impar usando
   `if/else`. (Pista: el operador `%` de nuevo).
   Entrada: `7`. Salida esperada: `Es impar`

2. Pide tres números enteros e imprime cuál es el mayor de los tres, usando
   `if/else if/else`.
   Entrada: `4`, `9`, `2`. Salida esperada: `El mayor es 9`

3. Pide una calificación numérica (0 a 100) y usa un **switch sin
   expresión** para imprimir la letra correspondiente: 90+ es "A", 80-89 es
   "B", 70-79 es "C", menor a 70 es "F".
   Entrada: `85`. Salida esperada: `B`

4. Pide un número del 1 al 7 y usa un `switch` normal (con expresión) para
   imprimir el nombre del día de la semana correspondiente (1 = Lunes, etc).
   Si el número no está entre 1 y 7, usa el `default` para imprimir un
   mensaje de error.
   Entrada: `3`. Salida esperada: `Miércoles`

5. (Un poco más retador) Sin usar operador ternario (porque Go no lo tiene),
   escribe un programa que reciba dos números y una letra que representa una
   operación (`+`, `-`, `*`, `/`), y usando `switch` calcule e imprima el
   resultado de aplicar esa operación a los dos números. Si la letra no es
   ninguna de esas cuatro, imprime un mensaje de "operación no válida".
   Entrada: `10`, `4`, `/`. Salida esperada: `2.5`

💡 **Reto extra:** Reescribe el ejercicio 1 usando la sintaxis de `if` con
declaración corta previa (`if resto := numero % 2; resto == 0 { ... }`) para
practicar ese patrón que vas a ver mucho en Go real.

---

## Día 4: Ciclos — for, range, break/continue, labels

### Teoría

Go tiene una sola palabra clave para ciclos: `for`. No existen `while` ni
`do-while` como palabras separadas — en Go, **`for` puede tomar tres formas**
distintas, y eso cubre todos los casos:

**Forma clásica** (como el `for` de C/Java): inicialización, condición,
actualización, separadas por `;`.

```go
for i := 0; i < 5; i++ {
    // ...
}
```

**Forma "while"**: solo la condición, sin inicialización ni actualización
explícita en la misma línea de `for`.

```go
for condicion {
    // ...
}
```

**Forma infinita**: sin nada, se detiene solo con un `break` interno.

```go
for {
    // ...
    if algo {
        break
    }
}
```

Esta unificación es una de las cosas que hace que Go se sienta minimalista:
en vez de tres palabras clave distintas para tres variantes de repetición, hay
una sola construcción flexible.

`range` es la forma idiomática de recorrer colecciones (arrays, slices, maps,
strings, channels). En vez de manejar un índice manualmente, `range` te da el
índice y el valor en cada iteración:

```go
for indice, valor := range coleccion {
    // ...
}
```

Si no necesitas el índice, se descarta con `_` (el "blank identifier" de Go,
que verás todo el tiempo): `for _, valor := range coleccion`.

`break` sale completamente del ciclo más cercano; `continue` salta a la
siguiente iteración sin ejecutar el resto del cuerpo. Igual que en la mayoría
de lenguajes.

Donde Go se pone interesante es con **ciclos anidados y labels**: si tienes un
`for` dentro de otro `for` y quieres hacer `break` o `continue` sobre el
ciclo **externo** (no el interno, que es el comportamiento por defecto),
puedes etiquetar el ciclo externo con un nombre seguido de `:` y referenciarlo:

```go
externo:
for i := 0; i < 5; i++ {
    for j := 0; j < 5; j++ {
        if j == 2 {
            continue externo // salta a la siguiente iteración de i, no de j
        }
    }
}
```

Esto es poco común en otros lenguajes de alto nivel (Python ni siquiera lo
soporta) pero en Go es una herramienta legítima y muy usada cuando procesas
estructuras de datos anidadas — algo que vas a hacer constantemente al
recorrer resultados de bases de datos o respuestas de APIs en tu futuro ERP.

### Ejercicios

1. Usa un `for` clásico para imprimir los números del 1 al 20.
   Salida esperada: `1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20`
   (o cada uno en su propia línea, como prefieras).

2. Usa un `for` en forma "while" para calcular la suma de todos los números
   del 1 al 100 (sin usar la fórmula matemática directa, hazlo sumando en
   cada iteración).
   Salida esperada: `5050`

3. Usa `continue` dentro de un `for` para imprimir únicamente los números
   pares entre 1 y 30.
   Salida esperada: `2 4 6 8 10 12 14 16 18 20 22 24 26 28 30`

4. Escribe un ciclo infinito (`for {}`) que le pida al usuario un número
   repetidamente, y se detenga con `break` únicamente cuando el usuario
   ingrese `0`. Mientras tanto, ve acumulando la suma de todos los números
   ingresados (sin contar el 0 final) y muéstrala al terminar.
   Entrada: `5`, `3`, `2`, `0`. Salida esperada: `Suma total: 10`

5. (Un poco más retador) Usa dos ciclos `for` anidados con `range` para
   imprimir una tabla de multiplicar del 1 al 5 (5 filas, 5 columnas), y usa
   un `label` con `continue` en el ciclo externo para **saltar por completo
   la fila del número 3** (no imprimas nada de esa fila, pero sí continúa
   con la 4 y la 5).
   Salida esperada (formato libre, pero la fila del 3 no debe aparecer):
   ```
   1x1=1 1x2=2 1x3=3 1x4=4 1x5=5
   2x1=2 2x2=4 2x3=6 2x4=8 2x5=10
   4x1=4 4x2=8 4x3=12 4x4=16 4x5=20
   5x1=5 5x2=10 5x3=15 5x4=20 5x5=25
   ```

💡 **Reto extra:** Escribe un ciclo que recorra un string carácter por
carácter usando `range` (Go te va a dar el índice de byte y el "rune" —
carácter Unicode — en cada iteración) y cuente cuántas veces aparece una letra
específica, por ejemplo la 'a'.

---

## Día 5: Funciones — parámetros, retornos, named returns, variádicas, closures

### Teoría

Declarar una función en Go se ve así:

```go
func nombre(parametro1 tipo1, parametro2 tipo2) tipoDeRetorno {
    return valor
}
```

Fíjate en dos detalles frente a otros lenguajes: el tipo va **después** del
nombre del parámetro (no antes, como en C/Java), y si dos parámetros
seguidos tienen el mismo tipo puedes escribir el tipo una sola vez:
`func sumar(a, b int) int`.

Lo que hace a Go distinto (y muy usado en el mundo real) es que las funciones
pueden devolver **múltiples valores** separados por comas, algo que la
mayoría de lenguajes clásicos no soportan de forma nativa:

```go
func dividir(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("no se puede dividir entre cero")
    }
    return a / b, nil
}
```

Este patrón — devolver "el resultado" y "un error" juntos — es EL patrón
idiomático de manejo de errores en Go, y lo vas a ver en el 90% de las
funciones de la librería estándar. Todavía no vamos a profundizar en `error`
como tipo (eso viene en semanas posteriores), pero ya lo estás viendo aparecer
aquí y en `strconv.Atoi` del día 1.

Los **retornos con nombre** (named returns) te dejan nombrar los valores de
retorno en la firma de la función, lo que los "pre-declara" como variables
dentro de la función, y permite un `return` "desnudo" (sin argumentos) que
devuelve lo que esas variables tengan en ese momento:

```go
func dividir(a, b int) (resultado int, err error) {
    if b == 0 {
        err = errors.New("no se puede dividir entre cero")
        return // devuelve resultado (0) y err
    }
    resultado = a / b
    return
}
```

Es opcional usarlos, pero mejoran la legibilidad cuando el propósito de cada
valor de retorno no es obvio solo con el tipo.

Las **funciones variádicas** aceptan un número variable de argumentos del
mismo tipo, usando `...tipo`. Dentro de la función, ese parámetro se comporta
como un slice (tema del día 6):

```go
func sumarTodos(numeros ...int) int {
    total := 0
    for _, n := range numeros {
        total += n
    }
    return total
}
// se puede llamar: sumarTodos(1, 2, 3) o sumarTodos(1, 2, 3, 4, 5, 6)
```

`fmt.Println` y `fmt.Printf` son en realidad funciones variádicas — por eso
puedes pasarles cualquier cantidad de argumentos.

Finalmente, Go trata las funciones como **valores de primera clase**: puedes
guardarlas en variables, pasarlas como argumento a otra función, y
retornarlas desde otra función. Cuando una función es retornada desde otra y
"recuerda" variables del entorno donde fue creada, eso es un **closure**:

```go
func creadorDeContador() func() int {
    contador := 0
    return func() int {
        contador++
        return contador
    }
}
// cada llamada a la función devuelta incrementa Y recuerda su propio contador
```

Los closures son la base de patrones muy usados en Go para configuración,
middlewares en servidores HTTP, y generadores de funciones reutilizables —
algo que vas a usar bastante cuando construyas tu backend.

### Ejercicios

1. Escribe una función `esMayorDeEdad(edad int) bool` que devuelva `true` si
   la edad es mayor o igual a 18, y llámala desde `main` con al menos dos
   valores distintos, imprimiendo el resultado.
   Entrada: `20` y `15`. Salida esperada: `true` y `false`

2. Escribe una función que reciba dos números y un operador (string, uno de
   `"+"`, `"-"`, `"*"`, `"/"`) y devuelva **dos valores**: el resultado
   (float64) y un booleano que indique si la operación fue válida (por
   ejemplo, `false` si el operador no es reconocido o si se intenta dividir
   entre cero).
   Entrada: `10, 0, "/"`. Salida esperada: `0 false`

3. Escribe una función variádica `promedio(numeros ...float64) float64` que
   calcule el promedio de todos los números que reciba, sin importar
   cuántos sean.
   Entrada: `4, 8, 15, 16`. Salida esperada: `10.75`

4. Escribe una función que use **named returns** para dividir dos enteros,
   devolviendo el cociente y el residuo (por ejemplo `dividirConResto(17, 5)`
   debe devolver `3` y `2`). Usa `return` desnudo al final.
   Entrada: `17, 5`. Salida esperada: `Cociente: 3, Residuo: 2`

5. (Un poco más retador) Escribe una función `creadorDeMultiplicador(factor
   int) func(int) int` que devuelva una función (closure) que multiplique
   cualquier número por el `factor` que se le pasó al crearla. Crea dos
   multiplicadores distintos (por ejemplo, uno "duplicador" con factor 2 y
   uno "triplicador" con factor 3) y demuestra que cada uno "recuerda" su
   propio factor de forma independiente.
   Entrada: `duplicador(5)`, `triplicador(5)`. Salida esperada: `10` y `15`

💡 **Reto extra:** Escribe una función de orden superior `aplicar(numeros
[]int, operacion func(int) int) []int` que reciba un slice de enteros y una
función, y devuelva un nuevo slice con la función aplicada a cada elemento
(esto es básicamente un "map" hecho a mano — más adelante Go 1.21+ trae
utilidades similares en el paquete `slices`, pero hacerlo a mano primero te
enseña cómo funciona por dentro).

---

## Día 6: Arrays vs Slices

### Teoría

Un **array** en Go tiene tamaño **fijo**, definido en su tipo:
`var numeros [5]int` es un array de exactamente 5 enteros, y ese 5 es parte
del tipo — un `[5]int` y un `[3]int` son tipos distintos y no puedes asignar
uno al otro. Esto es bastante rígido, y en la práctica **casi nadie usa
arrays directamente en Go**. Existen principalmente como la base sobre la que
se construyen los slices.

Un **slice** es una "vista" flexible y redimensionable sobre un array
subyacente (que Go maneja por ti, sin que tengas que pensarlo). Se declara
así:

```go
numeros := []int{1, 2, 3}       // slice literal, sin tamaño en los corchetes
otro := make([]int, 5)          // slice de 5 enteros, todos en 0
```

Fíjate en la diferencia sintáctica: el array lleva un número dentro de
`[ ]` (`[5]int`), el slice lo deja vacío (`[]int`). Es una diferencia sutil
pero crítica de leer bien.

Los slices tienen dos propiedades clave:

- `len(slice)` — cuántos elementos tiene **actualmente**.
- `cap(slice)` — cuántos elementos puede tener **antes de que Go necesite
  reservar más memoria internamente** (la capacidad del array subyacente).

`append(slice, valor)` agrega un elemento al final. Si la capacidad no
alcanza, Go crea automáticamente un array subyacente más grande, copia los
datos, y te devuelve un slice apuntando al nuevo array — **por eso siempre
debes reasignar el resultado**: `numeros = append(numeros, 4)`, nunca ignorar
el valor de retorno.

El **slicing** con `[a:b]` te da un sub-slice desde el índice `a` (incluido)
hasta `b` (excluido) del slice original: `numeros[1:3]` te da los elementos
en posición 1 y 2, no la 3. Ojo con esto: comparte memoria con el slice
original, así que modificar el sub-slice puede modificar el original — es
algo que sorprende a mucha gente que viene de otros lenguajes.

`copy(destino, origen)` copia elementos de un slice a otro **sin** compartir
memoria, a diferencia de simplemente reasignar o hacer slicing.

¿Por qué slices y no arrays? Porque en programación real casi nunca sabes de
antemano cuántos elementos vas a tener (¿cuántos usuarios va a tener tu ERP?
¿cuántas filas va a devolver una consulta a la base de datos?). Los slices te
dan esa flexibilidad, y son tan centrales en Go que la mayoría de funciones de
la librería estándar los usan en vez de arrays. Piensa en los slices como el
equivalente de las listas dinámicas de Python o los `ArrayList` de Java, pero
con una relación muy explícita y visible con la memoria que manejan por
debajo — otra muestra de la filosofía de Go de no esconderte lo que pasa por
debajo.

### Ejercicios

1. Declara un slice de 5 nombres (strings) usando slice literal, e imprime
   su longitud con `len()`.
   Entrada: `["Ana", "Luis", "Marta", "Pedro", "Sofía"]`. Salida esperada:
   `5`

2. Parte de un slice vacío de enteros (`[]int{}`) y usa `append` en un ciclo
   `for` para llenarlo con los números del 1 al 10. Imprime el slice
   completo al final.
   Salida esperada: `[1 2 3 4 5 6 7 8 9 10]`

3. Dado un slice de 10 números, usa slicing (`[a:b]`) para extraer solamente
   los primeros 3 y los últimos 3 elementos, imprimiendo ambos sub-slices por
   separado.
   Entrada: `[10 20 30 40 50 60 70 80 90 100]`. Salida esperada:
   `[10 20 30]` y `[80 90 100]`

4. Crea un slice de 5 enteros con `make`, imprime su `len` y `cap`. Luego
   usa `append` para agregar 3 elementos más y vuelve a imprimir `len` y
   `cap`, observando cómo cambia (o no) la capacidad. (Pista: la forma en que
   Go crece la capacidad no es "de 1 en 1", investiga y observa el patrón).

5. (Un poco más retador) Escribe una función `filtrarPares(numeros []int)
   []int` que reciba un slice de enteros y devuelva un **nuevo** slice que
   contenga solo los números pares, usando `append` sobre un slice vacío
   dentro de la función (no modifiques el slice original).
   Entrada: `[1 2 3 4 5 6 7 8 9 10]`. Salida esperada: `[2 4 6 8 10]`

💡 **Reto extra:** Declara un slice, saca un sub-slice de él con `[a:b]`, y
modifica un elemento del sub-slice. Imprime el slice original después del
cambio y comprueba si también cambió (esto demuestra que comparten memoria).
Luego repite el experimento pero usando `copy` en vez de slicing directo, y
compara el resultado.

---

## Día 7: Maps + repaso integrador de la semana

### Teoría

Un **map** en Go es una colección de pares clave-valor, el equivalente a un
diccionario en Python o un objeto/`Map` en JavaScript. Se declara así:

```go
edades := make(map[string]int)          // map vacío, listo para usar
precios := map[string]float64{           // map literal con valores iniciales
    "manzana": 2500,
    "pera":    3000,
}
```

El tipo `map[TipoClave]TipoValor` es explícito sobre qué tipo de clave y qué
tipo de valor maneja — no puedes mezclar claves de distintos tipos como en
JavaScript.

Las operaciones CRUD (crear, leer, actualizar, eliminar) son directas:

```go
edades["Christian"] = 30        // crear o actualizar
edad := edades["Christian"]     // leer
delete(edades, "Christian")     // eliminar
```

Iterar un map con `range` te da clave y valor en cada vuelta, **pero el orden
no está garantizado** — cada ejecución puede recorrer las llaves en un orden
distinto. Esto es intencional en Go (para evitar que el código dependa
accidentalmente de un orden), y es diferente de, por ejemplo, los objetos de
JavaScript moderno que sí preservan orden de inserción. Si necesitas orden,
tienes que ordenar las llaves tú mismo (por ejemplo con el paquete `sort`).

```go
for clave, valor := range edades {
    // orden no garantizado
}
```

El detalle más importante — y el que más confunde a quien viene de otros
lenguajes — es el **"comma-ok" idiom**: cuando accedes a una llave que no
existe en el map, Go **no** te da un error ni un `nil` — te da el "valor
cero" del tipo del valor (`0` para `int`, `""` para `string`), lo cual puede
esconder un bug si no tienes cuidado (¿el precio de la manzana es
legítimamente 0, o simplemente no existe esa llave?). Para distinguir esos
dos casos, el acceso a un map puede devolver un **segundo valor booleano**:

```go
precio, existe := precios["manzana"]
if !existe {
    // la llave no está en el map
}
```

Este patrón (`valor, ok := ...`) es tan común en Go que se le conoce
específicamente como "comma-ok idiom", y lo vas a volver a ver en otros
contextos (type assertions, lectura de channels) en semanas futuras.

### Ejercicios

1. Crea un map de `string` a `int` que represente el inventario de una
   tienda (producto -> cantidad). Agrega 3 productos, imprime el mapa
   completo, y luego actualiza la cantidad de uno de ellos.
   Entrada inicial: `{"manzanas": 50, "peras": 30, "uvas": 20}`. Luego
   actualiza `"peras"` a `45`.

2. Dado el map anterior, usa `range` para recorrerlo e imprimir cada
   producto junto con su cantidad en el formato `producto: cantidad`.

3. Usa el "comma-ok idiom" para verificar si `"bananas"` existe en el
   inventario del ejercicio 1. Si no existe, imprime un mensaje indicando
   que el producto no está registrado; si existe, imprime su cantidad.
   Salida esperada: `El producto 'bananas' no está registrado`

4. Elimina `"uvas"` del inventario usando `delete`, e imprime el mapa antes
   y después de eliminarlo para confirmar el cambio.

5. **Mini-reto integrador de la semana** — combina TODO lo visto en los 7
   días: Vas a construir un mini "sistema de calificaciones" de consola.
   - Usa un `map[string]int` donde la llave es el nombre de un estudiante y
     el valor es su calificación (0-100).
   - Usa un ciclo `for` para pedir al usuario, repetidamente, que ingrese un
     nombre y una calificación, agregándolos al map. El ciclo debe
     detenerse cuando el usuario ingrese el nombre `"fin"` (usa `break`).
   - Escribe una función `letra(calificacion int) string` que use un
     `switch` sin expresión para convertir la calificación numérica en una
     letra (A/B/C/F, como en el ejercicio 3 del día 3).
   - Recorre el map con `range`, y para cada estudiante imprime su nombre,
     calificación y letra correspondiente (usando la función anterior).
   - Al final, calcula y muestra el **promedio general** de todas las
     calificaciones (vas a necesitar sumar todos los valores del map y
     dividir por la cantidad de estudiantes — un slice intermedio para
     guardar solo las calificaciones puede ayudarte, si quieres practicar
     también slices).
   - Como toque final, guarda en un slice de strings los nombres de los
     estudiantes que sacaron letra "A", y muéstralos al final como "Cuadro
     de honor".

   Entrada de ejemplo:
   ```
   Christian 95
   Ana 72
   Luis 58
   fin
   ```
   Salida esperada (formato libre, pero debe incluir estos datos):
   ```
   Christian: 95 (A)
   Ana: 72 (C)
   Luis: 58 (F)
   Promedio general: 75.00
   Cuadro de honor: [Christian]
   ```

💡 **Reto extra:** Extiende el mini-reto para que, en vez de sobrescribir la
calificación si el mismo nombre se ingresa dos veces, se guarde un **slice de
calificaciones por estudiante** (`map[string][]int`) y el promedio final se
calcule sobre el promedio de cada estudiante individualmente. Este ejercicio
te obliga a combinar maps y slices anidados, algo muy común al modelar datos
reales en un backend.
