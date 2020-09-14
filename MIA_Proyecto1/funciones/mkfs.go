package funciones

import (
	"MIA_Proyecto1/estructuras"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"unsafe"
)

type mkfs struct {
	id      string
	unit    byte
	tipo    byte
	agregar int64
}

func NewMkfs(id string, unit byte, tipo byte, agregar int64) mkfs {
	e := mkfs{id, unit, tipo, agregar}
	return e
}

func (e mkfs) Ejecutar() {
	path := Obtener_ruta(e.id)
	particion := Obtener_particion(e.id)
	nombre_split := strings.Split(path, "/")
	nombre := nombre_split[len(nombre_split)-1]
	fmt.Println("nombre:", nombre)

	if particion.Part_start < 0 {
		Mensaje("La particion no exite en ese disco", 2)
		return
	}

	if particion.Part_tipo == 112 {
		fmt.Println(particion.Part_start)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			Mensaje("No existe un disco con ese nombre en el directorio", 2)
			return

		}

		sb := crear_SB(particion, nombre)

		avd := crear_AVD("/", 0, 1, 664)

		dd := crear_DD()
		dd.Dd_array_files[0] = crear_Block_File("user.txt", 0)

		usertxt := crear_Inodo(0, 0, 1)
		usertxt.I_array_bloques[0] = 0
		usertxt.I_array_bloques[1] = 1

		dbuser1 := crear_Block_Data("1,G,ROOT\n1,U,ROOT,ROOT,2")
		dbuser2 := crear_Block_Data("01113915\n")

		sb.Sb_inodos_free--
		sb.Sb_arbol_virutal_free--
		sb.Sb_detalle_directorio_free--
		sb.Sb_bloques_free--
		sb.Sb_bloques_free--

		BMAVD := crear_BitMap(sb.Sb_arbol_virtual_count)
		BMDD := crear_BitMap(sb.Sb_detalle_directorio_count)
		BMINODO := crear_BitMap(sb.Sb_inodos_count)
		BMBLOCK := crear_BitMap(sb.Sb_bloques_count)

		BMAVD[0] = "1"
		BMDD[0] = "1"
		BMINODO[0] = "1"
		BMBLOCK[0] = "1"
		BMBLOCK[1] = "1"

		set_super_bloque(path, sb, particion.Part_start)
		set_AVD(path, avd, sb.Sb_ap_arbol_directorio, 0)
		set_DD(path, dd, sb.Sb_ap_detalle_directorio, 0)
		set_Inodo(path, usertxt, sb.Sb_ap_tabla_inodo, 0)
		set_Block_Data(path, dbuser1, sb.Sb_ap_bloques, 0)
		set_Block_Data(path, dbuser2, sb.Sb_ap_bloques, 1)
		//set_Data(path, particion.Part_start, sb, BMAVD, BMDD, BMINODO, BMBLOCK)

		Mensaje("Se ha formateado la particion con exito", 1)
	} else {
		Mensaje("No se puede formatear una particion extendida", 2)
	}

}

