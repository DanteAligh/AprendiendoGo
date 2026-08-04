// Bases de datos con database/sql + SQLite (driver modernc.org/sqlite,
// 100% Go: no necesita CGO ni un compilador de C instalado en Windows).
// Material de referencia — corresponde al día 18 de nuestro calendario.
//
// Este archivo es REFERENCIA DE SINTAXIS, no la solución de los ejercicios.
//
// El driver ya está en el go.mod del repositorio. Cómo correrlo, desde la
// raíz del repositorio en PowerShell:
//
//	go run ./cmd/extra/sqlite
//
// El CONCEPTO (abrir la conexión, crear la tabla, Exec/Query con
// placeholders) es idéntico sin importar el motor: cuando en la semana 4
// pasemos a PostgreSQL, lo único que cambia es el driver y que los
// placeholders son $1, $2 en vez de ?.
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite" // import en blanco: solo registra el driver "sqlite" ante database/sql
)

const archivoDB = "productos.db"

// Producto es la fila de nuestra tabla, representada como struct Go.
type Producto struct {
	ID     int64
	Nombre string
	Precio float64
	Stock  int
}

func main() {
	// Empezamos limpio en cada corrida, para que el ejemplo sea repetible.
	os.Remove(archivoDB)
	defer os.Remove(archivoDB)

	// sql.Open NO conecta de inmediato: solo prepara el *sql.DB (un pool
	// de conexiones, seguro para usar desde varias goroutines a la vez,
	// como las de un servidor HTTP del día 18/19).
	db, err := sql.Open("sqlite", archivoDB)
	if err != nil {
		log.Fatal("error abriendo la base de datos:", err)
	}
	defer db.Close()

	// Ping fuerza una conexión real, para detectar problemas temprano.
	if err := db.Ping(); err != nil {
		log.Fatal("error conectando a la base de datos:", err)
	}

	if err := crearTabla(db); err != nil {
		log.Fatal("error creando tabla:", err)
	}

	fmt.Println("--- 1 y 2. Crear tabla + insertar productos ---")
	id1, err := insertarProducto(db, Producto{Nombre: "Teclado", Precio: 450.50, Stock: 12})
	if err != nil {
		log.Fatal("error insertando:", err)
	}
	fmt.Printf("  producto insertado con id=%d\n", id1)

	id2, err := insertarProducto(db, Producto{Nombre: "Mouse", Precio: 220.00, Stock: 30})
	if err != nil {
		log.Fatal("error insertando:", err)
	}
	fmt.Printf("  producto insertado con id=%d\n", id2)

	id3, err := insertarProducto(db, Producto{Nombre: "Monitor", Precio: 3200.00, Stock: 5})
	if err != nil {
		log.Fatal("error insertando:", err)
	}
	fmt.Printf("  producto insertado con id=%d\n", id3)
	fmt.Println()

	fmt.Println("--- 3. Listar todos los productos ---")
	if err := listarProductos(db); err != nil {
		log.Fatal("error listando:", err)
	}
	fmt.Println()

	fmt.Println("--- 4. Actualizar stock de un producto ---")
	filasAfectadas, err := actualizarStock(db, id1, 8)
	if err != nil {
		log.Fatal("error actualizando:", err)
	}
	fmt.Printf("  filas afectadas al actualizar id=%d: %d\n", id1, filasAfectadas)

	// Intentamos actualizar un id que no existe, para ver RowsAffected en 0.
	filasAfectadas, err = actualizarStock(db, 9999, 1)
	if err != nil {
		log.Fatal("error actualizando:", err)
	}
	fmt.Printf("  filas afectadas al actualizar id=9999 (no existe): %d\n", filasAfectadas)
	fmt.Println()

	fmt.Println("--- listar después de actualizar ---")
	if err := listarProductos(db); err != nil {
		log.Fatal("error listando:", err)
	}
	fmt.Println()

	fmt.Println("--- 5. Eliminar un producto ---")
	filasAfectadas, err = eliminarProducto(db, id2)
	if err != nil {
		log.Fatal("error eliminando:", err)
	}
	fmt.Printf("  filas afectadas al eliminar id=%d: %d\n", id2, filasAfectadas)
	fmt.Println()

	fmt.Println("--- listar después de eliminar ---")
	if err := listarProductos(db); err != nil {
		log.Fatal("error listando:", err)
	}
	fmt.Println()

	fmt.Println("--- Reto extra: transacción ---")
	if err := insertarVariosEnTransaccion(db); err != nil {
		log.Fatal("error en transacción:", err)
	}
	if err := listarProductos(db); err != nil {
		log.Fatal("error listando:", err)
	}
}

