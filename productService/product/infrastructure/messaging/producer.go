package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"micr_course/productService/config"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	name       = "productService"
	type_      = "direct"
	durable    = true
	autoDelete = false
	internal   = false
	noWait     = false
	mandatory  = false
	immediate  = false
)

type RabbitMQPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	cfg     *config.RabbitMQ
}

func NewProducerSetup(cfg *config.RabbitMQ) (*RabbitMQPublisher, error) {
	var conn *amqp.Connection
	var err error

	counts := 0
	maxRetries := 30
	retryDelay := 2 * time.Second

	for {
		conn, err = amqp.Dial(cfg.Addr)
		if err == nil {
			log.Println("✅ Successfully connected to RabbitMQ")
			break
		}

		counts++
		log.Printf("⏳ Failed to connect to RabbitMQ (Attempt %d/%d): %v. Retrying in %v...", counts, maxRetries, err, retryDelay)

		if counts >= maxRetries {
			log.Fatalf("❌ Could not connect to RabbitMQ after %d attempts: %v", maxRetries, err)
		}

		time.Sleep(retryDelay)
	}
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to RabbitMQ: %s", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Failed to open a channel: %s", err)
	}
	err = ch.ExchangeDeclare(
		cfg.Exchange,
		"direct",
		durable,
		autoDelete,
		internal,
		noWait,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("Failed to declare an exchange: %s", err)
	}
	return &RabbitMQPublisher{
		conn:    conn,
		channel: ch,
		cfg:     cfg,
	}, nil
}

func (p *RabbitMQPublisher) SendProductCreated(ctx context.Context, product interface{}) error {
	return p.publish("product.created", product)
}

func (p *RabbitMQPublisher) SendProductUpdated(ctx context.Context, product interface{}) error {
	return p.publish("product.updated", product)
}

func (p *RabbitMQPublisher) SendProductDeleted(ctx context.Context, product interface{}) error {
	return p.publish("product.deleted", product)
}

func (p *RabbitMQPublisher) Close() {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
}

func (p *RabbitMQPublisher) publish(routing string, payload interface{}) error {
	log.Printf("Publishing message to RabbitMQ: %s", p.cfg.Exchange)
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Failed to marshal payload: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err = p.channel.PublishWithContext(
		ctx,
		p.cfg.Exchange,
		routing,
		mandatory,
		immediate,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			MessageId:    uuid.New().String(),
			Timestamp:    time.Now(),
			Body:         body,
		})
	if err != nil {
		return fmt.Errorf("публикация сообщения %w", err)
	}
	return nil
}
