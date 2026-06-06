package capaControladores

import (
	"context"

	fachada "servidor.local/servidorDeAudios/capaFachada"
	"servidor.local/servidorDeAudios/modelos"
	pb "servidor.local/servidorDeAudios/serviciosAudio"
)

type ControladorAudios struct {
	pb.UnimplementedServiciosAudioServer
}

func (c *ControladorAudios) ObtenerTipos(ctx context.Context, req *pb.SolicitudVacia) (*pb.RespuestaTipos, error) {
	tipos := fachada.ObtenerTipos()

	var tiposDTO []*pb.TipoAudioDTO
	for _, t := range tipos {
		dto := &pb.TipoAudioDTO{
			Id:     int32(t.GetId()),
			Nombre: t.GetNombre(),
		}
		tiposDTO = append(tiposDTO, dto)
	}

	return &pb.RespuestaTipos{
		Estado: &pb.RespuestaBase{Codigo: 200, Mensaje: "OK"},
		Tipos:  tiposDTO,
	}, nil
}

func (c *ControladorAudios) ObtenerAudiosPorTipo(ctx context.Context, req *pb.SolicitudPorTipo) (*pb.RespuestaListaAudios, error) {
	audios, err := fachada.ObtenerAudiosPorTipo(int(req.GetIdTipo()))
	if err != nil {
		return &pb.RespuestaListaAudios{
			Estado: &pb.RespuestaBase{Codigo: 404, Mensaje: err.Error()},
		}, nil
	}

	var audiosDTO []*pb.ResumenAudioDTO
	for _, a := range audios {
		dto := &pb.ResumenAudioDTO{
			IdAudio: int32(a.GetIdAudio()),
			Titulo:  a.GetTitulo(),
			IdTipo:  int32(a.GetIdTipo()),
		}
		audiosDTO = append(audiosDTO, dto)
	}

	return &pb.RespuestaListaAudios{
		Estado: &pb.RespuestaBase{Codigo: 200, Mensaje: "OK"},
		Audios: audiosDTO,
	}, nil
}

func (c *ControladorAudios) ObtenerMetadata(ctx context.Context, req *pb.SolicitudPorId) (*pb.RespuestaMetadata, error) {
	dato, idTipo, err := fachada.ObtenerMetadata(int(req.GetIdAudio()))
	if err != nil {
		return &pb.RespuestaMetadata{
			Estado: &pb.RespuestaBase{Codigo: 404, Mensaje: err.Error()},
		}, nil
	}

	respuesta := &pb.RespuestaMetadata{
		Estado: &pb.RespuestaBase{Codigo: 200, Mensaje: "OK"},
	}

	switch idTipo {
	case 1:
		m := dato.(modelos.Musica)
		respuesta.Metadato = &pb.RespuestaMetadata_Musica{
			Musica: &pb.MusicaDTO{
				IdAudio:           int32(m.GetIdAudio()),
				Titulo:            m.GetTitulo(),
				ArtistaPrincipal:  m.GetArtistaPrincipal(),
				Album:             m.GetAlbum(),
				GeneroMusical:     m.GetGeneroMusical(),
				SelloDiscografico: m.GetSelloDiscografico(),
				AnioLanzamiento:   int32(m.GetAnioLanzamiento()),
			},
		}
	case 2:
		p := dato.(modelos.Podcast)
		respuesta.Metadato = &pb.RespuestaMetadata_Podcast{
			Podcast: &pb.PodcastDTO{
				IdAudio:                int32(p.GetIdAudio()),
				TituloPodcast:          p.GetTituloPodcast(),
				TituloEpisodio:         p.GetTituloEpisodio(),
				Anfitrion:              p.GetAnfitrion(),
				TemporadaEpisodio:      p.GetTemporadaEpisodio(),
				NotasShow:              p.GetNotasShow(),
				ClasificacionContenido: p.GetClasificacionContenido(),
			},
		}
	case 3:
		a := dato.(modelos.Audiolibro)
		respuesta.Metadato = &pb.RespuestaMetadata_Audiolibro{
			Audiolibro: &pb.AudiolibroDTO{
				IdAudio:   int32(a.GetIdAudio()),
				Titulo:    a.GetTitulo(),
				Autor:     a.GetAutor(),
				Narrador:  a.GetNarrador(),
				Editorial: a.GetEditorial(),
				Isbn:      a.GetIsbn(),
				Capitulo:  a.GetCapitulo(),
			},
		}
	case 4:
		r := dato.(modelos.RuidoBlanco)
		respuesta.Metadato = &pb.RespuestaMetadata_RuidoBlanco{
			RuidoBlanco: &pb.RuidoBlancoDTO{
				IdAudio:             int32(r.GetIdAudio()),
				TipoSonido:          r.GetTipoSonido(),
				FuenteAudio:         r.GetFuenteAudio(),
				UsoSugerido:         r.GetUsoSugerido(),
				ProveedorContenido:  r.GetProveedorContenido(),
				DuracionBucle:       int32(r.GetDuracionBucle()),
				FrecuenciaDominante: r.GetFrecuenciaDominante(),
			},
		}
	}

	return respuesta, nil
}
