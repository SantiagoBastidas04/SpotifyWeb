package capaAccesoDatos

import (
	"sync"

	"servidor.local/servidorDeAudios/modelos"
)

var tiposAudio = []modelos.TipoAudio{}
var musicas = []modelos.Musica{}
var podcasts = []modelos.Podcast{}
var audiolibros = []modelos.Audiolibro{}
var ruidosBlancos = []modelos.RuidoBlanco{}
var mu sync.Mutex
var proximoIdMusica = 103
var proximoIdPodcast = 203
var proximoIdAudiolibro = 303
var proximoIdRuidoBlanco = 403

func init() {
	cargarTipos()
	cargarMusicas()
	cargarPodcasts()
	cargarAudiolibros()
	cargarRuidosBlancos()
}

func cargarTipos() {
	t1 := modelos.TipoAudio{}
	t1.SetId(1)
	t1.SetNombre("Música")

	t2 := modelos.TipoAudio{}
	t2.SetId(2)
	t2.SetNombre("Podcast")

	t3 := modelos.TipoAudio{}
	t3.SetId(3)
	t3.SetNombre("Audiolibro")

	t4 := modelos.TipoAudio{}
	t4.SetId(4)
	t4.SetNombre("Ruido Blanco")

	tiposAudio = append(tiposAudio, t1, t2, t3, t4)
}

func cargarMusicas() {
	m1 := modelos.Musica{}
	m1.SetIdAudio(101)
	m1.SetTitulo("Blinding Lights")
	m1.SetArtistaPrincipal("The Weeknd")
	m1.SetAlbum("After Hours")
	m1.SetGeneroMusical("Synth-pop")
	m1.SetSelloDiscografico("Republic Records")
	m1.SetAnioLanzamiento(2019)

	m2 := modelos.Musica{}
	m2.SetIdAudio(102)
	m2.SetTitulo("Bohemian Rhapsody")
	m2.SetArtistaPrincipal("Queen")
	m2.SetAlbum("A Night at the Opera")
	m2.SetGeneroMusical("Rock")
	m2.SetSelloDiscografico("EMI")
	m2.SetAnioLanzamiento(1975)

	musicas = append(musicas, m1, m2)
}

func cargarPodcasts() {
	p1 := modelos.Podcast{}
	p1.SetIdAudio(201)
	p1.SetTituloPodcast("Lex Fridman Podcast")
	p1.SetTituloEpisodio("Elon Musk War, AI, Aliens")
	p1.SetAnfitrion("Lex Fridman")
	p1.SetTemporadaEpisodio("EP 400")
	p1.SetNotasShow("Conversación sobre el futuro de la humanidad")
	p1.SetClasificacionContenido("Para toda la familia")

	p2 := modelos.Podcast{}
	p2.SetIdAudio(202)
	p2.SetTituloPodcast("El Cartel Paranormal")
	p2.SetTituloEpisodio("Fantasmas en Bogotá")
	p2.SetAnfitrion("Carlos Trujillo")
	p2.SetTemporadaEpisodio("T2 EP15")
	p2.SetNotasShow("Historias de terror en la capital colombiana")
	p2.SetClasificacionContenido("Explícito")

	podcasts = append(podcasts, p1, p2)
}

func cargarAudiolibros() {
	a1 := modelos.Audiolibro{}
	a1.SetIdAudio(301)
	a1.SetTitulo("El Principito")
	a1.SetAutor("Antoine de Saint-Exupéry")
	a1.SetNarrador("Juan Gómez")
	a1.SetEditorial("Salamandra")
	a1.SetIsbn("978-84-9838-507-2")
	a1.SetCapitulo("Capítulo 1")

	a2 := modelos.Audiolibro{}
	a2.SetIdAudio(302)
	a2.SetTitulo("Cien Años de Soledad")
	a2.SetAutor("Gabriel García Márquez")
	a2.SetNarrador("Marco Antonio Sainz")
	a2.SetEditorial("Sudamericana")
	a2.SetIsbn("978-84-397-0476-6")
	a2.SetCapitulo("Capítulo 1")

	audiolibros = append(audiolibros, a1, a2)
}

