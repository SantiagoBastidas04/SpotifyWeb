
const REST_URL = 'http://localhost:7000';
let tipoActual = null;

// Definición de campos por tipo 
const CAMPOS = {
    1: [
        { id: 'artistaPrincipal', label: 'Artista Principal *', ph: 'Ej: The Weeknd', tipo: 'text' },
        { id: 'album', label: 'Álbum *', ph: 'Ej: After Hours', tipo: 'text' },
        { id: 'generoMusical', label: 'Género Musical *', ph: 'Ej: Synth-pop', tipo: 'text' },
        { id: 'selloDiscografico', label: 'Sello Discográfico', ph: 'Ej: Republic Records', tipo: 'text' },
        { id: 'anioLanzamiento', label: 'Año de Lanzamiento', ph: 'Ej: 2019', tipo: 'number' }
    ],
    2: [
        { id: 'tituloPodcast', label: 'Nombre del Podcast *', ph: 'Ej: Lex Fridman Podcast', tipo: 'text' },
        { id: 'tituloEpisodio', label: 'Título del Episodio *', ph: 'Ej: Elon Musk: War, AI', tipo: 'text' },
        { id: 'anfitrion', label: 'Anfitrión (Host) *', ph: 'Ej: Lex Fridman', tipo: 'text' },
        { id: 'temporadaEpisodio', label: 'Temporada / Episodio', ph: 'Ej: EP 400', tipo: 'text' },
        { id: 'notasShow', label: 'Notas del Show', ph: 'Descripción breve...', tipo: 'text' },
        { id: 'clasificacionContenido', label: 'Clasificación de Contenido', ph: 'Ej: Para toda la familia', tipo: 'text' }
    ],
    3: [
        { id: 'autor', label: 'Autor *', ph: 'Ej: Gabriel García Márquez', tipo: 'text' },
        { id: 'narrador', label: 'Narrador *', ph: 'Ej: Marco Antonio Sainz', tipo: 'text' },
        { id: 'editorial', label: 'Editorial', ph: 'Ej: Sudamericana', tipo: 'text' },
        { id: 'isbn', label: 'ISBN', ph: 'Ej: 978-84-397-0476-6', tipo: 'text' },
        { id: 'capitulo', label: 'Capítulo', ph: 'Ej: Capítulo 1', tipo: 'text' }
    ],
    4: [
        { id: 'tipoSonido', label: 'Tipo de Sonido *', ph: 'Ej: Ruido Blanco, Rosa, Marrón', tipo: 'text' },
        { id: 'fuenteAudio', label: 'Fuente del Audio *', ph: 'Ej: Lluvia, Ventilador, Bosque', tipo: 'text' },
        { id: 'usoSugerido', label: 'Uso Sugerido *', ph: 'Ej: Dormir, Concentración', tipo: 'text' },
        { id: 'proveedorContenido', label: 'Proveedor', ph: 'Ej: RelaxSounds', tipo: 'text' },
        { id: 'duracionBucle', label: 'Duración del Bucle (min)', ph: 'Ej: 60', tipo: 'number' },
        { id: 'frecuenciaDominante', label: 'Frecuencia Dominante', ph: 'Ej: Graves, Agudos', tipo: 'text' }
    ]
};

const NOMBRES_TIPO = { 1: 'Música', 2: 'Podcast', 3: 'Audiolibro', 4: 'Ruido Blanco' };

// Seleccionar tipo
function seleccionarTipo(idTipo) {
    tipoActual = idTipo;

    // Marcar botón activo
    document.querySelectorAll('.tipo-opt').forEach(el => el.classList.remove('active'));
    document.getElementById('tipo-' + idTipo).classList.add('active');

    // Renderizar campos específicos del tipo
    const contenedor = document.getElementById('campos-tipo');
    contenedor.innerHTML = '';

    CAMPOS[idTipo].forEach(campo => {
        const div = document.createElement('div');
        div.className = 'form-group';
        div.innerHTML = `
            <label>${campo.label}</label>
            <input type="${campo.tipo}" id="f-${campo.id}"
                   placeholder="${campo.ph}" ${campo.tipo === 'number' ? 'min="0"' : ''} />
        `;
        contenedor.appendChild(div);
    });

    // Limpiar resultado anterior y mostrar el formulario
    mostrarResultado('', '');
    document.getElementById('f-titulo').value = '';
    document.getElementById('f-archivo').value = '';
    document.getElementById('form-campos').style.display = 'block';
}

// Registrar audio
async function registrarAudio() {
    if (!tipoActual) {
        mostrarResultado('Selecciona un tipo de audio primero.', 'err');
        return;
    }

    const titulo = (document.getElementById('f-titulo').value || '').trim();
    const archivo = document.getElementById('f-archivo').files[0];

    if (!titulo) {
        mostrarResultado('El nombre del archivo es obligatorio.', 'err');
        return;
    }
    if (!archivo) {
        mostrarResultado('Debes seleccionar un archivo de audio.', 'err');
        return;
    }

    // Construir FormData con todos los campos
    const fd = new FormData();
    fd.append('tipo', String(tipoActual));
    fd.append('titulo', titulo);
    fd.append('archivo', archivo, titulo + '.mp3');

    // Campos específicos del tipo seleccionado
    CAMPOS[tipoActual].forEach(campo => {
        const el = document.getElementById('f-' + campo.id);
        if (el) fd.append(campo.id, el.value || '');
    });

    // Deshabilitar botón mientras se envía
    const btn = document.getElementById('btn-registrar');
    btn.disabled = true;
    btn.textContent = '⏳ Enviando...';

    try {
        const resp = await fetch(REST_URL + '/canciones/almacenamiento', {
            method: 'POST',
            body: fd
            // No pongas Content-Type: el navegador lo setea con el boundary correcto
        });

        const texto = await resp.text();
        let data;
        try { data = JSON.parse(texto); } catch (_) { data = { mensaje: texto }; }

        if (resp.ok) {
            mostrarResultado(
                '✅ Audio registrado correctamente.\n' +
                'ID asignado: ' + (data.idAudio || '—') + '\n' +
                (data.mensaje || ''),
                'ok'
            );
            // Limpiar formulario tras éxito
            document.getElementById('f-titulo').value = '';
            document.getElementById('f-archivo').value = '';
            CAMPOS[tipoActual].forEach(campo => {
                const el = document.getElementById('f-' + campo.id);
                if (el) el.value = '';
            });
        } else {
            mostrarResultado('❌ Error del servidor: ' + (data.mensaje || texto), 'err');
        }

    } catch (e) {
        mostrarResultado('❌ No se pudo conectar con el servidor: ' + e.message, 'err');
    } finally {
        btn.disabled = false;
        btn.textContent = '✅ Registrar audio';
    }
}

// Mostrar resultado
function mostrarResultado(msg, tipo) {
    const el = document.getElementById('resultado');
    if (!msg) { el.style.display = 'none'; return; }
    el.className = tipo;          // 'ok' o 'err'
    el.style.display = 'block';
    el.style.whiteSpace = 'pre-line';
    el.textContent = msg;
    el.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
}
