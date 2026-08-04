// Día 16: manejo de archivos con os, io, bufio y encoding/csv.
//
// Este archivo es REFERENCIA DE SINTAXIS, no la solución de los ejercicios.
// Crea archivos temporales en el directorio actual al correrlo. Corre con:
//
//	go run ./cmd/extra/archivos
package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const archivoTexto = "ejemplo_dia16.txt"
const archivoCSV = "ejemplo_dia16.csv"

func main() {
	defer limpiar() // borra los archivos de ejemplo al terminar, para no ensuciar la carpeta

	demoEscribirTexto()
	demoLeerTexto()
	demoEscribirCSV()
	demoLeerCSV()
	demoIoReader()
}

// -----------------------------------------------------------------------
// 1. Escribir un archivo de texto con os + bufio.
// -----------------------------------------------------------------------
//
// os.Create trunca (vacía) el archivo si ya existe, o lo crea si no existe.
// Go no tiene "excepciones": open/create devuelven (valor, error) y el
// error es tuyo para revisar explícitamente.
func demoEscribirTexto() {
	fmt.Println("--- 1. Escribir archivo de texto ---")

	archivo, err := os.Create(archivoTexto)
	if err != nil {
		fmt.Println("  error creando archivo:", err)
		return
	}
	defer archivo.Close() // se cierra al salir de la función, pase lo que pase

	// bufio.NewWriter evita hacer una llamada al sistema operativo por cada
	// escritura; acumula en un buffer y lo manda en bloque con Flush().
	escritor := bufio.NewWriter(archivo)
	lineas := []string{
		"primera línea",
		"segunda línea",
		"tercera línea",
	}
	for _, linea := range lineas {
		fmt.Fprintln(escritor, linea) // Fprintln agrega salto de línea al final
	}

	if err := escritor.Flush(); err != nil { // IMPORTANTE: sin Flush, puede que nada llegue al disco
		fmt.Println("  error al escribir:", err)
		return
	}
	fmt.Println("  archivo escrito:", archivoTexto)
	fmt.Println()
}

// -----------------------------------------------------------------------
// 2. Leer un archivo de texto línea por línea con bufio.Scanner.
// -----------------------------------------------------------------------
func demoLeerTexto() {
	fmt.Println("--- 2. Leer archivo de texto ---")

	archivo, err := os.Open(archivoTexto) // Open es solo lectura
	if err != nil {
		fmt.Println("  error abriendo archivo:", err)
		return
	}
	defer archivo.Close()

	escaner := bufio.NewScanner(archivo)
	numero := 1
	for escaner.Scan() { // avanza línea por línea
		fmt.Printf("  línea %d: %s\n", numero, escaner.Text())
		numero++
	}
	if err := escaner.Err(); err != nil {
		fmt.Println("  error leyendo:", err)
	}
	fmt.Println()
}

// -----------------------------------------------------------------------
// 3. Escribir un CSV con encoding/csv.
// -----------------------------------------------------------------------
func demoEscribirCSV() {
	fmt.Println("--- 3. Escribir CSV ---")

	archivo, err := os.Create(archivoCSV)
	if err != nil {
		fmt.Println("  error creando csv:", err)
		return
	}
	defer archivo.Close()

	escritor := csv.NewWriter(archivo)

	filas := [][]string{
		{"nombre", "precio", "stock"}, // encabezado
		{"Teclado", "450.50", "12"},
		{"Mouse", "220.00", "30"},
		{"Monitor", "3200.00", "5"},
		{"Cable HDMI", "90.00", "50"},
	}

	for _, fila := range filas {
		if err := escritor.Write(fila); err != nil {
			fmt.Println("  error escribiendo fila:", err)
			return
		}
	}
	escritor.Flush() // como con bufio, hay que vaciar el buffer explícitamente

	if err := escritor.Error(); err != nil {
		fmt.Println("  error al finalizar csv:", err)
		return
	}
	fmt.Println("  csv escrito:", archivoCSV)
	fmt.Println()
}

// -----------------------------------------------------------------------
// 4. Leer un CSV con encoding/csv y calcular el valor total del inventario.
// -----------------------------------------------------------------------
func demoLeerCSV() {
	fmt.Println("--- 4. Leer CSV y calcular valor total ---")

	archivo, err := os.Open(archivoCSV)
	if err != nil {
		fmt.Println("  error abriendo csv:", err)
		return
	}
	defer archivo.Close()

	lector := csv.NewReader(archivo)

	filas, err := lector.ReadAll() // lee todo el CSV de una vez a [][]string
	if err != nil {
		fmt.Println("  error leyendo csv:", err)
		return
	}

	valorTotal := 0.0
	// filas[0] es el encabezado (nombre,precio,stock), lo saltamos.
	for _, fila := range filas[1:] {
		precio, err := strconv.ParseFloat(fila[1], 64)
		if err != nil {
			fmt.Println("  precio inválido, se omite fila:", fila)
			continue
		}
		stock, err := strconv.Atoi(fila[2])
		if err != nil {
			fmt.Println("  stock inválido, se omite fila:", fila)
			continue
		}
		valorTotal += precio * float64(stock)
	}

	fmt.Printf("  valor total del inventario: %.2f\n", valorTotal)
	fmt.Println()
}

// -----------------------------------------------------------------------
// 5. Programar contra io.Reader, no contra *os.File.
// -----------------------------------------------------------------------
//
// contarPalabras funciona igual con un archivo, un string en memoria, o
// cualquier otra fuente que implemente io.Reader. Esa es la elegancia de
// las interfaces de Go: la función no sabe (ni le importa) de dónde vienen
// los bytes.
func contarPalabras(r io.Reader) (int, error) {
	datos, err := io.ReadAll(r)
	if err != nil {
		return 0, err
	}
	texto := strings.TrimSpace(string(datos))
	if texto == "" {
		return 0, nil
	}
	palabras := strings.Fields(texto) // separa por espacios/saltos de línea
	return len(palabras), nil
}

func demoIoReader() {
	fmt.Println("--- 5. Misma función, dos fuentes distintas (io.Reader) ---")

	// Fuente 1: un archivo real.
	archivo, err := os.Open(archivoTexto)
	if err != nil {
		fmt.Println("  error abriendo archivo:", err)
		return
	}
	defer archivo.Close()

	n, err := contarPalabras(archivo)
	if err != nil {
		fmt.Println("  error contando palabras del archivo:", err)
	} else {
		fmt.Printf("  palabras en el archivo: %d\n", n)
	}

	// Fuente 2: un string en memoria, sin tocar disco.
	lector := strings.NewReader("esto es una prueba de contar palabras en memoria")
	n2, err := contarPalabras(lector)
	if err != nil {
		fmt.Println("  error contando palabras del string:", err)
	} else {
		fmt.Printf("  palabras en el string: %d\n", n2)
	}
}

// limpiar borra los archivos de ejemplo generados por este programa.
func limpiar() {
	os.Remove(archivoTexto)
	os.Remove(archivoCSV)
}
