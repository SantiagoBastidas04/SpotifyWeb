package co.edu.unicauca.servidorcanciones.capaFachadaServices.DTO;

import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
public class ReproduccionDTO {

    // Nickname del usuario que reproduce o pausa
    private String nombreUsuario;
    private String idAudio;
}