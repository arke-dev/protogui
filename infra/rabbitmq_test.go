package infra

import (
	"context"
	"fmt"
	"testing"

	"github.com/lawmatsuyama/protogui/models"
)

func TestRabbitmqGet(t *testing.T) {
	rab := NewConsumerMQ(NewRabbitMQ("guest", "guest", "localhost", "", 5672))

	msgs, err := rab.Get(context.Background(), &models.GetMessagesRequest{Queue: "ff.registration.foreign-finder.dlq", Mode: models.Nack, Quantity: 4})
	if err != nil {
		t.Fatalf(err.Error())
	}

	for _, msg := range msgs {
		fmt.Println(string(msg.Payload))
		fmt.Println(msg.Type)
	}

}
