// Día 11: Interfaces — implementación implícita, any, type assertion y type switch
//
// Los ejercicios usan Sonido (Perro/Gato), Empleado (Fijo/PorHoras) y
// Procesar(dato any). Aquí mostramos el mismo concepto con un caso propio:
// un sistema de "notificaciones" que puede enviarse por distintos canales.
package main

import "fmt"

// Notificador es una interfaz pequeña, definida con exactamente el método
// que necesitamos consumir: Enviar. Cualquier tipo con este método la
// cumple automáticamente, sin declarar nada extra (duck typing).
type Notificador interface {
	Enviar(mensaje string) string
}

// EmailNotificador y SMSNotificador son dos tipos concretos, completamente
// independientes entre sí, que casualmente cumplen la misma interfaz.
type EmailNotificador struct {
	Destinatario string
}

func (e EmailNotificador) Enviar(mensaje string) string {
	return fmt.Sprintf("Correo enviado a %s: \"%s\"", e.Destinatario, mensaje)
}

type SMSNotificador struct {
	Numero string
}

func (s SMSNotificador) Enviar(mensaje string) string {
	return fmt.Sprintf("SMS enviado a %s: \"%s\"", s.Numero, mensaje)
}

// notificar acepta CUALQUIER cosa que cumpla Notificador. No le importa si
// es un Email o un SMS -- eso es justo la ventaja de programar contra la
// interfaz y no contra el tipo concreto.
func notificar(n Notificador, mensaje string) {
	fmt.Println(n.Enviar(mensaje))
}

// describirCanal usa TYPE SWITCH para reaccionar distinto según el tipo
// concreto que hay "adentro" de la interfaz, cuando sí nos importa el
// detalle específico de cada canal (por ejemplo, para logging o métricas).
func describirCanal(n Notificador) {
	switch v := n.(type) {
	case EmailNotificador:
		fmt.Println("  -> Canal: Email, destinatario:", v.Destinatario)
	case SMSNotificador:
		fmt.Println("  -> Canal: SMS, número:", v.Numero)
	default:
		fmt.Println("  -> Canal desconocido")
	}
}

// resumirEntrada usa "any" (interface{}) para aceptar literalmente cualquier
// valor, y un type switch para decidir qué hacer según su tipo dinámico.
// Se usa con moderación: aquí tiene sentido porque de verdad no sabemos de
// antemano qué tipo de dato puede llegar (por ejemplo, un valor leído de un
// JSON genérico en un backend real).
func resumirEntrada(valor any) string {
	switch v := valor.(type) {
	case int:
		return fmt.Sprintf("es un entero: %d", v)
	case string:
		return fmt.Sprintf("es un texto de %d caracteres", len(v))
	case bool:
		return fmt.Sprintf("es un booleano: %t", v)
	default:
		return "tipo no reconocido"
	}
}

func main() {
	correo := EmailNotificador{Destinatario: "soporte@empresa.com"}
	sms := SMSNotificador{Numero: "+52 555 123 4567"}

	notificar(correo, "Tu pedido fue enviado")
	notificar(sms, "Tu pedido fue enviado")

	// Un slice de la interfaz puede mezclar tipos concretos distintos.
	canales := []Notificador{correo, sms}
	fmt.Println("\nDetalle de cada canal:")
	for _, c := range canales {
		describirCanal(c)
	}

	// Type assertion clásica con ", ok": la forma segura, sin riesgo de panic.
	var n Notificador = correo
	if email, ok := n.(EmailNotificador); ok {
		fmt.Println("\nConfirmado por type assertion: es un EmailNotificador dirigido a", email.Destinatario)
	}
	if _, ok := n.(SMSNotificador); !ok {
		fmt.Println("Confirmado: NO es un SMSNotificador (ok fue false)")
	}

	// any / interface{} con type switch:
	fmt.Println("\nProbando resumirEntrada con distintos tipos:")
	entradas := []any{42, "hola mundo", true, 3.14}
	for _, e := range entradas {
		fmt.Println("-", resumirEntrada(e))
	}
}
