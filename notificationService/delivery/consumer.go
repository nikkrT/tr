package delivery

import (
	"context"
	"fmt"
	"notificationService/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

type NotificationDelivery interface {
	SendCreated(ctx context.Context, data []byte) error
	SendUpdated(ctx context.Context, data []byte) error
	SendDeleted(ctx context.Context, data []byte) error
}

const (
	name       = "notificationQueue"
	kind       = "direct"
	autoDelete = false
	internal   = false
	noWait     = false
	durable    = true
	noLocal    = false
	exclusive  = false
	autoAck    = false
)

type RabbitMQConsumer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	service    NotificationDelivery
	cfg        config.RabbitMQ
}

func NewRabbitMQConsumer(service NotificationDelivery, cfg config.RabbitMQ) (*RabbitMQConsumer, error) {
	conn, err := amqp.Dial(cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to RabbitMQ: %s", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Failed to open a channel: %s", err)
	}
	err = ch.ExchangeDeclare(
		cfg.Exchange,
		kind,
		durable,
		autoDelete,
		internal,
		noWait,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to declare an exchange: %s", err)
	}
	q, err := ch.QueueDeclare(
		name,
		durable,
		autoDelete,
		exclusive,
		noWait,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to declare a queue: %s", err)
	}

	for _, key := range cfg.Keys {
		ch.QueueBind(q.Name, key, cfg.Exchange, autoDelete, nil)
	}
	return &RabbitMQConsumer{
		connection: conn,
		channel:    ch,
		service:    service,
		cfg:        cfg,
	}, nil
}

func (r *RabbitMQConsumer) StartConsumers(workerPoolSize int, consumerTag string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer r.Close()

	err := r.channel.Qos(workerPoolSize, 0, false)
	if err != nil {
		return fmt.Errorf("Failed to set QoS: %s", err)
	}
	msgs, err := r.channel.Consume(
		name,
		consumerTag,
		autoAck,
		exclusive,
		noLocal,
		noWait,
		nil,
	)
	if err != nil {
		return fmt.Errorf("Failed to register a consumer: %s", err)
	}
	for i := 0; i < workerPoolSize; i++ {
		go r.worker(ctx, msgs)
	}

	chanErr := <-r.channel.NotifyClose(make(chan *amqp.Error))

	return chanErr
}

func (r *RabbitMQConsumer) worker(ctx context.Context, msgs <-chan amqp.Delivery) {
	for delivery := range msgs {
		fmt.Println("Got a message\n", delivery.Body)
		fmt.Println("Received a message: ", delivery.Body)
		switch delivery.RoutingKey {
		case r.cfg.Keys[0]:
			if err := r.service.SendCreated(ctx, delivery.Body); err != nil {
				fmt.Printf("Failed to send a message: %s", err)
				if err = delivery.Reject(false); err != nil {
					fmt.Printf("Failed to EVEN reject a message: %s", err)
				}
			} else {
				delivery.Ack(false)
			}
		case r.cfg.Keys[1]:
			if err := r.service.SendUpdated(ctx, delivery.Body); err != nil {
				fmt.Printf("Failed to send a message: %s", err)
				if err = delivery.Reject(false); err != nil {
					fmt.Printf("Failed to EVEN reject a message: %s", err)
				}
			} else {
				delivery.Ack(false)
			}
		case r.cfg.Keys[2]:
			if err := r.service.SendDeleted(ctx, delivery.Body); err != nil {
				fmt.Printf("Failed to send a message: %s", err)
				if err = delivery.Reject(false); err != nil {
					fmt.Printf("Failed to EVEN reject a message: %s", err)
				}
			} else {
				delivery.Ack(false)
			}
		}
	}
}
func (r *RabbitMQConsumer) Close() error {
	err := r.connection.Close()
	if err != nil {
		return fmt.Errorf("Failed to close a connection: %s", err)
	}
	return r.channel.Close()

}
