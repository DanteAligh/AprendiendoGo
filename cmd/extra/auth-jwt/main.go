// Autenticación: hashing de contraseñas (bcrypt) y JWT.
// Material de referencia — corresponde al día 24 de nuestro calendario.
//
// Este archivo usa DOS librerías EXTERNAS (no son de la librería estándar):
// golang.org/x/crypto/bcrypt y github.com/golang-jwt/jwt/v5. Ya están
// declaradas en el go.mod del repositorio, así que no tienes que instalar
// nada. Si alguna vez falla con "missing go.sum entry", con internet:
//
//	go mod tidy
//
// Cómo correrlo (desde la raíz del repositorio, en PowerShell):
//
//	go run ./cmd/extra/auth-jwt
//
// ---------------------------------------------------------------------------
// ¿POR QUÉ NUNCA GUARDAMOS CONTRASEÑAS EN TEXTO PLANO?
//
// Si tu base de datos de un ERP se filtra alguna vez (pasa, incluso a
// empresas grandes), NO querés que la tabla de usuarios tenga la columna
// "contraseña" con el valor real de cada contraseña. bcrypt es un algoritmo
// de HASHING diseñado específicamente para contraseñas: es lento a
// propósito (tiene un "costo" configurable) para que probar millones de
// combinaciones (fuerza bruta) sea carísimo computacionalmente. Además,
// genera un "salt" (dato aleatorio) distinto cada vez, así que la misma
// contraseña "1234" nunca produce el mismo hash dos veces — eso evita
// ataques con tablas precalculadas (rainbow tables).
//
// Importante: bcrypt NO es "cifrado" (encryption). El cifrado se puede
// revertir con una clave; el hashing NO se puede revertir. Para verificar
// un login, no "desencriptás" el hash: volvés a hashear la contraseña que
// te mandaron y le pedís a bcrypt que compare los dos hashes.
//
// ---------------------------------------------------------------------------
// ¿QUÉ ES UN JWT (JSON Web Token)?
//
// Un JWT es un string (con 3 partes separadas por puntos: header.payload.
// firma) que un servidor emite después de un login exitoso, y que el
// cliente (un frontend, una app móvil, otro servicio) manda de vuelta en
// cada petición futura (típicamente en la cabecera
// "Authorization: Bearer <token>") para probar quién es sin tener que
// volver a mandar usuario/contraseña en cada request.
//
// El payload contiene "claims" (afirmaciones): quién es el usuario, qué rol
// tiene, cuándo expira el token, etc. La FIRMA es lo que hace que el token
// sea confiable: el servidor la genera con una clave secreta que solo él
// conoce. Si alguien modifica el payload (por ejemplo, para cambiar su rol
// de "empleado" a "admin"), la firma ya no coincide y el servidor rechaza
// el token al validarlo. El JWT en sí NO es secreto ni está cifrado —
// cualquiera puede leer su contenido (es solo base64) — lo que lo protege
// es que no se puede FALSIFICAR sin la clave secreta del servidor.
//
// En un ERP real, este flujo es: login -> se verifica usuario/contraseña
// con bcrypt -> si es válido, se firma un JWT con los datos del usuario ->
// el cliente guarda ese token -> en cada petición futura, un middleware
// (como los del día 23) valida el token ANTES de dejar pasar la petición al
// handler real.
package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// claveSecreta firma y valida los tokens. En un proyecto real esto NUNCA va
// hardcodeado en el código: se lee desde una variable de entorno (ver día
// 23), y debe ser larga, aleatoria y distinta en cada entorno (desarrollo,
// producción). Acá va fija solo para que el ejemplo sea reproducible.
var claveSecreta = []byte("clave-secreta-de-ejemplo-cambiar-en-produccion")

// ClaimsUsuario define qué información viaja DENTRO del JWT. Además de
// nuestros propios campos (Usuario, Rol), incrustamos jwt.RegisteredClaims,
// que trae campos estándar de la especificación JWT como ExpiresAt
// (expiración) e IssuedAt (cuándo se emitió). Usar los campos estándar
// permite que cualquier librería JWT de cualquier lenguaje entienda tu
// token.
type ClaimsUsuario struct {
	Usuario string `json:"usuario"`
	Rol     string `json:"rol"`
	jwt.RegisteredClaims
}

