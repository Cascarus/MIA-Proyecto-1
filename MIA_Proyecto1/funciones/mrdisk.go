package funciones

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type mrdisk struct {
	path string
}

func NewMRDisk(path string) mrdisk {
	e := mrdisk{path}
	//fmt.Println("Creo el nuevo MRDIS")
	return e
}

func (e mrdisk) Ejecutar() {

	//fmt.Println("la ruta es: ", e.path)
	if _, err := os.Stat(e.path); os.IsNotExist(err) {
		Mensaje("No existe un disco con ese nombre en el directorio", 2)
		return

	} else {
		Mensaje("Desea eliminar el disco(Y/N): ", 0)
		reader := bufio.NewReader(os.Stdin)
		comando, _ := reader.ReadString('\n')
		com := strings.TrimRight(comando, "\n")

		if strings.ToLower(com) == "y" {
			r := os.Remove(e.path)

			if r != nil {
				log.Fatal(r)
			}
			Mensaje("El disco se ha eliminado exitosamente", 1)
			return
		}
		fmt.Println("")
	}

}
