package estructuras

type Nodo_SBoot struct {
	Sb_nombre_hd                         [16]byte
	Sb_arbol_virtual_count               int64
	Sb_detalle_directorio_count          int64
	Sb_inodos_count                      int64
	Sb_bloques_count                     int64
	Sb_arbol_virutal_free                int64
	Sb_detalle_directorio_free           int64
	Sb_inodos_free                       int64
	Sb_bloques_free                      int64
	Sb_date_creacion                     [20]byte
	Sb_date_ultimo_montaje               [20]byte
	Sb_montajes_count                    int64
	Sb_ap_bitmap_arbol_directorio        int64
	Sb_ap_arbol_directorio               int64
	Sb_ap_bitmap_detalle_directorio      int64
	Sb_ap_detalle_directorio             int64
	Sb_ap_bitmap_tabla_inodo             int64
	Sb_ap_tabla_inodo                    int64
	Sb_ap_bitmap_bloques                 int64
	Sb_ap_bloques                        int64
	Sb_ap_log                            int64
	Sb_size_struct_arbol_directorio      int64
	Sb_size_struct_detalle_directorio    int64
	Sb_size_struct_inodo                 int64
	Sb_size_struct_bloque                int64
	Sb_first_free_bit_arbol_directorio   int64
	Sb_first_free_bit_detalle_directorio int64
	Sb_first_free_bit_tabla_inodo        int64
	Sb_first_free_bit_bloques            int64
	Sb_magic_num                         [9]byte
}