func crear_SB(particion estructuras.Nodo_particion, nombre string) estructuras.Nodo_SBoot {
	SB := estructuras.Nodo_SBoot{}
	n := getN(particion.Part_size)
	fmt.Println("N:", n)

	copy(SB.Sb_nombre_hd[:], nombre)
	current_time := time.Now()
	tiempo_disco := string(current_time.Format("2006/01/02 15:04:05"))

	nn := estructuras.Nodo_AVD{}
	o := estructuras.Nodo_DD{}
	p := estructuras.Nodo_Inodo{}
	q := estructuras.Block_Data{}
	tam_AVD := int64(unsafe.Sizeof(nn))
	tam_DD := int64(unsafe.Sizeof(o))
	tam_Inodo := int64(unsafe.Sizeof(p))
	tam_BData := int64(unsafe.Sizeof(q))

	SB.Sb_arbol_virtual_count = n
	SB.Sb_detalle_directorio_count = n
	SB.Sb_inodos_count = 5 * n
	SB.Sb_bloques_count = 20 * n
	SB.Sb_arbol_virutal_free = n
	SB.Sb_detalle_directorio_free = n
	SB.Sb_inodos_free = 5 * n
	SB.Sb_bloques_free = 20 * n
	copy(SB.Sb_date_creacion[:], tiempo_disco)
	copy(SB.Sb_date_ultimo_montaje[:], tiempo_disco)
	SB.Sb_montajes_count = 0
	SB.Sb_ap_bitmap_arbol_directorio = particion.Part_start + int64(unsafe.Sizeof(SB))
	SB.Sb_ap_arbol_directorio = particion.Part_start + int64(unsafe.Sizeof(SB)) + n
	SB.Sb_ap_bitmap_detalle_directorio = particion.Part_start + int64(unsafe.Sizeof(SB)) + n + n*tam_AVD
	SB.Sb_ap_detalle_directorio = particion.Part_start + int64(unsafe.Sizeof(SB)) + n + n*tam_AVD + n
	SB.Sb_ap_bitmap_tabla_inodo = particion.Part_start + int64(unsafe.Sizeof(SB)) + n + n*tam_AVD + n + n*tam_DD
	SB.Sb_ap_tabla_inodo = particion.Part_start + int64(unsafe.Sizeof(SB)) + n + n*tam_AVD + n + n*tam_DD + 5*n
	SB.Sb_ap_bitmap_bloques = particion.Part_start + int64(unsafe.Sizeof(SB)) + n + n*tam_AVD + n + n*tam_DD + 5*n + 5*n*tam_Inodo
	SB.Sb_ap_bloques = particion.Part_start + int64(unsafe.Sizeof(SB)) + n + n*tam_AVD + n + n*tam_DD + 5*n + 5*n*tam_Inodo + 20*n
	SB.Sb_ap_log = particion.Part_start + int64(unsafe.Sizeof(SB)) + n + n*tam_AVD + n + n*tam_DD + 5*n + 5*n*tam_Inodo + 20*n + 20*n*tam_BData
	SB.Sb_size_struct_arbol_directorio = tam_AVD
	SB.Sb_size_struct_detalle_directorio = tam_DD
	SB.Sb_size_struct_inodo = tam_Inodo
	SB.Sb_size_struct_bloque = tam_BData
	SB.Sb_first_free_bit_arbol_directorio = 0
	SB.Sb_first_free_bit_detalle_directorio = 0
	SB.Sb_first_free_bit_tabla_inodo = 0
	SB.Sb_first_free_bit_bloques = 0
	copy(SB.Sb_magic_num[:], "201603127")

	return SB
}

func crear_AVD(name string, det int64, usr int64, perm int64) estructuras.Nodo_AVD {
	AVD := estructuras.Nodo_AVD{}

	copy(AVD.Avd_fecha_creacion[:], getTime())

	for i := 0; i < 6; i++ {
		AVD.Avd_ap_array_subdirectorios[i] = 0
	}

	AVD.Avd_ap_detalle_directorio = det
	AVD.Avd_ap_arbol_virtual_directorio = -1
	AVD.Avd_proper = usr

	return AVD
}

func crear_DD() estructuras.Nodo_DD {
	DD := estructuras.Nodo_DD{}

	for i := 0; i < 5; i++ {
		DD.Dd_array_files[i] = crear_Block_File("", -1)
	}

	DD.Dd_ap_detalle_directorio = -1

	return DD
}

func crear_Block_File(name string, ap_inodo int64) estructuras.Block_File {
	BF := estructuras.Block_File{}

	copy(BF.Dd_file_nombre[:], name)
	BF.Dd_file_ap_inodo = ap_inodo
	copy(BF.Dd_file_date_creation[:], getTime())
	copy(BF.Dd_file_date_modificacion[:], getTime())

	return BF
}

func crear_Inodo(inodo int64, size int64, proper int64) estructuras.Nodo_Inodo {
	INODO := estructuras.Nodo_Inodo{}

	INODO.I_count_inodo = inodo
	INODO.I_size_archivo = size
	INODO.I_count_bloques_asignados = 0

	for i := 0; i < 4; i++ {
		INODO.I_array_bloques[i] = -1
	}

	INODO.I_ap_indirecto = -1
	INODO.I_id_proper = proper

	return INODO
}

