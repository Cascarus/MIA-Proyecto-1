package funciones

import (
	"MIA_Proyecto1/estructuras"
	"bytes"
	"encoding/binary"
	"log"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

var montados []string
var discos []string
var name_particiones []string

func Mount(path string, name string) {
	id := generar_Nombre(path)
	montar := path + "_" + name + "_" + id

	if verificar_particion(path, name) {
		montados = append(montados, montar)

		Mensaje("Se ha montado la particion con exito!", 1)
		Mensaje("\tLas particiones montadas existentes son: ", 1)

		for i := 0; i < len(montados); i++ {
			if montados[i] != "" {
				dividir := strings.Split(montados[i], "_")
				Mensaje("\t path="+dividir[0]+" name="+dividir[1]+" id="+dividir[2], 1)
			}
		}
	} else {
		Mensaje("La particion no existe en ese disco", 2)
	}

}

func Unmount(id string) {
	encontrado := false
	ubicacion := 0

	for i := 0; i < len(montados); i++ {
		ubicacion = i
		montado := strings.Split(montados[i], "_")
		if montado[2] == strings.ToLower(id) {
			encontrado = true
			break
		}
	}

	if encontrado {
		montados[ubicacion] = ""

		Mensaje("Se ha desmontado la particion", 1)
		Mensaje("\tLas particiones montadas existentes son: ", 1)

		for i := 0; i < len(montados); i++ {
			if montados[i] != "" {
				dividir := strings.Split(montados[i], "_")
				Mensaje("\t path="+dividir[0]+" name="+dividir[1]+" id="+dividir[2], 1)
			}
		}
	}
}

func Mostrar_mounts() {
	Mensaje("Se ha montado la particion con exito!", 1)
	Mensaje("\tLas particiones montadas existentes son: ", 1)

	for i := 0; i < len(montados); i++ {
		if montados[i] != "" {
			dividir := strings.Split(montados[i], "_")
			Mensaje("\t path="+dividir[0]+" name="+dividir[1]+" id="+dividir[2], 1)
		}
	}
}

func generar_Nombre(path string) string {
	ruta := strings.Split(path, "/")
	nombre_D := strings.Split(ruta[len(ruta)-1], ".")
	nombre := nombre_D[0]

	discos_existentes := strings.Split(agregar(nombre), "_")

	//var convertir bool
	ch := 'a'
	caracter, err := strconv.ParseInt(discos_existentes[1], 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("paso el parseint de caracter")
	caracter += int64(ch)
	//ch = caracter

	num_part, err := strconv.ParseInt(discos_existentes[2], 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("paso el parseint de numpart")

	nombre = "vd" + string(caracter)
	nombre += discos_existentes[2]
	num_part++

	nuevo := discos_existentes[0] + "_" + discos_existentes[1] + "_" + strconv.FormatInt(num_part, 10)

	caracter, err = strconv.ParseInt(discos_existentes[1], 10, 64)
	discos[caracter] = nuevo

	//fmt.Println(nombre)
	return nombre
}

func agregar(nombre string) string {
	tamanio := len(discos)
	encontrado := false
	ubicacion := 0

	for i := 0; i < tamanio; i++ {
		nombreV := strings.Split(discos[i], "_")
		if strings.ToLower(nombreV[0]) == strings.ToLower(nombre) {
			encontrado = true
			break
		}
		ubicacion = i + 1
	}

	if encontrado == false {
		temp := nombre + "_" + strconv.Itoa(ubicacion) + "_1"
		discos = append(discos, temp)
	}
	//fmt.Println("discos[ubicacion]:", discos[ubicacion])
	return discos[ubicacion]
}

func Obtener_ruta(id string) string {
	rut := ""

	for a := 0; a < len(montados); a++ {
		if montados[a] != "" {
			dividir := strings.Split(montados[a], "_")

			if dividir[2] == strings.ToLower(id) {
				rut = dividir[0]
				break
			}
		}
	}

	return rut
}

func Obtener_particion(id string) estructuras.Nodo_particion {
	path := ""
	name := ""
	encontrado := false

	for a := 0; a < len(montados); a++ {
		if montados[a] != "" {
			dividir := strings.Split(montados[a], "_")

			if dividir[2] == strings.ToLower(id) {
				path = dividir[0]
				name = dividir[1]
				encontrado = true
				break
			}
		}
	}

	if encontrado {
		disco := estructuras.Nodo_Mbr{}

		var tamanio_mbr int = int(unsafe.Sizeof(disco))

		file, err := os.OpenFile(path, os.O_RDWR, 0777)
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}

		data := leerBytes(file, tamanio_mbr)
		buffer := bytes.NewBuffer(data)

		err = binary.Read(buffer, binary.BigEndian, &disco)
		if err != nil {
			log.Fatal("binary.Read failed", err)
		}

		for i := 0; i < 4; i++ {
			if validar_elim(disco, name, i) {
				return disco.Partition[i]
			}
		}
	}
	vacio := estructuras.Nodo_particion{}
	vacio.Part_start = -1
	return vacio
}

func Add_particion(path string, name string) {
	unir := strings.ToLower(path) + ">" + strings.ToLower(name)
	name_particiones = append(name_particiones, unir)
}

func verificar_particion(path string, name string) bool {
	unir := strings.ToLower(path) + ">" + strings.ToLower(name)
	tamanio := len(name_particiones)

	for i := 0; i < tamanio; i++ {
		if name_particiones[i] == unir {
			return true
		}
	}

	return false
}
