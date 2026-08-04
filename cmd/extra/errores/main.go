// Día 12: Manejo de errores idiomático
//
// Los ejercicios usan Dividir, ValidarEdad, Vender/ErrStockInsuficiente,
// CargarPedido/ProcesarPedido y ErrorValidacion. Aquí mostramos las mismas
// técnicas con un caso distinto: validar y procesar el login de un usuario.
package main

import (
	"errors"
	"fmt"
)

// ErrCredencialesInvalidas es un ERROR CENTINELA: un valor de error fijo y
// reconocible a nivel de paquete. Declararlo como variable (en vez de crear
// un errors.New nuevo cada vez) es lo que permite compararlo después con
// errors.Is, incluso si viaja envuelto dentro de otros errores.
var ErrCredencialesInvalidas = errors.New("usuario o contraseña incorrectos")

// ErrorDeCampo es un tipo de error PROPIO, con datos extra (Campo). Cumple
// la interfaz "error" porque tiene un método Error() string.
type ErrorDeCampo struct {
	Campo   string
	Detalle string
}

func (e *ErrorDeCampo) Error() string {
	return fmt.Sprintf("campo '%s': %s", e.Campo, e.Detalle)
}

// validarFormulario devuelve un *ErrorDeCampo (con datos estructurados) si
// algo del formulario de login está mal formado, ANTES de siquiera intentar
// autenticar. Separar "validación de forma" de "validación de credenciales"
// es un patrón común en backends reales.
func validarFormulario(usuario, contrasena string) error {
	if usuario == "" {
		return &ErrorDeCampo{Campo: "usuario", Detalle: "no puede estar vacío"}
	}
	if len(contrasena) < 4 {
		return &ErrorDeCampo{Campo: "contrasena", Detalle: "debe tener al menos 4 caracteres"}
	}
	return nil
}

// autenticar simula una verificación contra una base de datos: solo un
// usuario/contraseña "correcto" pasa. Cualquier otro caso devuelve el error
// centinela ErrCredencialesInvalidas.
func autenticar(usuario, contrasena string) error {
	if usuario == "admin" && contrasena == "1234" {
		return nil
	}
	return ErrCredencialesInvalidas
}

// login es la función de "capa superior": llama a las dos funciones
// anteriores y, si autenticar falla, ENVUELVE el error centinela con
// contexto extra usando %w. Envolver no reemplaza el error original: lo
// preserva dentro de la cadena para que errors.Is lo pueda encontrar más
// adelante, aunque el mensaje visible tenga más contexto.
func login(usuario, contrasena string) error {
	if err := validarFormulario(usuario, contrasena); err != nil {
		return fmt.Errorf("login: formulario inválido: %w", err)
	}
	if err := autenticar(usuario, contrasena); err != nil {
		return fmt.Errorf("login: no se pudo autenticar a '%s': %w", usuario, err)
	}
	return nil
}

func intentarLogin(usuario, contrasena string) {
	err := login(usuario, contrasena)

	// El patrón central de Go: comprobar el error explícitamente.
	if err == nil {
		fmt.Printf("Login de '%s' exitoso.\n", usuario)
		return
	}

	fmt.Println("Error recibido:", err)

	// errors.Is: ¿en el fondo de esta cadena está el error centinela?
	// Funciona aunque el error haya sido envuelto varias veces con %w.
	if errors.Is(err, ErrCredencialesInvalidas) {
		fmt.Println("  -> Diagnóstico: credenciales incorrectas, sugerir reintento.")
	}

	// errors.As: ¿en el fondo de esta cadena hay un error de ESTE TIPO
	// concreto? Si sí, nos lo asigna en campoErr para leer sus campos.
	var campoErr *ErrorDeCampo
	if errors.As(err, &campoErr) {
		fmt.Printf("  -> Diagnóstico: error de formulario en el campo '%s' (%s)\n", campoErr.Campo, campoErr.Detalle)
	}
	fmt.Println()
}

func main() {
	intentarLogin("admin", "1234")       // caso exitoso
	intentarLogin("", "abc")             // falla validación de formulario (usuario vacío)
	intentarLogin("admin", "abc")        // falla validación de formulario (contraseña corta)
	intentarLogin("admin", "clave-mala") // pasa validación, falla autenticación
}
