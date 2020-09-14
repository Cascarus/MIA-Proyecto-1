package funciones

import (
	"MIA_Proyecto1/estructuras"
	"bytes"
	"encoding/binary"
	"log"
	"os"
	"strings"
	"unsafe"
)

type fdisk struct {
	size     int64
	unit     byte
	path     string
	tipo     byte
	fit      string
	eliminar string
	name     string
	agregar  int64
	opcionFD int8
}

func NewFDisk(size int64, unit byte, path string, tipo byte, fit string, eliminar string, name string, agregar int64, opcionFD int8) fdisk {
	e := fdisk{size, unit, path, tipo, fit, eliminar, name, agregar, opcionFD}
	return e
}

func (e fdisk) Ejecutar() {
	if _, err := os.Stat(e.path); os.IsNotExist(err) {
		Mensaje("No existe un disco con ese nombre en el directorio", 2)
		return

	} else {

		if e.opcionFD == 0 {
			agregar_particion(e.size, e.unit, e.path, e.tipo, e.fit, e.name)
		} else if e.opcionFD == 1 {
			eliminar_particion(e.eliminar, e.path, e.name)
		}
	}
}

func agregar_particion(size int64, unit byte, path string, tipo byte, fit string, name string) {
	disco := estructuras.Nodo_Mbr{}

	var tamanio_mbr int = int(unsafe.Sizeof(disco))
	////fmt.Println("el tam de mbr es: ", tamanio_mbr)

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

	particion_nueva := estructuras.Nodo_particion{}
	particion_nueva.Part_status = 1
	particion_nueva.Part_size = asignar_size(unit, size)
	particion_nueva.Part_tipo = tipo
	particion_nueva.Part_fit = asignar_fit(fit)
	particion_nueva.Part_start = tipo_ajuste(disco, particion_nueva)

	if particion_nueva.Part_start == 0 {
		Mensaje("57->No hay suficiente espacio en el disco para crear la particion", 2)
		return
	}

	copy(particion_nueva.Part_name[:], name)

	insertado := false
	//switch para cada tipo de particion
	//Primaria
	if tipo == 112 {

		if espacio_disponible(disco)-particion_nueva.Part_size < 0 {
			Mensaje("87->No hay suficiente espacio en el disco para crear la particion", 2)
			return
		}

		if validar_name(disco, name) {
			Mensaje("Ya existe una particion con ese nombre", 2)
			return
		}

		for i := 0; i < 4; i++ {
			if disco.Partition[i].Part_status == 0 {
				disco.Partition[i] = particion_nueva
				insertado = true
				break
			}
		}

		if insertado == false {
			Mensaje("No se puede crear mas de 4 particiones primarias en el disco", 2)
			return

		} else {
			file.Seek(0, 0)
			var binario bytes.Buffer
			binary.Write(&binario, binary.BigEndian, &disco)
			escribirBytes(file, binario.Bytes())
			Add_particion(path, name)
			Mensaje("Se ha creado la particion con exito", 1)
		}

		//Extendida
	} else if tipo == 101 {

		if espacio_disponible(disco)-particion_nueva.Part_size < 0 {
			Mensaje("No hay suficiente espacio en el disco para crear la particion", 2)
			return
		}

		if validar_name(disco, name) {
			Mensaje("Ya existe una particion con ese nombre", 2)
			return
		}

		if disco.Logic_Exist {
			Mensaje("Ya existe una particion extendida en el disco", 2)
			return
		}

		ebr := estructuras.Nodo_particion{}

		for i := 0; i < 4; i++ {
			if disco.Partition[i].Part_status == 0 {
				disco.Partition[i] = particion_nueva
				ebr.Part_start = particion_nueva.Part_start
				disco.Logic_Exist = true
				insertado = true
				break
			}
		}

		if insertado == false {
			Mensaje("No se puede crear mas de 4 particiones primarias en el disco", 2)
			return

		} else {
			file.Seek(0, 0)
			var binario bytes.Buffer
			binary.Write(&binario, binary.BigEndian, &disco)
			escribirBytes(file, binario.Bytes())

			file.Seek(particion_nueva.Part_start, 0)
			var binario2 bytes.Buffer
			binary.Write(&binario2, binary.BigEndian, &ebr)
			escribirBytes(file, binario2.Bytes())
			Add_particion(path, name)
			Mensaje("Se ha creado una particion extendida con exito", 1)
		}
		//Logica
	} else if tipo == 108 {

		if disco.Logic_Exist == false {
			Mensaje("No se puede crear una particion logica porque no existe una extendida", 2)
			return
		}

		nueva_logica := estructuras.Nodo_ebr{}
		nueva_logica.Part_status = 1
		nueva_logica.Part_fit = asignar_fit(fit)
		nueva_logica.Part_size = asignar_size(unit, size)
		copy(nueva_logica.Part_name[:], name)

		EBR_inicial := estructuras.Nodo_ebr{}
		ubicacion_inicial := buscar_primer_EBR(disco)
		tam_extendida := buscar_tam_extendida(disco)

		file.Seek(ubicacion_inicial, 0)
		data := leerBytes(file, int(unsafe.Sizeof(EBR_inicial)))
		buffer := bytes.NewBuffer(data)

		err = binary.Read(buffer, binary.BigEndian, &EBR_inicial)
		if err != nil {
			log.Fatal("binary.Read failed", err)
		}

		var tam int64 = 0
		error := false
		for {
			temp := EBR_inicial
			//fmt.Println("tam extenda: ", tam_extendida)
			//fmt.Println("nueva logica: ", nueva_logica.Part_size)
			if tam+nueva_logica.Part_size <= tam_extendida {
				//fmt.Println("entro al if 1")
				//fmt.Println("part next: ", temp.Part_next)
				if temp.Part_next > 0 {
					//fmt.Println("entro al if 2")
					ubicacion_inicial += temp.Part_size
					tam += temp.Part_size

					file.Seek(ubicacion_inicial, 0)
					data1 := leerBytes(file, int(unsafe.Sizeof(EBR_inicial)))
					buffer2 := bytes.NewBuffer(data1)

					err = binary.Read(buffer2, binary.BigEndian, &EBR_inicial)
					if err != nil {
						log.Fatal("binary.Read failed", err)
					}

				} else {
					//fmt.Println("entro al else 2")
					if temp.Part_size < 1 {
						//fmt.Println("entro al if 3")
						nueva_logica.Part_start = ubicacion_inicial
						break
					} else {
						//fmt.Println("entro al else 3")
						//fmt.Println("ubicacion inicial: ", ubicacion_inicial)
						//fmt.Println("part size: ", temp.Part_size)
						temp.Part_next = ubicacion_inicial + temp.Part_size
						//fmt.Println("part next: ", temp.Part_next)

						file.Seek(ubicacion_inicial, 0)
						var binario33 bytes.Buffer
						binary.Write(&binario33, binary.BigEndian, &temp)
						escribirBytes(file, binario33.Bytes())

						file.Seek(ubicacion_inicial, 0)
						data := leerBytes(file, int(unsafe.Sizeof(EBR_inicial)))
						buffer := bytes.NewBuffer(data)

						err = binary.Read(buffer, binary.BigEndian, &EBR_inicial)
						if err != nil {
							log.Fatal("binary.Read failed", err)
						}

					}
				}
			} else {
				//fmt.Println("entro al else 1")
				error = true
				break
			}
		}

		if error {
			Mensaje("No hay suficionete espacio en el disco para crear la particion", 2)
			return
		} else {

			file.Seek(ubicacion_inicial, 0)
			var binario10 bytes.Buffer
			binary.Write(&binario10, binary.BigEndian, &nueva_logica)
			escribirBytes(file, binario10.Bytes())

			Mensaje("Se ha creado una particion logica con exito", 1)
		}
	}

}

