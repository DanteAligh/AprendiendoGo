package main

import "fmt"

func convertirTiempo(totalSegundos int) (int, int, int) {
	horas := totalSegundos / 3600
	minutos := (totalSegundos % 3600) / 60
	segundos := totalSegundos % 60
	return horas, minutos, segundos
}

func main() {
	t := []int{9875, 59}
	for _, totalSegundos := range t {
		horas, minutos, segundos := convertirTiempo(totalSegundos)
		fmt.Printf("%d: segundos son %d horas: %d minutos: %d s\n", totalSegundos, horas, minutos, segundos)
	}
}
