package modelos

type ResumenAudio struct {
	idAudio int
	titulo  string
	idTipo  int
}

func (r *ResumenAudio) GetIdAudio() int        { return r.idAudio }
func (r *ResumenAudio) SetIdAudio(id int)       { r.idAudio = id }

func (r *ResumenAudio) GetTitulo() string       { return r.titulo }
func (r *ResumenAudio) SetTitulo(titulo string)  { r.titulo = titulo }

func (r *ResumenAudio) GetIdTipo() int          { return r.idTipo }
func (r *ResumenAudio) SetIdTipo(idTipo int)     { r.idTipo = idTipo }
