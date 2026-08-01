package main

import "fmt"

func calcularIMC(peso, altura float64) (float64, bool) {

	if altura <= 0 {

		return 0, false
	}
	imc := peso / (altura * altura)
	return imc, true

}
func clasificarIMC(imc float64) string {
	if imc < 18.5 {
		return "Bajo peso"
	}
	if imc < 25 {
		return "Normal"
	}
	if imc < 30 {
		return "Sobrepeso"
	}
	return "Obesidad"
}
func main() {
	peso := 82.5
	altura := 1.78

	IMC, ok := calcularIMC(peso, altura)
	if !ok {
		fmt.Println("Altura no puede ser 0")
		return
	}
	clasificacion := clasificarIMC(IMC)
	fmt.Printf("Peso %.2f, Altura %.2f, IMC %.2f, Clasificación %s", peso, altura, IMC, clasificacion)
}
