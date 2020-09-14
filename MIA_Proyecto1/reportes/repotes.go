package reportes

import (
	"MIA_Proyecto1/estructuras"
	"MIA_Proyecto1/funciones"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

var conta_logica int

type reporte struct {
	path   string
	name   string
	inicio string
	fin    string
	EBR    string
	MBR    string
	ID     string
}

func NewReporte(path string, name string, id string) reporte {
	e := reporte{path: path, name: name, ID: id}
	e.inicio = "digraph D{\n\tgraph [pad=\"0.5\", nodesep=\"0.5\", ranksep=\"2\"];\n\tnode [shape=plain]\n\t //rankdir=LR;\n"
	e.fin = "\t\t</table>>]\n\n}"
	conta_logica = 0
	return e
}

func (e reporte) Generar() {

	if strings.ToLower(e.name) == "disk" {
		Rep_disk(e.path, e.ID)
	} else if strings.ToLower(e.name) == "mbr" {
		Rep_MBR(e.path, e.ID)
	} else if strings.ToLower(e.name) == "sb" {
		Rep_SB(e.path, e.ID)
	}
}

func Rep_disk(path string, id string) {
	incio := "digraph D{\n\tgraph [pad=\"0.5\", nodesep=\"0.5\", ranksep=\"2\"];\n\tnode [shape=plain]\n\t //rankdir=LR;\n"
	cuerpo := "\t nodito [shape=record, label=\""
	final := "\"];\n}"
	ruta := funciones.Obtener_ruta(id)
	carpeta := get_name_carpeta(path)

	disco := estructuras.Nodo_Mbr{}
	var tamanio_mbr int = int(unsafe.Sizeof(disco))

	if _, err := os.Stat(carpeta); os.IsNotExist(err) {
		funciones.Crear_directorio(carpeta)
	}

	file, err := os.OpenFile(ruta, os.O_RDWR, 0777)
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

	cuerpo += "MBR\\n " + strconv.Itoa(tamanio_mbr) + " | "

	inicio := int64(tamanio_mbr)

	for i := 0; i < 4; i++ {
		if disco.Partition[i].Part_status == 1 {

			if disco.Partition[i].Part_start > inicio {
				disponible := disco.Partition[i].Part_start - inicio
				cuerpo += "Libre\\n" + strconv.FormatInt(disponible, 10) + "|"
				cuerpo += "Primaria\\n" + obtener_nombre(disco.Partition[i].Part_name) + "\\n" + strconv.FormatInt(disco.Partition[i].Part_size, 10) + "|"
				inicio = disco.Partition[i].Part_start + disco.Partition[i].Part_size
			} else {
				inicio += disco.Partition[i].Part_size
				//Primaria
				if disco.Partition[i].Part_tipo == 112 {
					cuerpo += "Primaria\\n" + obtener_nombre(disco.Partition[i].Part_name) + "\\n" + strconv.FormatInt(disco.Partition[i].Part_size, 10) + "|"
				} else if disco.Partition[i].Part_tipo == 101 {
					contenido_logica := obtener_logicas(disco.Partition[i].Part_start, ruta)
					cuerpo += "{Extendida\\n" + strconv.FormatInt(disco.Partition[i].Part_size, 10) + "|{" + contenido_logica + "}}|"
				}
			}
		}
	}
	if disco.Mbr_tamanio > inicio {
		disponible := disco.Mbr_tamanio - inicio
		cuerpo += "Libre\\n" + strconv.FormatInt(disponible, 10)
	}

	path_split := strings.Split(path, ".")
	path_png := path_split[0] + ".png"

	f, err := os.Create("Reporte_disk.dot")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	f.WriteString(incio)
	f.WriteString(cuerpo)
	f.WriteString(final)
	f.Sync()

	fmt.Print("Creando Reporte")

	for i := 0; i < 3; i++ {
		fmt.Print(".")
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Println("")

	term := exec.Command("dot", "-Tpng", "Reporte_disk.dot", "-o", path_png)
	err2 := term.Run()
	if err2 != nil {
		log.Fatal(err2)
	}
	funciones.Mensaje("Se ha creado el reporte disk con exito", 1)
}

func Rep_MBR(path string, id string) {
	inicio := "digraph D{\n\tgraph [pad=\"0.5\", nodesep=\"0.5\", ranksep=\"2\"];\n\tnode [shape=plain]\n\t rankdir=LR;\n"
	fin := "\t\t</table>>]\n\n}"

	MBR := "\tnodo [label=<\n \t\t<table>\n\t\t<tr><td>Nombre</td><td>Valor</td></tr>\n"
	EBR := ""

	ruta := funciones.Obtener_ruta(id)
	carpeta := get_name_carpeta(path)

	disco := estructuras.Nodo_Mbr{}
	var tamanio_mbr int = int(unsafe.Sizeof(disco))

	if _, err := os.Stat(carpeta); os.IsNotExist(err) {
		funciones.Crear_directorio(carpeta)
	}

	file, err := os.OpenFile(ruta, os.O_RDWR, 0777)
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

	MBR += "\t\t<tr><td>mbr_tamaño</td><td>" + strconv.FormatInt(disco.Mbr_tamanio, 10) + "</td></tr>\n"
	MBR += "\t\t<tr><td>mbr_fecha_creacion</td><td>" + string(disco.Mbr_fecha_creacion[:19]) + "</td></tr>\n"
	MBR += "\t\t<tr><td>mbr_disk_signature</td><td>" + strconv.FormatInt(disco.Mbr_disk_signature, 10) + "</td></tr>\n"

	for i := 0; i < 4; i++ {
		status := int(disco.Partition[i].Part_status)
		tipo := string(disco.Partition[i].Part_tipo)
		fit := string(disco.Partition[i].Part_fit)
		start := strconv.FormatInt(disco.Partition[i].Part_start, 10)
		tamanio := strconv.FormatInt(disco.Partition[i].Part_size, 10)
		nombre := obtener_nombre(disco.Partition[i].Part_name)

		MBR += "\t\t<tr><td>part_status_" + strconv.Itoa(i+1) + "</td><td>" + strconv.Itoa(status) + "</td></tr>\n"
		MBR += "\t\t<tr><td>part_type_" + strconv.Itoa(i+1) + "</td><td>" + tipo + "</td></tr>\n"
		MBR += "\t\t<tr><td>part_fit_" + strconv.Itoa(i+1) + "</td><td>" + fit + "</td></tr>\n"
		MBR += "\t\t<tr><td>part_start_" + strconv.Itoa(i+1) + "</td><td>" + start + "</td></tr>\n"
		MBR += "\t\t<tr><td>part_size_" + strconv.Itoa(i+1) + "</td><td>" + tamanio + "</td></tr>\n"
		MBR += "\t\t<tr><td>part_name_" + strconv.Itoa(i+1) + "</td><td>" + nombre + "</td></tr>\n"

		if disco.Partition[i].Part_tipo == 101 {
			EBR = MBR_logicas(ruta, disco.Partition[i].Part_start)
		}
	}

	path_split := strings.Split(path, ".")
	path_png := path_split[0] + ".png"

	f, err := os.Create("Reporte_mbr.dot")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	f.WriteString(inicio)
	f.WriteString(EBR)
	f.WriteString(MBR)
	f.WriteString(fin)
	f.Sync()

	fmt.Print("Creando Reporte")

	for i := 0; i < 3; i++ {
		fmt.Print(".")
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Println("")

	term := exec.Command("dot", "-Tpng", "Reporte_mbr.dot", "-o", path_png)
	err2 := term.Run()
	if err2 != nil {
		log.Fatal(err2)
	}
	funciones.Mensaje("Se ha creado el reporte mbr con exito", 1)

}

func Rep_SB(path string, id string) {
	inicio := "digraph D{\n\tgraph [pad=\"0.5\", nodesep=\"0.5\", ranksep=\"2\"];\n\tnode [shape=plain]\n\t rankdir=LR;\n"
	fin := "\t\t</table>>]\n\n}"
	SB := "\tnodo [label=<\n \t\t<table>\n\t\t<tr><td>Nombre</td><td>Valor</td></tr>\n"

	ruta := funciones.Obtener_ruta(id)
	carpeta := get_name_carpeta(path)
	particion := funciones.Obtener_particion(id)

	superB := estructuras.Nodo_SBoot{}
	var tamanio_mbr int = int(unsafe.Sizeof(superB))

	if _, err := os.Stat(carpeta); os.IsNotExist(err) {
		funciones.Crear_directorio(carpeta)
	}

	file, err := os.OpenFile(ruta, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	file.Seek(particion.Part_start, 0)

	data := leerBytes(file, tamanio_mbr)
	buffer := bytes.NewBuffer(data)

	err = binary.Read(buffer, binary.BigEndian, &superB)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}

	SB += "\t\t<tr><td>sb_nombre_hd</td><td>" + obtener_nombre(superB.Sb_nombre_hd) + "</td></tr>\n"
	SB += "\t\t<tr><td>sb_arbol_virtual_count</td><td>" + strconv.FormatInt(superB.Sb_arbol_virtual_count, 10) + "</td></tr>\n"
	SB += "\t\t<tr><td>sb_detalle_directorio_count</td><td>" + strconv.FormatInt(superB.Sb_detalle_directorio_count, 10) + "</td></tr>\n"
	SB += "\t\t<tr><td>sb_inodos_count</td><td>" + strconv.FormatInt(superB.Sb_inodos_count, 10) + "</td></tr>\n"
	SB += "\t\t<tr><td>sb_bloques_count</td><td>" + strconv.FormatInt(superB.Sb_bloques_count, 10) + "</td></tr>\n"
	SB += "\t\t<tr><td>sb_arbol_virutal_free</td><td>" + strconv.FormatInt(superB.Sb_arbol_virutal_free, 10) + "</td></tr>\n"
	SB += "\t\t<tr><td>sb_detalle_directorio_free</td><td>" + strconv.FormatInt(superB.Sb_detalle_directorio_free, 10) + "</td></tr>\n"
	SB += "\t\t<tr><td>sb_inodos_free</td><td>" + strconv.FormatInt(superB.Sb_inodos_free, 10) + "</td></tr>\n"
	SB += "\t\t<tr><td>sb_bloques_free</td><td>" + strconv.FormatInt(superB.Sb_bloques_free, 10) + "</td></tr>\n"
	SB += "\t\t<tr><td>sb_date_creacion</td><td>" + string(superB.Sb_date_creacion[:19]) + "</td></tr>\n"
	SB += "\t\t<tr><td>sb_date_ultimo_montaje</td><td>" + string(superB.Sb_date_ultimo_montaje[:19]) + "</td></tr>\n"
	SB += "\t\t<tr><td>sb_montajes_count</td><td>" + strconv.FormatInt(superB.Sb_montajes_count, 10) + "</td></tr>\n"
	SB += "\t\t<tr><td>sb_ap_bitmap_arbol_directorio</td><td>" + strconv.FormatInt(superB.Sb_ap_bitmap_arbol_directorio, 10) + "</td></tr>\n"
	SB += "\t\t<tr><td>sb_ap_arbol_directorio</td><td>" + strconv.FormatInt(superB.Sb_ap_arbol_directorio, 10) + "</td></tr>\n"
	SB += "\t\t<tr><td>sb_ap_bitmap_detalle_directorio</td><td>" + strconv.FormatInt(superB.Sb_ap_bitmap_detalle_directorio, 10) + "</td></tr>\n"
	SB += "\t\t<tr><td>sb_ap_detalle_directorio</td><td>" + strconv.FormatInt(superB.Sb_ap_detalle_directorio, 10) + "</td></tr>\n"
	SB += "\t\t<tr><td>sb_ap_bitmap_tabla_inodo</td><td>" + strconv.FormatInt(superB.Sb_ap_bitmap_tabla_inodo, 10) + "</td></tr>\n"
	SB += "\t\t<tr><td>sb_ap_tabla_inodo</td><td>" + strconv.FormatInt(superB.Sb_ap_tabla_inodo, 10) + "</td></tr>\n"
	SB += "\t\t<tr><td>sb_ap_bitmap_bloques</td><td>" + strconv.FormatInt(superB.Sb_ap_bitmap_bloques, 10) + "</td></tr>\n"
	SB += "\t\t<tr><td>sb_ap_bloques</td><td>" + strconv.FormatInt(superB.Sb_ap_bloques, 10) + "</td></tr>\n"
	SB += "\t\t<tr><td>sb_ap_log</td><td>" + strconv.FormatInt(superB.Sb_ap_log, 10) + "</td></tr>\n"
	SB += "\t\t<tr><td>sb_size_struct_arbol_directorio</td><td>" + strconv.FormatInt(superB.Sb_size_struct_arbol_directorio, 10) + "</td></tr>\n"
	SB += "\t\t<tr><td>sb_size_struct_detalle_directorio</td><td>" + strconv.FormatInt(superB.Sb_size_struct_detalle_directorio, 10) + "</td></tr>\n"
	SB += "\t\t<tr><td>sb_size_struct_inodo</td><td>" + strconv.FormatInt(superB.Sb_size_struct_inodo, 10) + "</td></tr>\n"
	SB += "\t\t<tr><td>sb_size_struct_bloque</td><td>" + strconv.FormatInt(superB.Sb_size_struct_bloque, 10) + "</td></tr>\n"
	SB += "\t\t<tr><td>sb_first_free_bit_arbol_directorio</td><td>" + strconv.FormatInt(superB.Sb_first_free_bit_arbol_directorio, 10) + "</td></tr>\n"
	SB += "\t\t<tr><td>sb_first_free_bit_detalle_directorio</td><td>" + strconv.FormatInt(superB.Sb_first_free_bit_detalle_directorio, 10) + "</td></tr>\n"
	SB += "\t\t<tr><td>sb_first_free_bit_tabla_inodo</td><td>" + strconv.FormatInt(superB.Sb_first_free_bit_tabla_inodo, 10) + "</td></tr>\n"
	SB += "\t\t<tr><td>sb_first_free_bit_bloques</td><td>" + strconv.FormatInt(superB.Sb_first_free_bit_bloques, 10) + "</td></tr>\n"
	SB += "\t\t<tr><td>sb_magic_num</td><td>" + string(superB.Sb_magic_num[:9]) + "</td></tr>\n"

	path_split := strings.Split(path, ".")
	path_png := path_split[0] + ".png"

	f, err := os.Create("Reporte_sb.dot")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	f.WriteString(inicio)
	f.WriteString(SB)
	f.WriteString(fin)
	f.Sync()

	fmt.Print("Creando Reporte")

	for i := 0; i < 3; i++ {
		fmt.Print(".")
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Println("")

	term := exec.Command("dot", "-Tpng", "Reporte_sb.dot", "-o", path_png)
	err2 := term.Run()
	if err2 != nil {
		log.Fatal(err2)
	}
	funciones.Mensaje("Se ha creado el reporte sb con exito", 1)

}

func MBR_logicas(path string, inicio int64) string {
	ebr_actual := estructuras.Nodo_ebr{}
	var tamanio_ebr int = int(unsafe.Sizeof(ebr_actual))
	EBR := ""
	conta_logica++

	file, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	file.Seek(inicio, 0)
	data := leerBytes(file, tamanio_ebr)
	buffer := bytes.NewBuffer(data)

	err = binary.Read(buffer, binary.BigEndian, &ebr_actual)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}

	no_logica := conta_logica

	temp := ""
	if ebr_actual.Part_next > 0 {
		temp = MBR_logicas(path, ebr_actual.Part_next)
	}

	status := int(ebr_actual.Part_status)
	fit := string(ebr_actual.Part_fit)
	start := strconv.FormatInt(ebr_actual.Part_start, 10)
	tamanio := strconv.FormatInt(ebr_actual.Part_size, 10)
	sig := strconv.FormatInt(ebr_actual.Part_next, 10)
	nombre := obtener_nombre(ebr_actual.Part_name)

	EBR += temp
	EBR += "\tnodoL" + strconv.Itoa(no_logica) + "[label=<\n \t\t<table>\n"
	EBR += "\t\t<tr><td>EBR_" + strconv.Itoa(no_logica) + "</td></tr>\n"
	EBR += "\t\t<tr><td>Nombre</td><td>Valor</td></tr>\n"
	EBR += "\t\t<tr><td>part_status_" + strconv.Itoa(no_logica) + "</td><td>" + strconv.Itoa(status) + "</td></tr>\n"
	EBR += "\t\t<tr><td>part_fit_" + strconv.Itoa(no_logica) + "</td><td>" + fit + "</td></tr>\n"
	EBR += "\t\t<tr><td>part_start_" + strconv.Itoa(no_logica) + "</td><td>" + start + "</td></tr>\n"
	EBR += "\t\t<tr><td>part_size_" + strconv.Itoa(no_logica) + "</td><td>" + tamanio + "</td></tr>\n"
	EBR += "\t\t<tr><td>part_next_" + strconv.Itoa(no_logica) + "</td><td>" + sig + "</td></tr>\n"
	EBR += "\t\t<tr><td>part_name_" + strconv.Itoa(no_logica) + "</td><td>" + nombre + "</td></tr>\n"
	EBR += "\t\t</table>>]\n\n"

	return EBR

}

func leerBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number) //array de bytes

	_, err := file.Read(bytes) // Leido -> bytes
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}

func obtener_nombre(arreglo [16]byte) string {
	var nombre string

	for i := 0; i < 16; i++ {
		if arreglo[i] != 0 {
			nombre += string(arreglo[i])
		}
	}

	return nombre
}

func obtener_logicas(inicio int64, path string) string {
	ebr := estructuras.Nodo_ebr{}
	var tamanio_ebr int = int(unsafe.Sizeof(ebr))

	file, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	file.Seek(inicio, 0)
	data := leerBytes(file, tamanio_ebr)
	buffer := bytes.NewBuffer(data)

	err = binary.Read(buffer, binary.BigEndian, &ebr)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}

	aux := ""
	if ebr.Part_next > 0 {
		aux = obtener_logicas(ebr.Part_next, path)
	}

	if ebr.Part_size < 1 {
		return ""
	}

	contenido := "EBR|Logica\\n" + strconv.FormatInt(ebr.Part_size, 10)
	return contenido + "|" + aux
}

func get_name_carpeta(path string) string {
	ruta := "/"

	path_s := strings.Split(path, "/")

	for i := 0; i < (len(path_s) - 1); i++ {
		if path_s[i] != "" {
			ruta += path_s[i] + "/"
		}
	}

	return ruta
}