func eliminar_particion(delete string, path string, name string) {
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

	eliminado := false
	if delete == "fast" {
		for i := 0; i < 4; i++ {
			if validar_elim(disco, name, i) {
				disco.Partition[i].Part_status = 0
				eliminado = true
				break
			}
		}
	} else if delete == "full" {
		for i := 0; i < 4; i++ {
			if validar_elim(disco, name, i) {
				disco.Partition[i] = estructuras.Particion_T0()
				eliminado = true
				break
			}
		}
	}

	if eliminado == false {
		Mensaje("El nombre de la particion no existe", 2)
		return
	}

	file.Seek(0, 0)
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, &disco)
	escribirBytes(file, binario.Bytes())

	Mensaje("Se ha eliminado una particion con exito", 1)

}

func espacio_disponible(disco estructuras.Nodo_Mbr) int64 {
	espacio := disco.Mbr_tamanio - int64(unsafe.Sizeof(disco))

	for i := 0; i < 1; i++ {
		if disco.Partition[i].Part_status == 1 {
			espacio -= disco.Partition[i].Part_size
		}
	}
	return espacio
}

func asignar_size(unit byte, tamanio int64) int64 {
	if unit == 98 {
		return tamanio
	} else if unit == 107 {
		return tamanio * 1024
	} else {
		return tamanio * 1024 * 1024
	}
}

