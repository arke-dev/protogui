package service

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/lawmatsuyama/protogui/infra"
	"github.com/lawmatsuyama/protogui/models"
)

type Consumer interface {
	GetMessages(ctx context.Context, req *models.GetMessagesRequest) (string, error)
}

type consumer struct {
	consumerMQ    *infra.ConsumerMQ
	protoCompiler ProtoCompiler
}

func NewConsumer(consumerMQ *infra.ConsumerMQ, protoCompiler ProtoCompiler) *consumer {
	return &consumer{consumerMQ: consumerMQ, protoCompiler: protoCompiler}
}

func (c *consumer) GetMessages(ctx context.Context, req *models.GetMessagesRequest) (string, error) {
	if req.Quantity > 50 {
		return "", errors.New("too many messages to retrieve")
	}

	err := c.protoCompiler.RegisterProto(req.Path)
	if err != nil {
		return "", err
	}

	msgs, err := c.getMessages(ctx, req)
	if err != nil {
		return "", err
	}

	// res := make([]string, 0)
	var res string
	for _, msg := range msgs {
		m, err := c.protoCompiler.DecodeNoBase64("", msg.Type, msg.Payload)
		if err != nil {
			fmt.Printf("failed to decode message %v\n", err)
			s := base64.StdEncoding.EncodeToString(msg.Payload)
			m = string(s)
		}

		formatedMessage := c.formatMessage(msg.Type, m)
		if res == "" {
			res = formatedMessage
			continue
		}
		res = fmt.Sprintf("%s\n\n%s", res, formatedMessage)
	}

	return res, nil
}

func (c *consumer) formatMessage(typename string, msg string) string {
	line1 := strings.Repeat("-", 10)
	header := fmt.Sprintf("typename: %s", typename)
	line2 := strings.Repeat("-", 10)

	return fmt.Sprintf("%s\n%s\n%s\n%s", line1, header, line2, msg)
}

func (c *consumer) getMessages(ctx context.Context, req *models.GetMessagesRequest) ([]models.GetMessagesResponse, error) {

	return c.consumerMQ.Get(ctx, req)
}
