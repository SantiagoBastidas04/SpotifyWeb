
// Configuración
const REST_URL = 'http://localhost:7000';
const WS_URL = 'http://localhost:5000/ws';
const ICONOS = { 1: '🎵', 2: '🎙️', 3: '📚', 4: '🌊' };

// Estado de la aplicación
let stompClient = null;
let nickname = null;
let audioActual = null;
let suscripciones = [];
let cerrando = false;
let iniciandoStream = false;
//  Referencia al reproductor de audio 
const reproductor = document.getElementById('audio-player');

// Exponer log temprano para que bundle.js pueda usarlo
// (bundle.js llama window.__writeAudioLog en el evento 'data')
window.__writeAudioLog = function (msg, nivel) {
    escribirLog(msg, nivel || 'info');
};


//  LOGIN

function iniciarSesion() {
    const input = document.getElementById('nick-input');
    const nick = (input.value || '').trim();
    if (!nick) { input.focus(); return; }

    nickname = nick;
    document.getElementById('user-badge').textContent = '👤 ' + nick;
    document.getElementById('login-overlay').style.display = 'none';
    document.getElementById('app').style.display = 'flex';

    escribirLog('Bienvenido, ' + nick, 'ok');
    conectarWebSocket();
    cargarTipos();
    inicializarReproductor();
}

// Enter en el campo de nickname
document.getElementById('nick-input')
    .addEventListener('keydown', e => { if (e.key === 'Enter') iniciarSesion(); });


//  CATÁLOGO — llamadas REST

async function cargarTipos() {
    try {
        const resp = await fetch(REST_URL + '/canciones/tipos');
        const tipos = await resp.json();
        renderizarTipos(tipos);
    } catch (e) {
        document.getElementById('tipos-wrap').innerHTML =
            '<div class="placeholder" style="color:#e94560">Error cargando tipos: ' + e.message + '</div>';
    }
}

function renderizarTipos(tipos) {
    const contenedor = document.getElementById('tipos-wrap');
    contenedor.innerHTML = '';
    tipos.forEach(t => {
        const btn = document.createElement('button');
        btn.className = 'type-btn';
        btn.innerHTML = `<span class="type-icon">${ICONOS[t.id] || '🎧'}</span>${t.nombre}`;
        btn.onclick = () => seleccionarTipo(t, btn);
        contenedor.appendChild(btn);
    });
}

async function seleccionarTipo(tipo, btnEl) {
    document.querySelectorAll('.type-btn').forEach(b => b.classList.remove('active'));
    btnEl.classList.add('active');

    document.getElementById('detalle-sec').style.display = 'none';
    document.getElementById('audios-sec').style.display = 'block';
    document.getElementById('audios-sec-title').textContent =
        (ICONOS[tipo.id] || '🎧') + ' ' + tipo.nombre;
    document.getElementById('audios-wrap').innerHTML =
        '<div class="placeholder">Cargando audios...</div>';

    try {
        const resp = await fetch(REST_URL + '/canciones/por-tipo?idTipo=' + tipo.id);
        const audios = await resp.json();
        renderizarAudios(audios);
    } catch (e) {
        document.getElementById('audios-wrap').innerHTML =
            '<div class="placeholder" style="color:#e94560">Error: ' + e.message + '</div>';
    }
}

function renderizarAudios(audios) {
    const contenedor = document.getElementById('audios-wrap');
    contenedor.innerHTML = '';
    if (!audios || audios.length === 0) {
        contenedor.innerHTML = '<div class="placeholder">No hay audios disponibles</div>';
        return;
    }
    audios.forEach(a => {
        const btn = document.createElement('button');
        btn.className = 'audio-btn';
        btn.innerHTML = `<div class="at">${a.titulo}</div>
                         <div class="as">ID: ${a.idAudio} · Tipo: ${a.idTipo}</div>`;
        btn.onclick = () => seleccionarAudio(a, btn);
        contenedor.appendChild(btn);
    });
}

