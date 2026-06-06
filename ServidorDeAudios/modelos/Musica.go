package modelos

type Musica struct {
	idAudio           int
	titulo            string
	artistaPrincipal  string
	album             string
	generoMusical     string
	selloDiscografico string
	anioLanzamiento   int
}

func (m *Musica) GetIdAudio() int              { return m.idAudio }
func (m *Musica) SetIdAudio(id int)             { m.idAudio = id }

func (m *Musica) GetTitulo() string             { return m.titulo }
func (m *Musica) SetTitulo(titulo string)        { m.titulo = titulo }

func (m *Musica) GetArtistaPrincipal() string          { return m.artistaPrincipal }
func (m *Musica) SetArtistaPrincipal(artista string)    { m.artistaPrincipal = artista }

func (m *Musica) GetAlbum() string              { return m.album }
func (m *Musica) SetAlbum(album string)          { m.album = album }

func (m *Musica) GetGeneroMusical() string       { return m.generoMusical }
func (m *Musica) SetGeneroMusical(genero string) { m.generoMusical = genero }

func (m *Musica) GetSelloDiscografico() string        { return m.selloDiscografico }
func (m *Musica) SetSelloDiscografico(sello string)    { m.selloDiscografico = sello }

func (m *Musica) GetAnioLanzamiento() int           { return m.anioLanzamiento }
func (m *Musica) SetAnioLanzamiento(anio int)        { m.anioLanzamiento = anio }
