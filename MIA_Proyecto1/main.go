package main

import (
	"MIA_Proyecto1/funciones"
	"MIA_Proyecto1/interprete"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	reader := bufio.NewReader(os.Stdin)

	for {
		funciones.Mensaje(">", 0)
		comando, _ := reader.ReadString('\n')
		com := strings.TrimRight(comando, "\n")

		if strings.ToLower(com) == "exit" {
			break
		}

		inter := interprete.New(com)
		inter.Ejecutar()

	}

	fmt.Println("----------------------------El programa ha finalizado----------------------------")

}