async function seleccionarAudio(audio, btnEl) {
    document.querySelectorAll('.audio-btn').forEach(b => b.classList.remove('active'));
    btnEl.classList.add('active');

    document.getElementById('detalle-sec').style.display = 'block';
    document.getElementById('detalle-wrap').innerHTML =
        '<div class="placeholder">Cargando detalles...</div>';

    try {
        const resp = await fetch(REST_URL + '/canciones/metadata?idAudio=' + audio.idAudio);
        const meta = await resp.json();
        renderizarDetalle(meta);
    } catch (e) {
        document.getElementById('detalle-wrap').innerHTML =
            '<div class="placeholder" style="color:#e94560">Error: ' + e.message + '</div>';
    }
}

function renderizarDetalle(meta) {
    const tipo = meta.tipo;
    let filas = '';
    let titulo = meta.titulo || '—';

    if (tipo === 1) {
        filas = `
          <div class="meta-row"><span class="mk">Artista</span><span class="mv">${meta.artistaPrincipal || '—'}</span></div>
          <div class="meta-row"><span class="mk">Álbum</span><span class="mv">${meta.album || '—'}</span></div>
          <div class="meta-row"><span class="mk">Género</span><span class="mv">${meta.generoMusical || '—'}</span></div>
          <div class="meta-row"><span class="mk">Sello</span><span class="mv">${meta.selloDiscografico || '—'}</span></div>
          <div class="meta-row"><span class="mk">Año</span><span class="mv">${meta.anioLanzamiento || '—'}</span></div>`;

    } else if (tipo === 2) {
        titulo = meta.tituloPodcast || meta.titulo || '—';
        filas = `
          <div class="meta-row"><span class="mk">Episodio</span><span class="mv">${meta.tituloEpisodio || '—'}</span></div>
          <div class="meta-row"><span class="mk">Anfitrión</span><span class="mv">${meta.anfitrion || '—'}</span></div>
          <div class="meta-row"><span class="mk">Temporada</span><span class="mv">${meta.temporadaEpisodio || '—'}</span></div>
          <div class="meta-row"><span class="mk">Notas</span><span class="mv">${meta.notasShow || '—'}</span></div>
          <div class="meta-row"><span class="mk">Contenido</span><span class="mv">${meta.clasificacionContenido || '—'}</span></div>`;

    } else if (tipo === 3) {
        filas = `
          <div class="meta-row"><span class="mk">Autor</span><span class="mv">${meta.autor || '—'}</span></div>
          <div class="meta-row"><span class="mk">Narrador</span><span class="mv">${meta.narrador || '—'}</span></div>
          <div class="meta-row"><span class="mk">Editorial</span><span class="mv">${meta.editorial || '—'}</span></div>
          <div class="meta-row"><span class="mk">ISBN</span><span class="mv">${meta.isbn || '—'}</span></div>
          <div class="meta-row"><span class="mk">Capítulo</span><span class="mv">${meta.capitulo || '—'}</span></div>`;

    } else if (tipo === 4) {
        titulo = (meta.tipoSonido || '') + ' — ' + (meta.fuenteAudio || '');
        filas = `
          <div class="meta-row"><span class="mk">Tipo sonido</span><span class="mv">${meta.tipoSonido || '—'}</span></div>
          <div class="meta-row"><span class="mk">Fuente</span><span class="mv">${meta.fuenteAudio || '—'}</span></div>
          <div class="meta-row"><span class="mk">Uso sugerido</span><span class="mv">${meta.usoSugerido || '—'}</span></div>
          <div class="meta-row"><span class="mk">Duración bucle</span><span class="mv">${meta.duracionBucle || '—'} min</span></div>
          <div class="meta-row"><span class="mk">Proveedor</span><span class="mv">${meta.proveedorContenido || '—'}</span></div>`;
    }

    // Guardar estado del audio seleccionado
    audioActual = {
        idAudio: meta.idAudio,
        titulo: meta.titulo || titulo,   // nombre del archivo (sin extensión)
        formato: 'mp3',
        displayTitle: titulo,
        idTipo: tipo
    };

    document.getElementById('detalle-wrap').innerHTML = `
      <div class="meta-card">
        <div class="mc-title">${ICONOS[tipo] || '🎧'} ${titulo}</div>
        ${filas}
        <button class="btn-play" onclick="reproducirAudio()">▶ Reproducir</button>
      </div>`;
}


