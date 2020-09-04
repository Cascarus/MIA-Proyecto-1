package interprete

import (
	"MIA_Proyecto1/funciones"
	"fmt"
	"strconv"
	"strings"
)

var path string
var ruta string
var size int64
var name string
var nombre string
var unit byte
var tipo byte
var fit string
var eliminar string
var agregar int64
var IDs string
var opcionFD int8
var errorGeneral bool

type interprete struct {
	comando []string
}

func New(comand string) interprete {
	comansplit := strings.Split(comand, " -")
	coso := interprete{comando: comansplit}
	limpiar()
	return coso
}

func (e interprete) Ejecutar() {
	//imprimir()
	if strings.ToLower(e.comando[0]) == "mkdisk" {
		Mkdisk(e.comando[1:])
	} else if strings.ToLower(e.comando[0]) == "rmdisk" {
		Mrdisk(e.comando[1:])
	} else if strings.ToLower(e.comando[0]) == "fdisk" {
		Fdisk(e.comando[1:])
	} else if strings.ToLower(e.comando[0]) == "rep" {
		Rep(e.comando[1:])
	}
	//imprimir()
	//fmt.Println("\033[1;32mMENSAJE: El disco se ha creado exitosamente!\033[0m")
}

func Mkdisk(contenido []string) {
	unit = 'm'
	//var b_size, b_path, b_name bool

	Opciones_Parametro(contenido)

	if errorGeneral == true {
		return
	}

	if name == "default" {
		funciones.Mensaje("MKDISK debe de llevar un nombre", 2)
		return
	} else if path == "default" {
		funciones.Mensaje("MKDISK debe de llevar un path", 2)
		return
	} else if size < 1 {
		funciones.Mensaje("MKDISK debe de llevar un tamaño", 2)
		return
	} else if unit == 98 {
		funciones.Mensaje("MKDISK debe de usar k(kilobytes) o m(megabytes)", 2)
		return
	}

	verificar := strings.Split(name, ".")
	fmt.Println(name)
	//fmt.Println(verificar[0], " ", verificar[1])
	if len(verificar) == 2 && strings.ToLower(verificar[1]) != "dsk" {
		funciones.Mensaje("1. El parametro name debe de llevar nombre y la extencion .dsk", 2)
		return

	} else if len(verificar) == 2 && strings.ToLower(verificar[0]) == "" {
		funciones.Mensaje("2. El parametro name debe de llevar nombre y la extencion .dsk", 2)
		return

	} else if len(verificar) < 2 {
		funciones.Mensaje("El parametro name debe de llevar nombre y la extencion .dsk", 2)
		return
	}

	//fmt.Println("Se creara un disco con -size=", size, "-path=", path, "-name=", name, "-unit=", unit)
	mkd := funciones.NewMKDisk(size, path, name, unit)
	//mkd.Ejecutar(5846, 21, 3000)
	mkd.Ejecutar()
	//fmt.Println("Reading File: ")
	//mkd.ReadFile()
}

func Mrdisk(contenido []string) {
	var path string

	if len(contenido) == 1 {
		temp := strings.Split(contenido[0], "->")

		if strings.ToLower(temp[0]) == "-path" {
			sin_comillas := strings.Trim(temp[1], "\"")
			//sin_pslice := strings.TrimLeft(sin_comillas, "/")
			path = sin_comillas
		} else {
			funciones.Mensaje("MRDISK solo puede llevar Path", 2)
			return
		}
	} else if len(contenido) == 0 {
		funciones.Mensaje("MRDISK debe llevar un path", 2)
		return

	} else {
		funciones.Mensaje("MRDISK no puede llevar otro parametro que no sea path", 2)
		return
	}

	mrd := funciones.NewMRDisk(path)
	mrd.Ejecutar()

}

