package componenteConexionCola

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

// RabbitPublisher encapsula la conexión y canal con RabbitMQ.
type RabbitPublisher struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	queueName string
}

type NotificacionCancion struct {
	IdAudio           int    `json:"idAudio"`
	Titulo            string `json:"titulo"`
	Artista           string `json:"artista"`
	Genero            string `json:"genero"`
	FechaHoraRegistro string `json:"fechaHoraRegistro"`
}

// Crea una nueva conexion con rabbit
func NewRabbitPublisher() (*RabbitPublisher, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, fmt.Errorf("error conectando con RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("error abriendo canal RabbitMQ: %w", err)
	}

	q, err := ch.QueueDeclare(
		"notificaciones_canciones",
		false, false, false, false, nil,
	)
	if err != nil {
		return nil, fmt.Errorf("error declarando cola: %w", err)
	}

	return &RabbitPublisher{conn: conn, channel: ch, queueName: q.Name}, nil
}

func (p *RabbitPublisher) PublicarNotificacion(msg NotificacionCancion) error {
	if msg.FechaHoraRegistro == "" {
		msg.FechaHoraRegistro = time.Now().Format("2006-01-02 15:04:05")
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("error serializando notificación: %w", err)
	}

	err = p.channel.Publish(
		"",
		p.queueName,
		false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("error publicando en cola: %w", err)
	}

	fmt.Printf("[RabbitMQ] Notificación publicada: %s\n", string(body))
	return nil
}

// Cerrar libera el canal y la conexión con RabbitMQ.
func (p *RabbitPublisher) Cerrar() {
	p.channel.Close()
	p.conn.Close()
}
