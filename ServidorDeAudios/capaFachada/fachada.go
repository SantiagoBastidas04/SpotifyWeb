package capaFachada

import (
	"fmt"
	"strconv"

	repo "servidor.local/servidorDeAudios/capaAccesoDatos"
	"servidor.local/servidorDeAudios/modelos"
)



func ObtenerTipos() []modelos.TipoAudio {
	fmt.Println("[Fachada] ObtenerTipos llamado")
	return repo.ObtenerTipos()
}

func ObtenerAudiosPorTipo(idTipo int) ([]modelos.ResumenAudio, error) {
	fmt.Printf("[Fachada] ObtenerAudiosPorTipo llamado con idTipo=%d\n", idTipo)

	tipos := repo.ObtenerTipos()
	tipoValido := false
	for _, t := range tipos {
		if t.GetId() == idTipo {
			tipoValido = true
			break
		}
	}

	if !tipoValido {
		return nil, fmt.Errorf("tipo de audio con id %d no existe", idTipo)
	}

	return repo.ObtenerAudiosPorTipo(idTipo), nil
}

func ObtenerMetadata(idAudio int) (interface{}, int, error) {
	fmt.Printf("[Fachada] ObtenerMetadata llamado con idAudio=%d\n", idAudio)

	idTipo := repo.ObtenerTipoPorAudio(idAudio)

	switch idTipo {
	case 1:
		musica, encontrado := repo.ObtenerMusica(idAudio)
		if !encontrado {
			return nil, -1, fmt.Errorf("musica con id %d no encontrada", idAudio)
		}
		return musica, 1, nil
	case 2:
		podcast, encontrado := repo.ObtenerPodcast(idAudio)
		if !encontrado {
			return nil, -1, fmt.Errorf("podcast con id %d no encontrado", idAudio)
		}
		return podcast, 2, nil
	case 3:
		audiolibro, encontrado := repo.ObtenerAudiolibro(idAudio)
		if !encontrado {
			return nil, -1, fmt.Errorf("audiolibro con id %d no encontrado", idAudio)
		}
		return audiolibro, 3, nil
	case 4:
		ruidoBlanco, encontrado := repo.ObtenerRuidoBlanco(idAudio)
		if !encontrado {
			return nil, -1, fmt.Errorf("ruido blanco con id %d no encontrado", idAudio)
		}
		return ruidoBlanco, 4, nil
	default:
		return nil, -1, fmt.Errorf("audio con id %d no encontrado", idAudio)
	}
}

func ObtenerTodosLosAudios() ([]modelos.ResumenAudio, error) {
	fmt.Println("[Fachada] ObtenerTodosLosAudios llamado")
	var todos []modelos.ResumenAudio
	for idTipo := 1; idTipo <= 4; idTipo++ {
		audios := repo.ObtenerAudiosPorTipo(idTipo)
		todos = append(todos, audios...)
	}
	return todos, nil
}



func AlmacenarAudio(idTipo int, campos map[string]string, data []byte) (int, error) {
	titulo := campos["titulo"]
	fmt.Printf("[Fachada] AlmacenarAudio idTipo=%d titulo='%s'\n", idTipo, titulo)

	// Paso 1: guardar el archivo MP3 en disco con el nombre del titulo.
	// El ServidorStreaming busca los archivos como: audios/<titulo>.mp3
	if err := repo.GuardarArchivoFisico(titulo, data); err != nil {
		return 0, fmt.Errorf("error guardando archivo: %w", err)
	}

	// Paso 2: registrar los metadatos en memoria segun el tipo
	var idAsignado int

	switch idTipo {
	case 1:

		m := modelos.Musica{}
		m.SetTitulo(titulo)
		m.SetArtistaPrincipal(campos["artistaPrincipal"])
		m.SetAlbum(campos["album"])
		m.SetGeneroMusical(campos["generoMusical"])
		m.SetSelloDiscografico(campos["selloDiscografico"])
		anio, _ := strconv.Atoi(campos["anioLanzamiento"])
		m.SetAnioLanzamiento(anio)
		idAsignado = repo.AgregarMusica(m)

	case 2:
		// Podcast
		p := modelos.Podcast{}
		p.SetTituloPodcast(campos["tituloPodcast"])
		p.SetTituloEpisodio(campos["tituloEpisodio"])
		p.SetAnfitrion(campos["anfitrion"])
		p.SetTemporadaEpisodio(campos["temporadaEpisodio"])
		p.SetNotasShow(campos["notasShow"])
		p.SetClasificacionContenido(campos["clasificacionContenido"])
		idAsignado = repo.AgregarPodcast(p)

	case 3:
		// Audiolibro
		a := modelos.Audiolibro{}
		a.SetTitulo(titulo)
		a.SetAutor(campos["autor"])
		a.SetNarrador(campos["narrador"])
		a.SetEditorial(campos["editorial"])
		a.SetIsbn(campos["isbn"])
		a.SetCapitulo(campos["capitulo"])
		idAsignado = repo.AgregarAudiolibro(a)

	case 4:
		// Ruido Blanco
		rb := modelos.RuidoBlanco{}
		rb.SetTipoSonido(campos["tipoSonido"])
		rb.SetFuenteAudio(campos["fuenteAudio"])
		rb.SetUsoSugerido(campos["usoSugerido"])
		rb.SetProveedorContenido(campos["proveedorContenido"])
		duracion, _ := strconv.Atoi(campos["duracionBucle"])
		rb.SetDuracionBucle(duracion)
		rb.SetFrecuenciaDominante(campos["frecuenciaDominante"])
		idAsignado = repo.AgregarRuidoBlanco(rb)

	default:
		return 0, fmt.Errorf("tipo de audio invalido: %d", idTipo)
	}

	fmt.Printf("[Fachada] Audio registrado en memoria con id=%d\n", idAsignado)


	return idAsignado, nil
}



