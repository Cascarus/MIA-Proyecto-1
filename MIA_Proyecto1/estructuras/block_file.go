package estructuras

type Block_File struct {
	Dd_file_nombre            [20]byte
	Dd_file_ap_inodo          int64
	Dd_file_date_creation     [19]byte
	Dd_file_date_modificacion [19]byte
}