func Fdisk(contenido []string) {
	unit = 'k'
	opcionFD = 0
	fit = "wf"
	Opciones_Parametro(contenido)

	if errorGeneral == true {
		return
	}

	if opcionFD == 0 { //crear particion
		if name == "default" {
			funciones.Mensaje("FDISK debe de llevar un nombre", 2)
			return
		} else if path == "default" {
			funciones.Mensaje("FDISK debe de llevar un path", 2)
			return
		} else if size < 1 {
			funciones.Mensaje("FDISK debe de llevar un tamaño", 2)
			return
		}
	} else if opcionFD == 1 { //eliminar
		if name == "default" {
			funciones.Mensaje("FDISK debe de llevar un nombre", 2)
			return
		} else if path == "default" {
			funciones.Mensaje("FDISK debe de llevar un path", 2)
			return
		}
	} else if opcionFD == 2 { //agregar
		if name == "default" {
			funciones.Mensaje("FDISK debe de llevar un nombre", 2)
			return
		} else if path == "default" {
			funciones.Mensaje("FDISK debe de llevar un path", 2)
			return
		}

	}

	//fmt.Println("Se ejecuatara FDISK con -size=", size, "-path=", path, "-name=", name, "-unit=", unit, "-type=", tipo, "-fit=", fit, "-delete=", eliminar, "-add=", agregar)
	fdisk := funciones.NewFDisk(size, unit, path, tipo, fit, eliminar, name, agregar, opcionFD)
	fdisk.Ejecutar()
}

func Rep(contenido []string) {

	Opciones_Parametro(contenido)

	if errorGeneral == true {
		return
	}

	if nombre == "default" {
		funciones.Mensaje("Rep debe de llevar el nombre del reporte que desea generar", 2)
		return
	} else if path == "default" {
		funciones.Mensaje("Rep debe de llevar un path", 2)
		return
	} else if IDs == "default" {
		funciones.Mensaje("Rep debe de llevar el id de una particion", 2)
		return
	}
}