func cargarRuidosBlancos() {
	r1 := modelos.RuidoBlanco{}
	r1.SetIdAudio(401)
	r1.SetTipoSonido("Ruido Rosa")
	r1.SetFuenteAudio("Lluvia suave")
	r1.SetUsoSugerido("Dormir")
	r1.SetProveedorContenido("RelaxSounds")
	r1.SetDuracionBucle(60)
	r1.SetFrecuenciaDominante("Graves")

	r2 := modelos.RuidoBlanco{}
	r2.SetIdAudio(402)
	r2.SetTipoSonido("Ruido Blanco")
	r2.SetFuenteAudio("Ventilador")
	r2.SetUsoSugerido("Concentración")
	r2.SetProveedorContenido("FocusAudio")
	r2.SetDuracionBucle(30)
	r2.SetFrecuenciaDominante("Agudos")

	ruidosBlancos = append(ruidosBlancos, r1, r2)
}

func ObtenerTipos() []modelos.TipoAudio {
	return tiposAudio
}

func ObtenerAudiosPorTipo(idTipo int) []modelos.ResumenAudio {
	var resultado []modelos.ResumenAudio

	switch idTipo {
	case 1:
		for _, m := range musicas {
			r := modelos.ResumenAudio{}
			r.SetIdAudio(m.GetIdAudio())
			r.SetTitulo(m.GetTitulo())
			r.SetIdTipo(1)
			resultado = append(resultado, r)
		}
	case 2:
		for _, p := range podcasts {
			r := modelos.ResumenAudio{}
			r.SetIdAudio(p.GetIdAudio())
			r.SetTitulo(p.GetTituloPodcast() + " - " + p.GetTituloEpisodio())
			r.SetIdTipo(2)
			resultado = append(resultado, r)
		}
	case 3:
		for _, a := range audiolibros {
			r := modelos.ResumenAudio{}
			r.SetIdAudio(a.GetIdAudio())
			r.SetTitulo(a.GetTitulo())
			r.SetIdTipo(3)
			resultado = append(resultado, r)
		}
	case 4:
		for _, rb := range ruidosBlancos {
			r := modelos.ResumenAudio{}
			r.SetIdAudio(rb.GetIdAudio())
			r.SetTitulo(rb.GetTipoSonido() + " - " + rb.GetFuenteAudio())
			r.SetIdTipo(4)
			resultado = append(resultado, r)
		}
	}

	return resultado
}

func ObtenerMusica(idAudio int) (modelos.Musica, bool) {
	for _, m := range musicas {
		if m.GetIdAudio() == idAudio {
			return m, true
		}
	}
	return modelos.Musica{}, false
}

func ObtenerPodcast(idAudio int) (modelos.Podcast, bool) {
	for _, p := range podcasts {
		if p.GetIdAudio() == idAudio {
			return p, true
		}
	}
	return modelos.Podcast{}, false
}

func ObtenerAudiolibro(idAudio int) (modelos.Audiolibro, bool) {
	for _, a := range audiolibros {
		if a.GetIdAudio() == idAudio {
			return a, true
		}
	}
	return modelos.Audiolibro{}, false
}

func ObtenerRuidoBlanco(idAudio int) (modelos.RuidoBlanco, bool) {
	for _, r := range ruidosBlancos {
		if r.GetIdAudio() == idAudio {
			return r, true
		}
	}
	return modelos.RuidoBlanco{}, false
}

func ObtenerTipoPorAudio(idAudio int) int {
	if idAudio >= 101 && idAudio <= 199 {
		return 1
	}
	if idAudio >= 201 && idAudio <= 299 {
		return 2
	}
	if idAudio >= 301 && idAudio <= 399 {
		return 3
	}
	if idAudio >= 401 && idAudio <= 499 {
		return 4
	}

	if idAudio >= 500 {
		return 1
	}
	return -1
}

func AgregarMusica(m modelos.Musica) int {
	mu.Lock()
	defer mu.Unlock()
	proximoIdMusica++
	m.SetIdAudio(proximoIdMusica)
	musicas = append(musicas, m)
	return proximoIdMusica
}
func AgregarPodcast(p modelos.Podcast) int {
	mu.Lock()
	defer mu.Unlock()
	proximoIdPodcast++
	p.SetIdAudio(proximoIdPodcast)
	podcasts = append(podcasts, p)
	return proximoIdPodcast
}
func AgregarAudiolibro(a modelos.Audiolibro) int {
	mu.Lock()
	defer mu.Unlock()
	proximoIdAudiolibro++
	a.SetIdAudio(proximoIdAudiolibro)
	audiolibros = append(audiolibros, a)
	return proximoIdAudiolibro
}
func AgregarRuidoBlanco(r modelos.RuidoBlanco) int {
	mu.Lock()
	defer mu.Unlock()
	proximoIdRuidoBlanco++
	r.SetIdAudio(proximoIdRuidoBlanco)
	ruidosBlancos = append(ruidosBlancos, r)
	return proximoIdRuidoBlanco
}