func asignar_fit(fit string) byte {
	if fit == "bf" {
		return 'b'
	} else if fit == "ff" {
		return 'f'
	} else {
		return 'w'
	}
}

func buscar_primer_EBR(disco estructuras.Nodo_Mbr) int64 {
	tam_mbr := int64(unsafe.Sizeof(disco))

	for i := 0; i < 4; i++ {
		if disco.Partition[i].Part_tipo == 101 {
			return tam_mbr
		}
		tam_mbr += disco.Partition[i].Part_size
	}

	return 0
}

func buscar_tam_extendida(disco estructuras.Nodo_Mbr) int64 {
	for i := 0; i < 4; i++ {
		////fmt.Println("tipo particion: ", disco.Partition[i].Part_tipo)
		if disco.Partition[i].Part_tipo == 101 {
			return disco.Partition[i].Part_size
		}
	}
	return 0
}

func validar_name(disco estructuras.Nodo_Mbr, name string) bool {
	for i := 0; i < 4; i++ {

		if disco.Partition[i].Part_status == 0 {
			break
		}

		var parti string
		if len(name) <= 16 {
			parti = string(disco.Partition[i].Part_name[:len(name)])
		} else {
			parti = string(disco.Partition[i].Part_name[:16])
		}

		if strings.ToLower(parti) == strings.ToLower(name) {
			return true
		}
	}

	return false
}

func validar_elim(disco estructuras.Nodo_Mbr, name string, pos int) bool {
	var parti string
	if len(name) <= 16 {
		parti = string(disco.Partition[pos].Part_name[:len(name)])
	} else {
		parti = string(disco.Partition[pos].Part_name[:16])
	}

	if strings.ToLower(parti) == strings.ToLower(name) {
		return true
	}
	return false
}

func tipo_ajuste(disco estructuras.Nodo_Mbr, particion estructuras.Nodo_particion) int64 {
	var inicio int64 = int64(unsafe.Sizeof(disco))
	encontrado := false

	for i := 0; i < 4; i++ {
		if disco.Partition[i].Part_status == 1 {
			if disco.Partition[i].Part_start > inicio {
				disponible := disco.Partition[i].Part_start - inicio

				if disponible >= particion.Part_size {
					encontrado = true
					return inicio
				} else {
					inicio = disco.Partition[i].Part_start + disco.Partition[i].Part_size
				}
			} else if disco.Partition[i].Part_start == inicio {
				inicio += disco.Partition[i].Part_size
			}
		}
	}

	if !encontrado {
		if disco.Mbr_tamanio > inicio {
			disponible := disco.Mbr_tamanio - inicio
			if disponible >= particion.Part_size {
				return inicio
			} else {
				return 0
			}
		} else {
			return 0
		}
	}
	return 0
}