//  STREAMING — gRPC-Web (bundle.js)

function reproducirAudio() {
    if (!audioActual) return;

    const { idAudio, titulo, formato, displayTitle } = audioActual;

    const btnPlay = document.querySelector('.btn-play');
    if (btnPlay) { btnPlay.disabled = true; btnPlay.textContent = '⏳ Descargando...'; }

    const npEl = document.getElementById('now-playing');
    npEl.className = 'now-playing';
    npEl.textContent = (ICONOS[audioActual.idTipo] || '🎵') + ' ' + displayTitle;
    document.getElementById('stream-status').textContent = 'Descargando fragmentos...';

    reproductor.style.display = 'none';
    reproductor.src = '';

    // Suscribirse PRIMERO y esperar antes de iniciar el stream
    suscribirseAlAudio(String(idAudio));
    iniciandoStream = true;
    setTimeout(() => {
        // Notificar reproducción después de suscribirse
        notificarReproduccion(String(idAudio));

        escribirLog('Iniciando streaming: ' + titulo + '.' + formato, 'info');

        if (typeof window.iniciar_streaming_cancion === 'function') {
            window.iniciar_streaming_cancion(titulo, formato);
        } else if (typeof window.iniciarStreamGRPCImpl === 'function') {
            window.iniciarStreamGRPCImpl(titulo, formato);
        } else {
            escribirLog('Error: bundle.js no encontrado o no inicializado', 'err');
        }

        setTimeout(() => {
            if (btnPlay) { btnPlay.disabled = false; btnPlay.textContent = '▶ Reproducir'; }
        }, 2500);
    }, 300); // 300ms para que STOMP complete la suscripción
}


//  REPRODUCTOR — eventos del elemento <audio>

function inicializarReproductor() {

    reproductor.addEventListener('play', function () {
        habilitarReacciones(true);
        document.getElementById('stream-status').textContent = '▶ Reproduciendo';
        escribirLog('Reproducción iniciada ▶', 'ok');
        // Notificar al servidor que se reanuda la reproducción
        if (audioActual && !iniciandoStream) {
            notificarReproduccion(String(audioActual.idAudio));
        }
        iniciandoStream = false;
    });

    reproductor.addEventListener('pause', function () {
        if (cerrando) return; // ignorar pausa provocada por cerrarSesion
        document.getElementById('stream-status').textContent = '⏸ Pausado';
        escribirLog('Reproducción pausada ⏸', 'info');
        if (audioActual) {
            notificarPausa(String(audioActual.idAudio));
        }
    });
    reproductor.addEventListener('ended', function () {
        habilitarReacciones(false);
        document.getElementById('stream-status').textContent = '✅ Finalizado';
        escribirLog('Audio finalizado ✅', 'ok');
        if (audioActual) {
            notificarPausa(String(audioActual.idAudio));
        }
    });

    reproductor.addEventListener('error', function () {
        // Ignorar errores cuando el src está vacío (estado inicial)
        if (!reproductor.src || reproductor.src === window.location.href) return;
        escribirLog('Error en el reproductor de audio ❌', 'err');
        document.getElementById('stream-status').textContent = 'Error al reproducir el audio.';
    });
}


//  WEBSOCKET — conexión STOMP sobre SockJS

