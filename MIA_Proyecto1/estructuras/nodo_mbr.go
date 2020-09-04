package estructuras

type Nodo_Mbr struct {
	Mbr_tamanio        int64
	Mbr_fecha_creacion [20]byte
	Mbr_disk_signature int64
	Partition          [4]Nodo_particion
	Cant_Partitions    int8
	Logic_Exist        bool
}
