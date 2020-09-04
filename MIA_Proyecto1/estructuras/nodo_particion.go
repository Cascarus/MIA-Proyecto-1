package estructuras

type Nodo_particion struct {
	Part_status byte
	Part_tipo   byte
	Part_fit    byte
	Part_start  int64
	Part_size   int64
	Part_name   [16]byte
}