func crear_Block_Data(block string) estructuras.Block_Data {
	DATA := estructuras.Block_Data{}

	copy(DATA.Db_data[:], block)

	return DATA
}

func crear_BitMap(n int64) []string {
	var bitmap []string

	for i := 0; i < int(n); i++ {
		bitmap = append(bitmap, "0")
	}

	return bitmap
}

func set_super_bloque(path string, sb estructuras.Nodo_SBoot, pos int64) {

	file, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
		Mensaje("El disco no existe en la Ubicacion", 2)
	}

	file.Seek(pos, 0)
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, &sb)
	escribirBytes(file, binario.Bytes())

	Mensaje("Se inserto el Super Boot con exito", 1)
}

func set_AVD(path string, AVD estructuras.Nodo_AVD, inicio int64, n int64) {

	file, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		Mensaje("El disco no existe en la ubicacion", 2)
		log.Fatal(err)
	}

	file.Seek(inicio+n*int64(unsafe.Sizeof(AVD)), 0)
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, &AVD)
	escribirBytes(file, binario.Bytes())
	Mensaje("Se inserto el AVD con exito", 1)
}

func set_DD(path string, DD estructuras.Nodo_DD, inicio int64, n int64) {

	file, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		Mensaje("El disco no existe en la ubicacion", 2)
		log.Fatal(err)
	}

	file.Seek(inicio+n*int64(unsafe.Sizeof(DD)), 0)
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, &DD)
	escribirBytes(file, binario.Bytes())
	Mensaje("Se inserto el DD con exito", 1)
}

func set_Inodo(path string, INODO estructuras.Nodo_Inodo, inicio int64, n int64) {

	file, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		Mensaje("El disco no existe en la ubicacion", 2)
		log.Fatal(err)
	}

	file.Seek(inicio+n*int64(unsafe.Sizeof(INODO)), 0)
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, &INODO)
	escribirBytes(file, binario.Bytes())
	Mensaje("Se inserto el INODO con exito", 1)
}

func set_Block_Data(path string, BLOCK_DATA estructuras.Block_Data, inicio int64, n int64) {

	file, err := os.OpenFile(path, os.O_RDWR, 0777)
	defer file.Close()
	if err != nil {
		Mensaje("El disco no existe en la ubicacion", 2)
		log.Fatal(err)
	}

	file.Seek(inicio+n*int64(unsafe.Sizeof(BLOCK_DATA)), 0)
	var binario bytes.Buffer
	binary.Write(&binario, binary.BigEndian, &BLOCK_DATA)
	escribirBytes(file, binario.Bytes())
	Mensaje("Se inserto el Block Data con exito", 1)
}

/*
func set_Data(path string, inicio int64, sb estructuras.Nodo_SBoot, BMAVD string, BMDD string, BMINODO string, BMBLOCK string) {
	set_super_bloque(path, sb, inicio)

	if BMAVD != "" {
		set_BitMap(path, BMAVD, sb.Sb_ap_bitmap_arbol_directorio)
		//sb.Sb_first_free_bit_arbol_directorio
	}
}

func set_BitMap(path string, bitmap string, pos int64) {

}
*/
func getN(size int64) int64 {
	m := estructuras.Nodo_SBoot{}
	n := estructuras.Nodo_AVD{}
	o := estructuras.Nodo_DD{}
	p := estructuras.Nodo_Inodo{}
	q := estructuras.Block_Data{}
	r := estructuras.Bitacora{}

	tam_SB := int64(unsafe.Sizeof(m))
	tam_AVD := int64(unsafe.Sizeof(n))
	tam_DD := int64(unsafe.Sizeof(o))
	tam_Inodo := int64(unsafe.Sizeof(p))
	tam_BData := int64(unsafe.Sizeof(q))
	tam_Bitacora := int64(unsafe.Sizeof(r))

	return int64((size - (2 * tam_SB)) / (27 + tam_AVD + tam_DD + (5*tam_Inodo + (20 * tam_BData) + tam_Bitacora)))
}

func getTime() string {
	current_time := time.Now()
	tiempo_disco := string(current_time.Format("2006/01/02 15:04:05"))

	return tiempo_disco
}
