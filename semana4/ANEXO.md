# Anexo — Cómo investigar por tu cuenta (semana 4)

En esta semana usas librerías de **terceros** (no de la librería
estándar): JWT, bcrypt, y en tu futuro ERP probablemente un SDK de IA.
`go doc` sigue funcionando, pero ya no basta — estas librerías viven en
GitHub y su documentación real suele estar en el `README`, no solo en los
comentarios del código.

## Paso 1: ¿esta librería es de fiar?

Antes de instalar cualquier paquete de un desconocido, revisa en
pkg.go.dev (buscando `pkg.go.dev/<módulo>`, por ejemplo
`pkg.go.dev/github.com/golang-jwt/jwt/v5`):

- **"Imported by"**: cuántos otros proyectos la usan. Números altos =
  mucha gente ya la probó en producción.
- **Fecha de la última versión**: si no se actualizó en años y tiene
  issues de seguridad abiertos, piénsalo dos veces.
- El link al repositorio de GitHub, para leer el `README` y ver si tiene
  ejemplos claros de uso.

## Paso 2: `go get` y el significado de `go.mod` / `go.sum`

```powershell
go get github.com/golang-jwt/jwt/v5
```

Esto agrega la dependencia a tu `go.mod` (qué versión usas) y a tu
`go.sum` (un hash de verificación de esa versión exacta, para que nadie
pueda "cambiarte" el código de la librería sin que Go lo note). Nunca
edites `go.sum` a mano; si algo se corrompe, borra ambos archivos de
dependencias y corre `go mod tidy` para que Go los reconstruya.

```powershell
go mod tidy
```

Este comando limpia dependencias que ya no usas y agrega las que faltan,
comparando tu código real contra lo que declaras en `go.mod`. Córrelo
después de agregar o quitar imports.

## Paso 3: cuando la librería no tiene `go doc` útil

Muchas librerías de terceros documentan bien sus tipos (`go doc` sigue
funcionando sobre cualquier módulo descargado), pero el **flujo de uso**
(en qué orden llamas las cosas) casi siempre está en el `README` de
GitHub, no en los comentarios de una función suelta. Para JWT, por
ejemplo, `go doc` te dice qué hace `jwt.NewWithClaims`, pero el README te
muestra el flujo completo: crear claims -> firmar -> obtener el string ->
en otro lado, parsear y validar.

Regla práctica: **`go doc` para la firma exacta de una función puntual,
el README de GitHub para entender el flujo completo de la librería.**

## Paso 4: variables de entorno y secretos

Nunca hardcodees una clave secreta (la de firmar JWT, una API key de un
proveedor de IA) directo en el código, ni la subas a git. El patrón
estándar:

```powershell
$env:MI_CLAVE_SECRETA = "algo-largo-y-aleatorio"
```

(Eso es PowerShell. En Linux o dentro de un contenedor Docker, la misma idea
se escribe `export MI_CLAVE_SECRETA="..."`.)

Y en Go la lees con `os.Getenv("MI_CLAVE_SECRETA")`. En desarrollo local
se suele usar un archivo `.env` (que **nunca** se commitea — por eso ya
está en tu `.gitignore`) y una librería como `godotenv` para cargarlo. En
producción, las variables de entorno las configura la plataforma donde
despliegues (no un archivo).

## Esto es exactamente lo que vas a hacer con una API de IA

El día que conectes tu ERP a una API de IA real (Anthropic, OpenAI, la que
sea), el proceso es el mismo que en el ejemplo `cmd\extra\cliente-http`: la clave de API va en una
variable de entorno, nunca en el código; revisas el README/SDK oficial
para el flujo de autenticación; y usas `go doc` solo para las firmas
puntuales de las funciones del SDK. No hay magia adicional — es el mismo
patrón que ya practicaste esta semana, a mayor escala.
