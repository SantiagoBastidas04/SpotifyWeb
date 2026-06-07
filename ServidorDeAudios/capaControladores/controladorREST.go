package capaControladores

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	fachada "servidor.local/servidorDeAudios/capaFachada"
)

// ─────────────────────────────────────────────
//  CORS
//  Permite que el navegador (cliente web y admin
//  web) haga fetch() a este servidor sin que el
//  navegador bloquee la respuesta.
// ─────────────────────────────────────────────

func agregarCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}


func manejarPreflight(w http.ResponseWriter, r *http.Request) bool {
	agregarCORS(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return true
	}
	return false
}



type ResumenAudioDTO struct {
	IdAudio int    `json:"idAudio"`
	Titulo  string `json:"titulo"`
	IdTipo  int    `json:"idTipo"`
}

type TipoAudioDTO struct {
	Id     int    `json:"id"`
	Nombre string `json:"nombre"`
}

type RespuestaAlmacenamiento struct {
	IdAudio int    `json:"idAudio"`
	Mensaje string `json:"mensaje"`
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

// ─────────────────────────────────────────────
//  POST /canciones/almacenamiento
//  Usado por el Administrador web para registrar
//  un audio nuevo con sus metadatos y archivo MP3.
// ─────────────────────────────────────────────

func AlmacenarAudio(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[REST] POST /canciones/almacenamiento recibido")

	if manejarPreflight(w, r) {
		return
	}
	agregarCORS(w)

	if r.Method != http.MethodPost {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(50 << 20); err != nil {
		http.Error(w, "Error parseando formulario multipart", http.StatusBadRequest)
		return
	}

	// Leer el archivo de audio
	archivo, _, err := r.FormFile("archivo")
	if err != nil {
		http.Error(w, "Error leyendo campo 'archivo'", http.StatusBadRequest)
		return
	}
	defer archivo.Close()

	data, err := io.ReadAll(archivo)
	if err != nil {
		http.Error(w, "Error leyendo bytes del archivo", http.StatusInternalServerError)
		return
	}

	// El campo "tipo" indica el tipo de audio:
	// "1" = Musica, "2" = Podcast, "3" = Audiolibro, "4" = Ruido Blanco
	tipoStr := r.FormValue("tipo")
	idTipo, err := strconv.Atoi(tipoStr)
	if err != nil || idTipo < 1 || idTipo > 4 {
		http.Error(w, "Campo 'tipo' invalido. Debe ser 1, 2, 3 o 4", http.StatusBadRequest)
		return
	}

	// Leer todos los campos del formulario en un mapa.
	// El cliente web y el admin web solo mandan los campos
	// del tipo que corresponde.
	campos := map[string]string{
		"titulo":                 r.FormValue("titulo"),
		"artistaPrincipal":       r.FormValue("artistaPrincipal"),
		"album":                  r.FormValue("album"),
		"generoMusical":          r.FormValue("generoMusical"),
		"selloDiscografico":      r.FormValue("selloDiscografico"),
		"anioLanzamiento":        r.FormValue("anioLanzamiento"),
		"tituloPodcast":          r.FormValue("tituloPodcast"),
		"tituloEpisodio":         r.FormValue("tituloEpisodio"),
		"anfitrion":              r.FormValue("anfitrion"),
		"temporadaEpisodio":      r.FormValue("temporadaEpisodio"),
		"notasShow":              r.FormValue("notasShow"),
		"clasificacionContenido": r.FormValue("clasificacionContenido"),
		"autor":                  r.FormValue("autor"),
		"narrador":               r.FormValue("narrador"),
		"editorial":              r.FormValue("editorial"),
		"isbn":                   r.FormValue("isbn"),
		"capitulo":               r.FormValue("capitulo"),
		"tipoSonido":             r.FormValue("tipoSonido"),
		"fuenteAudio":            r.FormValue("fuenteAudio"),
		"usoSugerido":            r.FormValue("usoSugerido"),
		"proveedorContenido":     r.FormValue("proveedorContenido"),
		"duracionBucle":          r.FormValue("duracionBucle"),
		"frecuenciaDominante":    r.FormValue("frecuenciaDominante"),
	}

	fmt.Printf("[REST] Almacenando audio tipo=%d titulo='%s'\n", idTipo, campos["titulo"])

	idAsignado, err := fachada.AlmacenarAudio(idTipo, campos, data)
	if err != nil {
		http.Error(w, "Error almacenando audio: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(RespuestaAlmacenamiento{
		IdAudio: idAsignado,
		Mensaje: fmt.Sprintf("Audio '%s' almacenado con id %d", campos["titulo"], idAsignado),
	})
}

// ─────────────────────────────────────────────
//  GET /canciones/tipos
//  Devuelve la lista de tipos de audio en JSON.
//  El cliente web lo usa para mostrar el selector
//  de tipo antes de listar los audios.
// ─────────────────────────────────────────────

func ObtenerTiposREST(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[REST] GET /canciones/tipos recibido")

	if manejarPreflight(w, r) {
		return
	}
	agregarCORS(w)

	tipos := fachada.ObtenerTipos()

	var dtos []TipoAudioDTO
	for _, t := range tipos {
		dtos = append(dtos, TipoAudioDTO{
			Id:     t.GetId(),
			Nombre: t.GetNombre(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dtos)
}

// ─────────────────────────────────────────────
//  GET /canciones/por-tipo?idTipo=1
//  Devuelve la lista de audios de un tipo dado.
//  El cliente web lo usa cuando el usuario
//  selecciona un tipo para ver su catalogo.
// ─────────────────────────────────────────────

func ObtenerAudiosPorTipoREST(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[REST] GET /canciones/por-tipo recibido")

	if manejarPreflight(w, r) {
		return
	}
	agregarCORS(w)

	idTipoStr := r.URL.Query().Get("idTipo")
	idTipo, err := strconv.Atoi(idTipoStr)
	if err != nil {
		http.Error(w, "Parametro 'idTipo' invalido", http.StatusBadRequest)
		return
	}

	audios, err := fachada.ObtenerAudiosPorTipo(idTipo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var dtos []ResumenAudioDTO
	for _, a := range audios {
		dtos = append(dtos, ResumenAudioDTO{
			IdAudio: a.GetIdAudio(),
			Titulo:  a.GetTitulo(),
			IdTipo:  a.GetIdTipo(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dtos)
}

// ─────────────────────────────────────────────
//  GET /canciones/metadata?idAudio=101
//  Devuelve todos los metadatos de un audio.
//  El cliente web lo usa cuando el usuario
//  selecciona un audio para ver su detalle.
// ─────────────────────────────────────────────

func ObtenerMetadataREST(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[REST] GET /canciones/metadata recibido")

	if manejarPreflight(w, r) {
		return
	}
	agregarCORS(w)

	idAudioStr := r.URL.Query().Get("idAudio")
	idAudio, err := strconv.Atoi(idAudioStr)
	if err != nil {
		http.Error(w, "Parametro 'idAudio' invalido", http.StatusBadRequest)
		return
	}

	dato, idTipo, err := fachada.ObtenerMetadata(idAudio)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Construir el DTO plano segun el tipo
	dto := fachada.ConvertirMetadataADTO(dato, idTipo, idAudio)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dto)
}

// ─────────────────────────────────────────────
//  GET /canciones
//  Lista todos los audios sin filtrar por tipo.
// ─────────────────────────────────────────────

func ListarAudios(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[REST] GET /canciones recibido")

	if manejarPreflight(w, r) {
		return
	}
	agregarCORS(w)

	audios, err := fachada.ObtenerTodosLosAudios()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var dtos []ResumenAudioDTO
	for _, a := range audios {
		dtos = append(dtos, ResumenAudioDTO{
			IdAudio: a.GetIdAudio(),
			Titulo:  a.GetTitulo(),
			IdTipo:  a.GetIdTipo(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dtos)
}