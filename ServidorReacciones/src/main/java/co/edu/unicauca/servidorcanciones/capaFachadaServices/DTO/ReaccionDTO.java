package co.edu.unicauca.servidorcanciones.capaFachadaServices.DTO;

import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
public class ReaccionDTO {

    
    private String nombreUsuario;
    private String idAudio;
    private String tipoReaccion;
}