// hashearContrasena envuelve bcrypt.GenerateFromPassword. bcrypt.DefaultCost
// (actualmente 10) es un buen punto de partida: más alto = más lento de
// calcular = más seguro contra fuerza bruta, pero también más lento para tu
// propio servidor en cada login. Se ajusta según cuánta capacidad de cómputo
// tenga tu servidor.
func hashearContrasena(contrasena string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(contrasena), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hasheando contraseña: %w", err)
	}
	return string(hash), nil
}

// verificarContrasena compara una contraseña en texto plano (la que llega
// en un intento de login) contra un hash bcrypt ya guardado (el que
// tendrías, por ejemplo, en una columna de tu base de datos de usuarios).
func verificarContrasena(hash, contrasena string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(contrasena))
	return err == nil
}

// generarToken crea y firma un JWT para un usuario autenticado. El token
// expira en 1 hora: los JWT casi siempre tienen una expiración corta,
// porque si un token se filtra, un atacante solo podría usarlo mientras
// siga vigente (no hay forma de "revocar" un JWT ya emitido sin
// infraestructura adicional).
func generarToken(usuario, rol string) (string, error) {
	claims := ClaimsUsuario{
		Usuario: usuario,
		Rol:     rol,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// jwt.NewWithClaims arma el token en memoria; SignedString lo firma
	// con nuestra clave secreta usando el algoritmo HS256 (HMAC-SHA256) y
	// devuelve el string final "header.payload.firma".
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(claveSecreta)
}

// validarToken parsea un string JWT y verifica su firma. Si alguien
// modificó el payload sin conocer claveSecreta, jwt.ParseWithClaims
// devuelve un error acá mismo — nunca deberías confiar en un token sin
// pasar por esta validación.
func validarToken(tokenString string) (*ClaimsUsuario, error) {
	claims := &ClaimsUsuario{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		// Esta función le dice a la librería CÓMO obtener la clave para
		// verificar la firma. En sistemas más grandes (ej. login con
		// Google), acá se buscaría la clave pública correspondiente al
		// "key id" del token. En nuestro caso simple, siempre es la misma
		// clave secreta compartida.
		return claveSecreta, nil
	})
	if err != nil {
		return nil, fmt.Errorf("token inválido: %w", err)
	}
	if !token.Valid {
		return nil, errors.New("token no es válido")
	}
	return claims, nil
}

func main() {
	fmt.Println("=== 1. Hashing de contraseñas con bcrypt ===")
	contrasenaOriginal := "miContraseñaSegura123"

	hash, err := hashearContrasena(contrasenaOriginal)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Contraseña original: %s\n", contrasenaOriginal)
	fmt.Printf("Hash guardado (esto es lo que iría en la base de datos): %s\n", hash)

	fmt.Println("\n--- Simulando un intento de login ---")
	fmt.Println("Intento con la contraseña correcta:",
		verificarContrasena(hash, "miContraseñaSegura123"))
	fmt.Println("Intento con una contraseña incorrecta:",
		verificarContrasena(hash, "otraContrasena"))

	fmt.Println("\n=== 2. Generar y validar un JWT ===")
	token, err := generarToken("ana", "admin")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Token generado:\n%s\n", token)

	fmt.Println("\n--- Validando el token (caso correcto) ---")
	claims, err := validarToken(token)
	if err != nil {
		fmt.Println("Error validando:", err)
	} else {
		fmt.Printf("Token válido. Usuario: %s | Rol: %s | Expira: %s\n",
			claims.Usuario, claims.Rol, claims.ExpiresAt.Time.Format(time.RFC3339))
	}

	fmt.Println("\n--- Validando un token manipulado (debe fallar) ---")
	// Simulamos que alguien interceptó el token y le cambió un carácter del
	// payload (la parte del medio) para intentar hacerse pasar por admin
	// sin conocer la clave secreta. La verificación de firma debe rechazarlo.
	tokenManipulado := token[:len(token)-5] + "AAAAA"
	_, err = validarToken(tokenManipulado)
	if err != nil {
		fmt.Println("Como se esperaba, el token manipulado fue rechazado:", err)
	} else {
		fmt.Println("ALERTA: esto no debería pasar, el token manipulado fue aceptado")
	}
}
