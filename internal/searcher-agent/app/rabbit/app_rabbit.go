package app_rabbit

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log/slog"
)

type RabbitApp struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	msgs       <-chan amqp.Delivery
	log        *slog.Logger
}

// const rabbitUrl = "amqp://guest:guest@localhost:5672/"

func failOnError(err error, msg string) {
	if err != nil {
		panic(fmt.Errorf("%s: %w", msg, err))
	}
}

func New(rabbitUrl string, queueName string, log *slog.Logger) *RabbitApp {
	const op = "rabbitmq.RabbitApp.New"
	log = log.With(slog.String("op", op))

	conn, err := amqp.Dial(rabbitUrl)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")
	log = log.With("queue", q.Name)

	log.Info("Connected to RabbitMQ")
	return &RabbitApp{
		connection: conn,
		channel:    ch,
		log:        log,
		msgs:       msgs,
	}
}

func (r *RabbitApp) Close() {
	r.channel.Close()
	r.connection.Close()
}

func (r *RabbitApp) Run() error {
	const op = "rabbitmq.RabbitApp.Run"
	log := r.log.With(slog.String("op", op))
	log.Info(" [*] Waiting for messages. To exit press CTRL+C")
	for d := range r.msgs {
		logger := r.log.With(slog.String("messageID", d.MessageId))
		logger.Info("Received a message")
		if err := processMessage(d, logger); err != nil {
			return err
		}
	}
	return nil
}

func processMessage(d amqp.Delivery, log *slog.Logger) error {
	const op = "rabbitmq.processMessage"
	log = log.With(slog.String("op", op))
	log.Info("Processing message",
		slog.String("messageID", d.MessageId),
		slog.Any("body", d.Body))
	return nil
}

func (r *RabbitApp) MustRun() {
	if err := r.Run(); err != nil {
		panic(err)
	}
}
