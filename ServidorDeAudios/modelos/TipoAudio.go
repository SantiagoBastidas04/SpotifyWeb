package modelos

type TipoAudio struct {
	id     int
	nombre string
}

func (t *TipoAudio) GetId() int {
	return t.id
}

func (t *TipoAudio) SetId(id int) {
	t.id = id
}

func (t *TipoAudio) GetNombre() string {
	return t.nombre
}

func (t *TipoAudio) SetNombre(nombre string) {
	t.nombre = nombre
}
