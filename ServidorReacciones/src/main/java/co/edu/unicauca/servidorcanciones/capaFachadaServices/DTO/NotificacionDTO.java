package co.edu.unicauca.servidorcanciones.capaFachadaServices.DTO;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.List;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class NotificacionDTO {
    private String nombreUsuario;
    private String accion;
    private List<String> oyentesActuales;
}
