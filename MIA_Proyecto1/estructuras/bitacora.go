package estructuras

type Bitacora struct {
	log_tipo_operacion int64
	log_tipo           byte
	log_nombre         [256]byte
	log_contenido      [256]byte
	log_fecha          [19]byte
}
