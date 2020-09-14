package estructuras

type Nodo_particion struct {
	Part_status byte
	Part_tipo   byte
	Part_fit    byte
	Part_start  int64
	Part_size   int64
	Part_name   [16]byte
}

func Particion_T0() Nodo_particion {
	e := Nodo_particion{Part_status: '0', Part_tipo: '0', Part_fit: '0'}
	str := "null"
	copy(e.Part_name[:], str)
	return e
}
