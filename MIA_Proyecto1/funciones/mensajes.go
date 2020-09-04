package funciones

import (
	"fmt"
	"log"
	"os"
)

func Mensaje(mensaje string, tipo int8) {
	if tipo == 0 {
		salida := fmt.Sprintf("\033[1;36m%s\033[0m", mensaje)
		fmt.Print(salida)

	} else if tipo == 1 {
		salida := fmt.Sprintf("\033[1;32m%s\n\033[0m", mensaje)
		fmt.Println(salida)

	} else if tipo == 2 {
		temp := "ERROR: " + mensaje
		salida := fmt.Sprintf("\033[1;31m%s\n\033[0m", temp)
		fmt.Println(salida)
	}
}

func escribirBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)

	if err != nil {
		log.Fatal(err)
	}
}

func leerBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number) //array de bytes

	_, err := file.Read(bytes) // Leido -> bytes
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}
