package rabbitProvider

import (
	"context"
	"fmt"
	rabbitDefines "github.com/getz-devs/librakeeper-server/lib/rabbit/getz.rabbitProto.v1"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
	"log/slog"
)

type RabbitConfig struct {
	RabbitUrl string
	QueueName string
}

type RabbitService struct {
	log *slog.Logger
	ch  *amqp.Channel
	q   amqp.Queue
}

func New(rabbitConfig RabbitConfig, log *slog.Logger) *RabbitService {
	const op = "rabbitProvider.New"
	log = log.With(slog.String("op", op))

	conn, err := amqp.Dial(rabbitConfig.RabbitUrl)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		rabbitConfig.QueueName, // name
		false,                  // durable
		false,                  // delete when unused
		false,                  // exclusive
		false,                  // no-wait
		nil,                    // arguments
	)
	failOnError(err, "Failed to declare a queue")

	return &RabbitService{
		log: log,
		ch:  ch,
		q:   q,
	}
}

func (s *RabbitService) sendMessage(ctx context.Context, message []byte) error {
	const op = "rabbitProvider.RabbitService.SendMessage"
	s.log.With(slog.String("op", op))
	err := s.ch.PublishWithContext(
		ctx,
		"",       // exchange
		s.q.Name, // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		})
	if err != nil {
		s.log.Error(err.Error())
		return fmt.Errorf("%s: %w", op, err)
	}

	s.log.Info(" [x] Sent ", slog.String("message", string(message)))

	return nil
}
func (s *RabbitService) AddRequest(ctx context.Context, isbn string) error {
	protoMessage := &rabbitDefines.ISBNMessage{
		Isbn: isbn,
	}
	out, err := proto.Marshal(protoMessage)
	if err != nil {
		s.log.Error(err.Error())
		return err
	}
	return s.sendMessage(ctx, out)
}

func (s *RabbitService) Close() {
	const op = "rabbitProvider.RabbitService.Close"
	s.log.With(slog.String("op", op))
	err := s.ch.Close()
	if err != nil {
		s.log.Error(err.Error())
		return
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		panic(fmt.Errorf("%s: %w", msg, err))
	}
}
