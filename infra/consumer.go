package infra

import (
	"context"

	"github.com/arke-dev/protogui/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumerMQ struct {
	rabbitCli *RabbitMQ
}

type MessageRabbit struct {
	Body []byte
	Type string
}

func NewConsumerMQ(rabbitCli *RabbitMQ) *ConsumerMQ {
	return &ConsumerMQ{rabbitCli: rabbitCli}
}

func (r *ConsumerMQ) Get(ctx context.Context, req *models.GetMessagesRequest) ([]models.GetMessagesResponse, error) {
	_, ch := r.rabbitCli.GetRabbitmqAMQP()

	acks := make([]amqp.Delivery, 0)
	defer func() {
		for _, n := range acks {
			if req.Mode == models.Ack {
				n.Ack(true)
				continue
			}
			n.Nack(true, true)
		}
	}()

	messages := make([]models.GetMessagesResponse, 0)

	var (
		msg amqp.Delivery
		ok  bool
		err error
	)

	for range req.Quantity {
		msg, ok, err = ch.Get(req.Queue, false)
		if !ok || err != nil {
			break
		}

		acks = append(acks, msg)
		messages = append(messages, models.GetMessagesResponse{
			Payload: msg.Body,
			Type:    msg.Type,
		})
	}

	return messages, err
}
