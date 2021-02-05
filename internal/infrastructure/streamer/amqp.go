package streamer

import (
	"encoding/json"
	"github.com/oleglarionov/otusgolang_finalproject/internal/application/event"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"log"
)

type AMQPConfig struct {
	Dsn   string
	Queue string
}

type AMQPStreamer struct {
	cfg  AMQPConfig
	conn *amqp.Connection
	ch   *amqp.Channel
	q    amqp.Queue
}

func NewAMQPStreamer(cfg AMQPConfig) *AMQPStreamer {
	return &AMQPStreamer{cfg: cfg}
}

func (s *AMQPStreamer) Push(event event.Event) error {
	bytes, err := json.Marshal(event)
	if err != nil {
		return errors.WithStack(err)
	}

	err = s.ch.Publish(
		"",
		s.q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        bytes,
		})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (s *AMQPStreamer) Connect() error {
	conn, err := amqp.Dial(s.cfg.Dsn)
	if err != nil {
		return errors.WithStack(err)
	}
	s.conn = conn
	log.Println("connection to rabbit established")

	ch, err := conn.Channel()
	if err != nil {
		return errors.WithStack(err)
	}
	s.ch = ch

	q, err := s.ch.QueueDeclare(
		s.cfg.Queue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return errors.WithStack(err)
	}
	s.q = q

	return nil
}

func (s *AMQPStreamer) Close() {
	if s.ch != nil {
		s.ch.Close()
	}

	if s.conn != nil {
		s.conn.Close()
	}
}
