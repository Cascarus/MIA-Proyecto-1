package estructuras

type Nodo_Inodo struct {
	I_count_inodo             int64
	I_size_archivo            int64
	I_count_bloques_asignados int64
	I_array_bloques           [4]int64
	I_ap_indirecto            int64
	I_id_proper               int64
}