function conectarWebSocket() {
    try {
        const socket = new SockJS(WS_URL);
        stompClient = Stomp.over(socket);
        stompClient.debug = null;

        stompClient.connect(
            { nickname: nickname },
            function onConectado() {
                document.getElementById('ws-badge').textContent = '● Conectado';
                document.getElementById('ws-badge').className = 'badge on';
                escribirLog('Servidor de reacciones conectado ✓', 'ok');
            },
            function onError(err) {
                document.getElementById('ws-badge').textContent = '● Desconectado';
                document.getElementById('ws-badge').className = 'badge off';
                escribirLog('WebSocket desconectado, reintentando en 5 s...', 'err');
                setTimeout(conectarWebSocket, 5000);
            }
        );
    } catch (e) {
        escribirLog('No se pudo conectar al servidor de reacciones: ' + e.message, 'err');
    }
}

function suscribirseAlAudio(idAudio) {
    console.log('suscribirseAlAudio llamado para audio:', idAudio, 
                'suscripciones activas:', suscripciones.length);
    if (!stompClient || !stompClient.connected) return;

    // Cancelar suscripciones anteriores
    suscripciones.forEach(s => { try { s.unsubscribe(); } catch (_) { } });
    suscripciones = [];

    // Canal 1: reacciones de otros oyentes
    suscripciones.push(
        stompClient.subscribe('/audio/' + idAudio + '/reacciones', function (msg) {
            const data = JSON.parse(msg.body);
            // No mostrar nuestra propia reacción (el servidor ya la validó)
            agregarNotificacion(
                data.tipoReaccion + ' <b>' + data.nombreUsuario + '</b> reaccionó',
                'react'
            );
            mostrarEmojiFlotante(data.tipoReaccion);
        })
    );

    // Canal 2: nuevo oyente se une
    suscripciones.push(
        stompClient.subscribe('/audio/' + idAudio + '/reproduciendo', function (msg) {
            const data = JSON.parse(msg.body);
            agregarNotificacion(
                '🎵 <b>' + data.nombreUsuario + '</b> comenzó a escuchar este audio',
                'info'
            );
            escribirLog('👤 ' + data.nombreUsuario + ' se unió', 'info');
        })
    );

    // Canal 3: oyente pausa
    suscripciones.push(
        stompClient.subscribe('/audio/' + idAudio + '/pausas', function (msg) {
            const data = JSON.parse(msg.body);
            agregarNotificacion(
                '⏸ <b>' + data.nombreUsuario + '</b> pausó el audio',
                'warn'
            );
            escribirLog('⏸ ' + data.nombreUsuario + ' pausó', 'info');
        })
    );
    // Canal 4: oyente sale
    suscripciones.push(
        stompClient.subscribe('/audio/' + idAudio + '/salidas', function (msg) {
            const data = JSON.parse(msg.body);
            agregarNotificacion(
                '🚪 <b>' + data.nombreUsuario + '</b> salió del audio',
                'warn'
            );
            escribirLog('🚪 ' + data.nombreUsuario + ' salió', 'info');
        })
    );
    escribirLog('Suscrito a canales del audio ' + idAudio, 'ok');
}

function notificarReproduccion(idAudio) {
    if (!stompClient || !stompClient.connected) return;
    stompClient.send('/app/reproducir', {}, JSON.stringify({
        nombreUsuario: nickname,
        idAudio: idAudio
    }));
}

function notificarPausa(idAudio) {
    if (!stompClient || !stompClient.connected) return;
    stompClient.send('/app/pausar', {}, JSON.stringify({
        nombreUsuario: nickname,
        idAudio: idAudio
    }));
}


//  REACCIONES

function enviarReaccion(emoji) {
    if (!stompClient || !stompClient.connected || !audioActual) return;
    stompClient.send('/app/reaccion', {}, JSON.stringify({
        nombreUsuario: nickname,
        idAudio: String(audioActual.idAudio),
        tipoReaccion: emoji
    }));
    // Feedback visual inmediato para quien envía
    //mostrarEmojiFlotante(emoji);
}