type MetadataDTO struct {
	IdAudio int    `json:"idAudio"`
	Tipo    int    `json:"tipo"`
	Titulo  string `json:"titulo"`

	// Musica
	ArtistaPrincipal  string `json:"artistaPrincipal,omitempty"`
	Album             string `json:"album,omitempty"`
	GeneroMusical     string `json:"generoMusical,omitempty"`
	SelloDiscografico string `json:"selloDiscografico,omitempty"`
	AnioLanzamiento   int    `json:"anioLanzamiento,omitempty"`

	// Podcast
	TituloPodcast          string `json:"tituloPodcast,omitempty"`
	TituloEpisodio         string `json:"tituloEpisodio,omitempty"`
	Anfitrion              string `json:"anfitrion,omitempty"`
	TemporadaEpisodio      string `json:"temporadaEpisodio,omitempty"`
	NotasShow              string `json:"notasShow,omitempty"`
	ClasificacionContenido string `json:"clasificacionContenido,omitempty"`

	// Audiolibro
	Autor     string `json:"autor,omitempty"`
	Narrador  string `json:"narrador,omitempty"`
	Editorial string `json:"editorial,omitempty"`
	Isbn      string `json:"isbn,omitempty"`
	Capitulo  string `json:"capitulo,omitempty"`

	// Ruido Blanco
	TipoSonido          string `json:"tipoSonido,omitempty"`
	FuenteAudio         string `json:"fuenteAudio,omitempty"`
	UsoSugerido         string `json:"usoSugerido,omitempty"`
	ProveedorContenido  string `json:"proveedorContenido,omitempty"`
	DuracionBucle       int    `json:"duracionBucle,omitempty"`
	FrecuenciaDominante string `json:"frecuenciaDominante,omitempty"`
}

func ConvertirMetadataADTO(dato interface{}, idTipo int, idAudio int) MetadataDTO {
	dto := MetadataDTO{IdAudio: idAudio, Tipo: idTipo}

	switch idTipo {
	case 1:
		m := dato.(modelos.Musica)
		dto.Titulo = m.GetTitulo()
		dto.ArtistaPrincipal = m.GetArtistaPrincipal()
		dto.Album = m.GetAlbum()
		dto.GeneroMusical = m.GetGeneroMusical()
		dto.SelloDiscografico = m.GetSelloDiscografico()
		dto.AnioLanzamiento = m.GetAnioLanzamiento()
	case 2:
		p := dato.(modelos.Podcast)
		dto.Titulo = p.GetTituloPodcast()
		dto.TituloPodcast = p.GetTituloPodcast()
		dto.TituloEpisodio = p.GetTituloEpisodio()
		dto.Anfitrion = p.GetAnfitrion()
		dto.TemporadaEpisodio = p.GetTemporadaEpisodio()
		dto.NotasShow = p.GetNotasShow()
		dto.ClasificacionContenido = p.GetClasificacionContenido()
	case 3:
		a := dato.(modelos.Audiolibro)
		dto.Titulo = a.GetTitulo()
		dto.Autor = a.GetAutor()
		dto.Narrador = a.GetNarrador()
		dto.Editorial = a.GetEditorial()
		dto.Isbn = a.GetIsbn()
		dto.Capitulo = a.GetCapitulo()
	case 4:
		rb := dato.(modelos.RuidoBlanco)
		dto.Titulo = rb.GetTipoSonido() + " - " + rb.GetFuenteAudio()
		dto.TipoSonido = rb.GetTipoSonido()
		dto.FuenteAudio = rb.GetFuenteAudio()
		dto.UsoSugerido = rb.GetUsoSugerido()
		dto.ProveedorContenido = rb.GetProveedorContenido()
		dto.DuracionBucle = rb.GetDuracionBucle()
		dto.FrecuenciaDominante = rb.GetFrecuenciaDominante()
	}

	return dto
}