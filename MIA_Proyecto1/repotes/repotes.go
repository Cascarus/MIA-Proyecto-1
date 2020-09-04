package reportes

import (
	"MIA_Proyecto1/estructuras"
	"bytes"
	"encoding/binary"
	"log"
	"os"
	"os/exec"
	"strings"
	"unsafe"
)

type reporte struct {
	path   string
	name   string
	inicio string
	fin    string
	EBR    string
	MBR    string
}

func NewReporte(path string, name string) reporte {
	e := reporte{path: path, name: name}
	e.inicio = "digraph D{\n\tgraph [pad=\"0.5\", nodesep=\"0.5\", ranksep=\"2\"];\n\tnode [shape=plain]\n\trankdir=LR;\n"
	e.fin = "\t\t</table>>]\n\n}"
	return e
}

func (e reporte) Generar() {

	if strings.ToLower(e.name) == "disk" {
		Rep_disk(e.inicio, e.path)
	}
}

func Rep_disk(incio string, path string) {
	cuerpo := "\t nodito [shape=record, label=\""
	final := "\"];"

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

	cuerpo += "MBR\n" + string(tamanio_mbr) + " | "

	inicio := int64(tamanio_mbr)
	for i := 0; i < 4; i++ {
		if disco.Partition[i].Part_status == 1 {
			if disco.Partition[i].Part_start > inicio {
				disponible := disco.Partition[i].Part_start - inicio
				cuerpo += "Libre\n" + string(disponible) + "|"
				cuerpo += "Primaria =\n" + string(disco.Partition[i].Part_name[:len(disco.Partition[i].Part_name)]) + "\n" + string(disco.Partition[i].Part_size) + "|"
				inicio += disco.Partition[i].Part_start + disco.Partition[i].Part_size

			} else {
				inicio += disco.Partition[i].Part_size
				//Primaria
				if disco.Partition[i].Part_tipo == 112 {
					cuerpo += "Primaria =\n" + string(disco.Partition[i].Part_name[:len(disco.Partition[i].Part_name)]) + "\n" + string(disco.Partition[i].Part_size) + "|"
				} else if disco.Partition[i].Part_tipo == 101 {
					contenido_logica := ""
					cuerpo += "{logica|{" + contenido_logica + "}}|"
				}
			}
		}
	}

	if disco.Mbr_tamanio > inicio {
		disponible := disco.Mbr_tamanio - inicio
		cuerpo += "Libre\n" + string(disponible)
	}

	f, err := os.Create("Reporte_disk.dot")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	f.WriteString(incio)
	f.WriteString(cuerpo)
	f.WriteString(final)
	f.Sync()

	term := exec.Command("dot", "-Tpng", "Reporte_disk.dot", "-o", path)
	err2 := term.Run()
	if err2 != nil {
		log.Fatal(err2)
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
