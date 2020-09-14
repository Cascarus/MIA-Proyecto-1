package estructuras

type Nodo_AVD struct {
	Avd_fecha_creacion              [19]byte
	Avd_nombre_directorio           [20]byte
	Avd_ap_array_subdirectorios     [6]int64
	Avd_ap_detalle_directorio       int64
	Avd_ap_arbol_virtual_directorio int64
	Avd_proper                      int64
}
