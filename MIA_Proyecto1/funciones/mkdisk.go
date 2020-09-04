package funciones

import (
	"MIA_Proyecto1/estructuras"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
	"unsafe"
)

type mkdisk struct {
	size int64
	path string
	name string
	unit byte
}

func NewMKDisk(size int64, path string, name string, unit byte) mkdisk {
	e := mkdisk{size, path, name, unit}
	return e
}

func (e mkdisk) ReadFile() {
	//Abrimos/creamos un archivo.
	directorio_disco := e.path + e.name
	file, err := os.Open(directorio_disco)
	defer file.Close()
	if err != nil { //validar que no sea nulo.
		log.Fatal(err)
	}

	//Declaramos variable de tipo mbr
	m := estructuras.Nodo_Mbr{}
	//Obtenemos el tamanio del mbr
	var size int = int(unsafe.Sizeof(m))

	//Lee la cantidad de <size> bytes del archivo
	data := leerBytes(file, size)
	//Convierte la data en un buffer,necesario para
	//decodificar binario
	buffer := bytes.NewBuffer(data)

	//Decodificamos y guardamos en la variable m
	err = binary.Read(buffer, binary.BigEndian, &m)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}

	//Se imprimen los valores guardados en el struct
	fmt.Println(m)
	//fmt.Printf("tamanio: %d\nsignature: %d\nfecha: %d", m.Mbr_tamanio, m.Mbr_disk_signature, m.Mbr_fecha_creacion)
	fmt.Println("tamanio: ", m.Mbr_tamanio, " signature: ", m.Mbr_disk_signature)
}

func (e mkdisk) Ejecutar() {
	var tam int64
	kb := 1024
	mb := 1024 * 1024

	if e.unit == 107 {
		tam = e.size * int64(kb)
	} else if e.unit == 109 {
		tam = e.size * int64(mb)
	}

	directorio_disco := e.path + e.name
	//fmt.Println("la ruta es: ", directorio_disco)

	if _, err := os.Stat(directorio_disco); os.IsNotExist(err) {
		Crear_directorio(e.path)
		file, err := os.Create(directorio_disco)
		defer file.Close()
		if err != nil {
			log.Fatal(err)
		}

		var otro int8 = 0

		//fmt.Println(unsafe.Sizeof(otro))
		//Escribimos un 0 en el inicio del archivo.
		var binario bytes.Buffer
		binary.Write(&binario, binary.BigEndian, &otro)
		escribirBytes(file, binario.Bytes())
		//Nos posicionamos en el byte 1023 (primera posicion es 0)
		file.Seek(tam-1, 0) // segundo parametro: 0, 1, 2.     0 -> Inicio, 1-> desde donde esta el puntero, 2 -> Del fin para atras

		//Escribimos un 0 al final del archivo.
		var binario2 bytes.Buffer
		binary.Write(&binario2, binary.BigEndian, &otro)
		escribirBytes(file, binario2.Bytes())

		//----------------------------------------------------------------------- //
		//Escribimos nuestro struct en el inicio del archivo

		file.Seek(0, 0) // nos posicionamos en el inicio del archivo.

		//Asignamos valores a los atributos del struct.
		disco := estructuras.Nodo_Mbr{}
		disco.Mbr_tamanio = tam
		disco.Mbr_disk_signature = int64(rand.Intn(1000))
		disco.Logic_Exist = false
		disco.Cant_Partitions = 0
		current_time := time.Now()
		tiempo_disco := string(current_time.Format("2006/01/02 15:04:05"))

		copy(disco.Mbr_fecha_creacion[:], tiempo_disco)

		//Escribimos struct.
		var binario3 bytes.Buffer
		binary.Write(&binario3, binary.BigEndian, &disco)
		escribirBytes(file, binario3.Bytes())

		Mensaje("El disco se ha creado exitosamente", 1)
	} else {
		Mensaje("Ya existe un disco con ese nombre en el directorio", 2)
		return
	}

}

func Crear_directorio(path string) {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		//fmt.Println("No existe")
		//fmt.Println("el path es: ", path)
		errDir := os.MkdirAll(path, 0777)
		if errDir != nil {
			log.Fatal(err)
			//fmt.Println("No se ha creado!!")
		}
		//fmt.Println("se ha creado!!")

	} else {
		//fmt.Println("Ya existe!!")
	}
}
