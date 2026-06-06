package co.edu.unicauca.servidorcanciones.capaAccesoADatos;

import org.springframework.stereotype.Repository;

import java.util.Collections;
import java.util.List;
import java.util.concurrent.ConcurrentHashMap;

@Repository
public class RepositorioCancion {

    // Mapa 1: idAudio -> lista de nicknames que lo estan reproduciendo
    private final ConcurrentHashMap<String, List<String>> audioUsuarios = new ConcurrentHashMap<>();

    // Mapa 2: sessionId -> nickname del usuario
    private final ConcurrentHashMap<String, String> sessionNombreUsuarios = new ConcurrentHashMap<>();

    // Mapa 3: sessionId -> idAudio que esta reproduciendo esa sesion
    private final ConcurrentHashMap<String, String> sessionAudio = new ConcurrentHashMap<>();



    public void registrarSesion(String sessionId, String nombreUsuario) {
        sessionNombreUsuarios.put(sessionId, nombreUsuario);
    }

    public void eliminarSesion(String sessionId) {
        sessionNombreUsuarios.remove(sessionId);
    }

    public String getNickname(String sessionId) {
        return sessionNombreUsuarios.get(sessionId);
    }

    // ─────────────────────────────────────────────
    //  Reproduccion
    // ─────────────────────────────────────────────

   
    public List<String> comenzarReproduccion(String sessionId, String idAudio, String nombreUsuario) {
        sessionAudio.put(sessionId, idAudio);

        List<String> oyentes = audioUsuarios.computeIfAbsent(
                idAudio,
                k -> Collections.synchronizedList(new java.util.ArrayList<>())
        );

        // Copiar ANTES de agregar al nuevo — son los que ya estaban
        List<String> oyentesAnteriores = List.copyOf(oyentes);
        oyentes.add(nombreUsuario);
        return oyentesAnteriores;
    }

    
    public List<String> pausarReproduccion(String sessionId) {
        String idAudio = sessionAudio.remove(sessionId);
        if (idAudio == null) return Collections.emptyList();

        List<String> oyentes = audioUsuarios.getOrDefault(idAudio, Collections.emptyList());

        // CORRECCION: sacar al usuario de la lista antes de copiarla
        String nickname = sessionNombreUsuarios.get(sessionId);
        if (nickname != null) {
            oyentes.remove(nickname);
        }

        return List.copyOf(oyentes);
    }

    public String getAudioActual(String sessionId) {
        return sessionAudio.get(sessionId);
    }

    public List<String> getOyentes(String idAudio) {
        return audioUsuarios.getOrDefault(idAudio, Collections.emptyList());
    }
}