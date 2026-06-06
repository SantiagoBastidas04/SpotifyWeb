package modelos

type Podcast struct {
	idAudio                int
	tituloPodcast          string
	tituloEpisodio         string
	anfitrion              string
	temporadaEpisodio      string
	notasShow              string
	clasificacionContenido string
}

func (p *Podcast) GetIdAudio() int                         { return p.idAudio }
func (p *Podcast) SetIdAudio(id int)                        { p.idAudio = id }

func (p *Podcast) GetTituloPodcast() string                 { return p.tituloPodcast }
func (p *Podcast) SetTituloPodcast(titulo string)            { p.tituloPodcast = titulo }

func (p *Podcast) GetTituloEpisodio() string                { return p.tituloEpisodio }
func (p *Podcast) SetTituloEpisodio(titulo string)           { p.tituloEpisodio = titulo }

func (p *Podcast) GetAnfitrion() string                     { return p.anfitrion }
func (p *Podcast) SetAnfitrion(anfitrion string)             { p.anfitrion = anfitrion }

func (p *Podcast) GetTemporadaEpisodio() string             { return p.temporadaEpisodio }
func (p *Podcast) SetTemporadaEpisodio(temporada string)     { p.temporadaEpisodio = temporada }

func (p *Podcast) GetNotasShow() string                     { return p.notasShow }
func (p *Podcast) SetNotasShow(notas string)                 { p.notasShow = notas }

func (p *Podcast) GetClasificacionContenido() string              { return p.clasificacionContenido }
func (p *Podcast) SetClasificacionContenido(clasificacion string)  { p.clasificacionContenido = clasificacion }