// crearTabla usa "IF NOT EXISTS" para poder correrse varias veces sin
// error si la tabla ya existía.
func crearTabla(db *sql.DB) error {
	sqlCrear := `
	CREATE TABLE IF NOT EXISTS productos (
		id     INTEGER PRIMARY KEY AUTOINCREMENT,
		nombre TEXT    NOT NULL,
		precio REAL    NOT NULL,
		stock  INTEGER NOT NULL
	);`

	_, err := db.Exec(sqlCrear) // Exec: para sentencias que no devuelven filas
	return err
}

// insertarProducto usa placeholders (?) en vez de concatenar el SQL con
// datos del usuario. Esto evita SQL injection: el driver escapa los
// valores de forma segura y la sentencia se compila una sola vez
// (prepared statement) aunque la reutilices con distintos valores.
func insertarProducto(db *sql.DB, p Producto) (int64, error) {
	resultado, err := db.Exec(
		"INSERT INTO productos (nombre, precio, stock) VALUES (?, ?, ?)",
		p.Nombre, p.Precio, p.Stock,
	)
	if err != nil {
		return 0, err
	}
	return resultado.LastInsertId()
}

// listarProductos usa Query (no QueryRow) porque esperamos múltiples filas.
func listarProductos(db *sql.DB) error {
	rows, err := db.Query("SELECT id, nombre, precio, stock FROM productos ORDER BY id")
	if err != nil {
		return err
	}
	defer rows.Close() // SIEMPRE cerrar rows, o se quedan conexiones abiertas

	for rows.Next() {
		var p Producto
		if err := rows.Scan(&p.ID, &p.Nombre, &p.Precio, &p.Stock); err != nil {
			return err
		}
		fmt.Printf("  id=%d nombre=%s precio=%.2f stock=%d\n", p.ID, p.Nombre, p.Precio, p.Stock)
	}

	// rows.Err() reporta errores que hayan ocurrido DURANTE la iteración
	// (no solo al abrir la consulta) — es fácil olvidarlo y es buena
	// práctica revisarlo siempre.
	return rows.Err()
}

// actualizarStock usa QueryRow implícitamente vía Exec + RowsAffected para
// saber si el id realmente existía (RowsAffected == 0 significa que no).
func actualizarStock(db *sql.DB, id int64, nuevoStock int) (int64, error) {
	resultado, err := db.Exec("UPDATE productos SET stock = ? WHERE id = ?", nuevoStock, id)
	if err != nil {
		return 0, err
	}
	return resultado.RowsAffected()
}

func eliminarProducto(db *sql.DB, id int64) (int64, error) {
	resultado, err := db.Exec("DELETE FROM productos WHERE id = ?", id)
	if err != nil {
		return 0, err
	}
	return resultado.RowsAffected()
}

// insertarVariosEnTransaccion demuestra el reto extra del día 20: agrupar
// varias operaciones en una transacción. Si algo falla a la mitad, Rollback
// deshace TODO, dejando la base de datos como si nada hubiera pasado. En
// un ERP real esto es crítico para operaciones como "transferir stock
// entre dos almacenes": no puedes restar en uno y fallar antes de sumar en
// el otro, o el sistema queda en un estado inconsistente.
func insertarVariosEnTransaccion(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	productos := []Producto{
		{Nombre: "Cable HDMI", Precio: 90.00, Stock: 50},
		{Nombre: "Audífonos", Precio: 350.00, Stock: 15},
	}

	for _, p := range productos {
		if _, err := tx.Exec(
			"INSERT INTO productos (nombre, precio, stock) VALUES (?, ?, ?)",
			p.Nombre, p.Precio, p.Stock,
		); err != nil {
			tx.Rollback() // si algo falla, deshacemos todo lo hecho en esta transacción
			return err
		}
	}

	return tx.Commit() // confirma los cambios de forma atómica
}