func Opciones_Parametro(contenido []string) {

	for i := 0; i < len(contenido); i++ {
		temp := strings.Split(contenido[i], "->")

		if strings.ToLower(temp[0]) == "size" {
			size_temp, err := strconv.ParseInt(temp[1], 10, 64)
			if err != nil || size_temp <= 0 {
				funciones.Mensaje("Size solo pude utilizar numeros enteros positivos", 2)
				errorGeneral = true
				return
			}
			size = size_temp

		} else if strings.ToLower(temp[0]) == "path" {
			sin_comillas := strings.Trim(temp[1], "\"")
			path = sin_comillas
			//fmt.Println(path)

		} else if strings.ToLower(temp[0]) == "ruta" {
			sin_comillas := strings.Trim(temp[1], "\"")
			ruta = sin_comillas

		} else if strings.ToLower(temp[0]) == "name" {
			sin_comillas := strings.Trim(temp[1], "\"")

			if strings.ToLower(sin_comillas) == "" {
				funciones.Mensaje("Name debe de llevar un nombre", 2)
				errorGeneral = true
				return
			}
			name = sin_comillas

		} else if strings.ToLower(temp[0]) == "nombre" {

			if strings.ToLower(temp[1]) == "mbr" {
				nombre = temp[1]
			} else if strings.ToLower(temp[1]) == "disk" {
				nombre = temp[1]
			} else if strings.ToLower(temp[1]) == "sb" {
				nombre = temp[1]
			} else if strings.ToLower(temp[1]) == "bm_ardir" {
				nombre = temp[1]
			} else if strings.ToLower(temp[1]) == "bm_detdir" {
				nombre = temp[1]
			} else if strings.ToLower(temp[1]) == "bm_inode" {
				nombre = temp[1]
			} else if strings.ToLower(temp[1]) == "bm_block" {
				nombre = temp[1]
			} else if strings.ToLower(temp[1]) == "bitacora" {
				nombre = temp[1]
			} else if strings.ToLower(temp[1]) == "directorio" {
				nombre = temp[1]
			} else if strings.ToLower(temp[1]) == "tree_file" {
				nombre = temp[1]
			} else if strings.ToLower(temp[1]) == "tree_directorio" {
				nombre = temp[1]
			} else if strings.ToLower(temp[1]) == "tree_complete" {
				nombre = temp[1]
			} else if strings.ToLower(temp[1]) == "ls" {
				nombre = temp[1]
			} else {
				funciones.Mensaje("Nombre debe de un valor valido", 2)
				errorGeneral = true
				return
			}

		} else if strings.ToLower(temp[0]) == "unit" {
			if strings.ToLower(temp[1]) == "m" {
				unit = 'm'
			} else if strings.ToLower(temp[1]) == "k" {
				unit = 'k'
			} else if strings.ToLower(temp[1]) == "b" {
				unit = 'b'
			} else {
				funciones.Mensaje("Unit solo puede utilizar m(megabytes), k(kilobytes) o b(bytes)", 2)
				errorGeneral = true
				return
			}
		} else if strings.ToLower(temp[0]) == "type" {
			if strings.ToLower(temp[1]) == "p" {
				tipo = 'p'
			} else if strings.ToLower(temp[1]) == "e" {
				tipo = 'e'
			} else if strings.ToLower(temp[1]) == "l" {
				tipo = 'l'
			} else {
				funciones.Mensaje("Type solo puede utilizar p(primaria), e(extendida) o l(logica)", 2)
				errorGeneral = true
				return
			}
		} else if strings.ToLower(temp[0]) == "fit" {
			if strings.ToLower(temp[1]) == "bf" {
				fit = "bf"
			} else if strings.ToLower(temp[1]) == "ff" {
				fit = "ff"
			} else if strings.ToLower(temp[1]) == "wf" {
				fit = "wf"
			} else {
				funciones.Mensaje("Fit solo puede utilizar bf(best fit), ff(first fit) o wf(worst fit)", 2)
				errorGeneral = true
				return
			}
		} else if strings.ToLower(temp[0]) == "delete" {
			if strings.ToLower(temp[1]) == "fast" {
				eliminar = "fast"
			} else if strings.ToLower(temp[1]) == "full" {
				eliminar = "full"
			} else {
				funciones.Mensaje("Delete solo puede utilizar fast o full", 2)
				errorGeneral = true
				return
			}
			opcionFD = 1
		} else if strings.ToLower(temp[0]) == "add" {
			temp_add, err := strconv.ParseInt(temp[1], 10, 64)
			if err != nil || temp_add == 0 {
				funciones.Mensaje("Add solo puede utilizar un numero diferente de 0", 2)
				errorGeneral = true
				return
			}
			agregar = temp_add
			opcionFD = 2
		} else if strings.ToLower(temp[0]) == "id" {
			if strings.ToLower(temp[1]) == "" {
				funciones.Mensaje("Id debe de llevar un nombre", 2)
				errorGeneral = true
				return
			}
			IDs = temp[1]
		}
	}
}

func limpiar() {
	path = "default" // ya
	ruta = "default"
	size = -1        //ya
	name = "default" //ya
	nombre = "default"
	unit = '0'           //ya
	tipo = 'p'           //ya
	fit = "default"      //ya
	eliminar = "default" //ya
	agregar = -1
	IDs = "default"
	opcionFD = -1
	errorGeneral = false

}

func imprimir() {
	fmt.Println("path: ", path)
	fmt.Println("size: ", size)
	fmt.Println("name: ", name)
	fmt.Println("unit: ", unit)
	fmt.Println("tipo: ", tipo)
	fmt.Println("fit: ", fit)
	fmt.Println("delete: ", eliminar)
	fmt.Println("add: ", agregar)
	fmt.Println("IDs: ", IDs)
	fmt.Println("opcionFD: ", opcionFD)
	fmt.Println("ErrorG: ", errorGeneral)
}