function habilitarReacciones(habilitar) {
    ['br-1', 'br-2', 'br-3', 'br-4', 'br-5', 'br-6'].forEach(id => {
        const btn = document.getElementById(id);
        if (btn) btn.disabled = !habilitar;
    });
}

function mostrarEmojiFlotante(emoji) {
    const el = document.createElement('div');
    el.className = 'efloat';
    el.textContent = emoji;
    el.style.left = Math.floor(Math.random() * 60 + 20) + '%';
    el.style.bottom = '180px';
    document.body.appendChild(el);
    setTimeout(() => el.remove(), 2200);
}


//  HELPERS DE UI

function agregarNotificacion(html, tipo) {
    const lista = document.getElementById('notif-list');
    const ph = lista.querySelector('.placeholder');
    if (ph) ph.remove();

    const div = document.createElement('div');
    div.className = 'notif ' + tipo;
    const hora = new Date().toLocaleTimeString();
    div.innerHTML = `<span class="nt">${hora}</span>${html}`;
    lista.insertBefore(div, lista.firstChild);

    // Mantener máximo 40 notificaciones
    while (lista.children.length > 40) lista.lastChild.remove();
}

function escribirLog(msg, nivel) {
    const area = document.getElementById('log-area');
    if (!area) { console.log('[LOG]', msg); return; }

    const div = document.createElement('div');
    div.className = 'le ' + (nivel || 'info');
    div.textContent = '[' + new Date().toLocaleTimeString() + '] ' + msg;
    area.appendChild(div);
    area.scrollTop = area.scrollHeight;

    // Conservar máximo 60 entradas
    while (area.children.length > 60) area.firstChild.remove();
}


function pedirCancion(titulo, formato) {
    if (typeof window.iniciar_streaming_cancion === 'function') {
        return window.iniciar_streaming_cancion(titulo, formato || 'mp3');
    }
    if (typeof window.iniciarStreamGRPCImpl === 'function') {
        return window.iniciarStreamGRPCImpl(titulo, formato || 'mp3');
    }
    escribirLog('Error: implementación gRPC no encontrada', 'err');
}
function cerrarSesion() {
    cerrando = true;

    if (audioActual && stompClient && stompClient.connected) {
        notificarSalida(String(audioActual.idAudio));
    }
    if (stompClient && stompClient.connected) {
        stompClient.disconnect();
    }

    reproductor.pause();
    reproductor.src = '';
    reproductor.style.display = 'none';
    habilitarReacciones(false);

    const npEl = document.getElementById('now-playing');
    npEl.className = 'now-playing empty';
    npEl.textContent = 'Selecciona un audio del catálogo para reproducir';
    document.getElementById('stream-status').textContent = '';

    document.getElementById('notif-list').innerHTML =
        '<div class="placeholder">Las notificaciones aparecerán aquí...</div>';
    document.getElementById('log-area').innerHTML = '';

    document.getElementById('ws-badge').textContent = '● Desconectado';
    document.getElementById('ws-badge').className = 'badge off';
    document.getElementById('user-badge').textContent = '👤 —';

    document.getElementById('app').style.display = 'none';
    document.getElementById('login-overlay').style.display = 'flex';

    nickname = null;
    audioActual = null;
    suscripciones = [];
    document.getElementById('nick-input').value = '';

    cerrando = false; // resetear bandera
}

// Cierre abrupto del navegador
window.addEventListener('beforeunload', function () {
    if (audioActual && stompClient && stompClient.connected) {
        notificarSalida(String(audioActual.idAudio));
        stompClient.disconnect();
    }
});

function notificarSalida(idAudio) {
    if (!stompClient || !stompClient.connected) return;
    stompClient.send('/app/salir', {}, JSON.stringify({
        nombreUsuario: nickname,
        idAudio: idAudio
    }));
}