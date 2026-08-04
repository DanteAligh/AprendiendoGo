// Día 8: Structs
//
// Este ejemplo NO resuelve los ejercicios del día (esos usan Producto,
// Punto, Auto, Articulo...). Aquí modelamos algo distinto: una pequeña
// ficha de "empleado de soporte técnico" para mostrar la sintaxis de
// structs, instanciación, structs anidados y structs embebidos.
package main

import "fmt"

// Direccion es un struct simple que vamos a anidar dentro de otro.
type Direccion struct {
	Calle  string
	Ciudad string
	CP     string
}

// DatosPersonales agrupa información básica de una persona.
// Lo vamos a EMBEBER (no anidar con nombre) dentro de Empleado más abajo,
// para que veas la diferencia práctica entre las dos técnicas.
type DatosPersonales struct {
	Nombre string
	Edad   int
}

// Empleado combina:
//   - un campo normal (Puesto)
//   - un campo anidado CON nombre (Direccion Direccion) -> acceso vía .Direccion.Calle
//   - un campo embebido SIN nombre (DatosPersonales)      -> acceso directo vía .Nombre
type Empleado struct {
	DatosPersonales // embebido: sus campos se "promueven" a Empleado
	Puesto          string
	Direccion       Direccion // anidado con nombre: hay que pasar por .Direccion
}

func main() {
	// --- Instanciación con literal y nombres de campo (la forma recomendada) ---
	// Usar nombres de campo hace que el código no se rompa si el struct cambia
	// de orden en el futuro, y es mucho más legible que un literal posicional.
	e1 := Empleado{
		DatosPersonales: DatosPersonales{Nombre: "Marta", Edad: 34},
		Puesto:          "Soporte Nivel 2",
		Direccion: Direccion{
			Calle:  "Av. Siempre Viva 742",
			Ciudad: "Springfield",
			CP:     "00000",
		},
	}

	// Gracias al embebido, accedemos a Nombre y Edad DIRECTO, sin pasar por
	// e1.DatosPersonales.Nombre (aunque esa forma larga también funciona).
	fmt.Println("Empleado:", e1.Nombre, "-", e1.Edad, "años")
	// El campo anidado CON nombre sí necesita el camino completo:
	fmt.Println("Vive en:", e1.Direccion.Calle, e1.Direccion.Ciudad)

	// --- new(): reserva memoria para un struct "vacío" (zero value) y
	// devuelve un puntero a él. Se usa poco en la práctica, pero es bueno
	// reconocerlo. e2 es de tipo *Empleado.
	e2 := new(Empleado)
	fmt.Printf("Empleado recién creado con new(): %+v\n", *e2)
	// Nota: los campos string quedan en "" y los int en 0 (zero values).

	// Podemos llenarlo después, como si fuera cualquier struct (Go
	// desreferencia automáticamente para asignar campos vía puntero).
	e2.Nombre = "Carlos"
	e2.Edad = 29
	e2.Puesto = "Soporte Nivel 1"
	fmt.Println("Ahora e2 tiene nombre:", e2.Nombre)

	// --- var sin inicializar: también arranca en zero value ---
	var e3 Empleado
	fmt.Printf("e3 (var sin inicializar): %+v\n", e3)

	// --- %+v muestra nombres de campo, útil para depurar rápido ---
	fmt.Printf("e1 completo: %+v\n", e1)

	// --- Un slice de structs, muy común al modelar listas de registros ---
	equipo := []Empleado{
		e1,
		*e2,
		{DatosPersonales: DatosPersonales{Nombre: "Iris", Edad: 41}, Puesto: "Lead"},
	}

	fmt.Println("\nEquipo de soporte:")
	for _, emp := range equipo {
		fmt.Printf("- %s (%d años), puesto: %s\n", emp.Nombre, emp.Edad, emp.Puesto)
	}
}
