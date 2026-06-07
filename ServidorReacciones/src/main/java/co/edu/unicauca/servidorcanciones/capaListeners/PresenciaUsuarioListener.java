package co.edu.unicauca.servidorcanciones.capaListeners;

import co.edu.unicauca.servidorcanciones.capaAccesoADatos.RepositorioCancion;
import co.edu.unicauca.servidorcanciones.capaFachadaServices.DTO.NotificacionDTO;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.event.EventListener;
import org.springframework.messaging.simp.SimpMessagingTemplate;
import org.springframework.messaging.simp.stomp.StompHeaderAccessor;
import org.springframework.stereotype.Component;
import org.springframework.web.socket.messaging.SessionConnectEvent;
import org.springframework.web.socket.messaging.SessionDisconnectEvent;

import java.util.List;

@Component
public class PresenciaUsuarioListener {

    @Autowired
    private RepositorioCancion repositorioCancion;

    @Autowired
    private SimpMessagingTemplate messagingTemplate;

    /**
     * Se dispara cuando un cliente abre la conexion WebSocket.
     * Guarda la relacion sessionId -> nickname en el repositorio.
     */
    @EventListener
    public void handleSessionConnected(SessionConnectEvent event) {
        StompHeaderAccessor accessor = StompHeaderAccessor.wrap(event.getMessage());
        String sessionId = accessor.getSessionId();
        String nickname = accessor.getFirstNativeHeader("nickname");

        if (nickname != null && sessionId != null) {
            repositorioCancion.registrarSesion(sessionId, nickname);
            System.out.println("[CONEXION]  " + nickname + " con sesion " + sessionId);
        }
    }

    /**
     * Se dispara cuando el cliente cierra la pestana o pierde conexion
     * sin haber enviado /app/pausar explicitamente.
     * El servidor limpia el estado y notifica a los oyentes restantes
     * por el Canal 3 (/audio/{idAudio}/pausas).
     */
    @EventListener
    public void handleSessionDisconnect(SessionDisconnectEvent event) {
        StompHeaderAccessor accessor = StompHeaderAccessor.wrap(event.getMessage());
        String sessionId = accessor.getSessionId();

        String nickname = repositorioCancion.getNickname(sessionId);
        String idAudio = repositorioCancion.getAudioActual(sessionId);

        if (repositorioCancion.yaSalioLimpiamente(sessionId)) {
            repositorioCancion.eliminarSesion(sessionId);
            System.out.println("[DESCONEXION LIMPIA] " + nickname + " con sesion " + sessionId);
            return;
        }

        // Salida abrupta 
        List<String> oyentesRestantes = repositorioCancion.pausarReproduccion(sessionId);
        repositorioCancion.eliminarSesion(sessionId);

        if (nickname != null) {
            System.out.println("[DESCONEXION ABRUPTA] " + nickname + " con sesion " + sessionId);

            if (idAudio != null && !oyentesRestantes.isEmpty()) {
                NotificacionDTO notificacion = new NotificacionDTO(nickname, "PAUSO", oyentesRestantes);
                messagingTemplate.convertAndSend("/audio/" + idAudio + "/pausas", notificacion);
            }
        }
    }
}