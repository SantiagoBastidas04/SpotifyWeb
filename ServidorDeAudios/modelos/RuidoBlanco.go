package modelos

type RuidoBlanco struct {
	idAudio             int
	tipoSonido          string
	fuenteAudio         string
	usoSugerido         string
	proveedorContenido  string
	duracionBucle       int
	frecuenciaDominante string
}

func (r *RuidoBlanco) GetIdAudio() int                    { return r.idAudio }
func (r *RuidoBlanco) SetIdAudio(id int)                   { r.idAudio = id }

func (r *RuidoBlanco) GetTipoSonido() string               { return r.tipoSonido }
func (r *RuidoBlanco) SetTipoSonido(tipo string)            { r.tipoSonido = tipo }

func (r *RuidoBlanco) GetFuenteAudio() string              { return r.fuenteAudio }
func (r *RuidoBlanco) SetFuenteAudio(fuente string)         { r.fuenteAudio = fuente }

func (r *RuidoBlanco) GetUsoSugerido() string              { return r.usoSugerido }
func (r *RuidoBlanco) SetUsoSugerido(uso string)            { r.usoSugerido = uso }

func (r *RuidoBlanco) GetProveedorContenido() string             { return r.proveedorContenido }
func (r *RuidoBlanco) SetProveedorContenido(proveedor string)     { r.proveedorContenido = proveedor }

func (r *RuidoBlanco) GetDuracionBucle() int               { return r.duracionBucle }
func (r *RuidoBlanco) SetDuracionBucle(duracion int)        { r.duracionBucle = duracion }

func (r *RuidoBlanco) GetFrecuenciaDominante() string            { return r.frecuenciaDominante }
func (r *RuidoBlanco) SetFrecuenciaDominante(frecuencia string)   { r.frecuenciaDominante = frecuencia }
