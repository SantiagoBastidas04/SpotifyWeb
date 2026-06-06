package co.edu.unicauca.servidorcanciones.capaControladores;

import co.edu.unicauca.servidorcanciones.capaAccesoADatos.RepositorioCancion;
import co.edu.unicauca.servidorcanciones.capaFachadaServices.DTO.NotificacionDTO;
import co.edu.unicauca.servidorcanciones.capaFachadaServices.DTO.ReaccionDTO;
import co.edu.unicauca.servidorcanciones.capaFachadaServices.DTO.ReproduccionDTO;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.messaging.handler.annotation.MessageMapping;
import org.springframework.messaging.handler.annotation.Payload;
import org.springframework.messaging.simp.SimpMessageHeaderAccessor;
import org.springframework.messaging.simp.SimpMessagingTemplate;
import org.springframework.stereotype.Controller;

import java.util.List;

@Controller
public class CancionController {

    @Autowired
    private SimpMessagingTemplate messagingTemplate;

    @Autowired
    private RepositorioCancion repositorioCancion;

    // -----------------------------------------------------------------------
    // CANAL 1 — Reacciones
    // Cliente envia:  /app/reaccion
    // Servidor publica en: /audio/{idAudio}/reacciones
    // Reciben: todos los que estan reproduciendo ese audio
    // -----------------------------------------------------------------------
    @MessageMapping("/reaccion")
    public void enviarReaccion(@Payload ReaccionDTO reaccion,
                               SimpMessageHeaderAccessor headerAccessor) {

        String sessionId = headerAccessor.getSessionId();
        String idAudio = repositorioCancion.getAudioActual(sessionId);

        // Verificar que el usuario realmente este reproduciendo ese audio
        if (idAudio == null || !idAudio.equals(reaccion.getIdAudio())) return;

        List<String> oyentes = repositorioCancion.getOyentes(idAudio);
        if (!oyentes.isEmpty()) {
            messagingTemplate.convertAndSend(
                    "/audio/" + idAudio + "/reacciones", reaccion);
        }
    }

    // -----------------------------------------------------------------------
    // CANAL 2 — Comenzar reproduccion
    // Cliente envia:  /app/reproducir
    // Servidor publica en: /audio/{idAudio}/reproduciendo
    // Reciben: los oyentes que ya estaban antes de que llegara el nuevo
    // -----------------------------------------------------------------------
    @MessageMapping("/reproducir")
    public void comenzarReproduccion(@Payload ReproduccionDTO evento,
                                     SimpMessageHeaderAccessor headerAccessor) {

        String sessionId = headerAccessor.getSessionId();
        String idAudio = evento.getIdAudio();
        String nombreUsuario = evento.getNombreUsuario();

        if (idAudio == null || idAudio.isEmpty() || nombreUsuario == null || nombreUsuario.isEmpty()) {
            System.err.println("[ERROR] Datos invalidos en reproduccion: idAudio=" + idAudio
                    + ", nombreUsuario=" + nombreUsuario);
            return;
        }

        // Registra al usuario y obtiene los oyentes ANTERIORES
        List<String> oyentesAnteriores = repositorioCancion.comenzarReproduccion(
                sessionId, idAudio, nombreUsuario);

        // Lista actualizada (ya incluye al nuevo)
        List<String> oyentesActuales = repositorioCancion.getOyentes(idAudio);

        NotificacionDTO notificacion = new NotificacionDTO(nombreUsuario, "SE_UNE", oyentesActuales);

        // Solo notificar si ya habia alguien escuchando
        if (!oyentesAnteriores.isEmpty()) {
            messagingTemplate.convertAndSend(
                    "/audio/" + idAudio + "/reproduciendo", notificacion);
        }
    }

    // -----------------------------------------------------------------------
    // CANAL 3 — Pausar reproduccion
    // CAMBIADO: endpoint era /app/detener, ahora es /app/pausar
    // porque el requerimiento dice "cuando el cliente pausa"
    // Cliente envia:  /app/pausar
    // Servidor publica en: /audio/{idAudio}/pausas
    // Reciben: los oyentes que siguen escuchando
    // -----------------------------------------------------------------------
    @MessageMapping("/pausar")
    public void pausarReproduccion(@Payload ReproduccionDTO evento,
                                   SimpMessageHeaderAccessor headerAccessor) {

        String sessionId = headerAccessor.getSessionId();
        String idAudio = evento.getIdAudio();
        String nombreUsuario = evento.getNombreUsuario();

        if (idAudio == null || idAudio.isEmpty() || nombreUsuario == null || nombreUsuario.isEmpty()) {
            System.err.println("[ERROR] Datos invalidos en pausa: idAudio=" + idAudio
                    + ", nombreUsuario=" + nombreUsuario);
            return;
        }

        List<String> oyentesRestantes = repositorioCancion.pausarReproduccion(sessionId);

        NotificacionDTO notificacion = new NotificacionDTO(nombreUsuario, "PAUSO", oyentesRestantes);

        if (!oyentesRestantes.isEmpty()) {
            messagingTemplate.convertAndSend(
                    "/audio/" + idAudio + "/pausas", notificacion);
        }
    }
}