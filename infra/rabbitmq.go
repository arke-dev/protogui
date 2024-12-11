package infra

import (
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	user     string
	password string
	host     string
	vhost    string
	port     int
	connAMQP *amqp.Connection
	ch       *amqp.Channel
	mux      sync.Mutex
}

func NewRabbitMQ(user, password, host, vhost string, port int) *RabbitMQ {
	return &RabbitMQ{
		user:     user,
		password: password,
		host:     host,
		vhost:    vhost,
		port:     port,
	}
}

func (r *RabbitMQ) GetRabbitmqAMQP() (*amqp.Connection, *amqp.Channel) {
	r.mux.Lock()
	defer r.mux.Unlock()

	if r.connAMQP == nil {
		conn, err := r.createConnection("rabbiter")
		if err != nil {
			panic(fmt.Sprintf("failed to connect in rabbitmq %v\n", err))
		}
		r.connAMQP = conn
	}

	if r.ch == nil || r.ch.IsClosed() {
		fmt.Println("open amqp channel")
		ch, err := r.connAMQP.Channel()
		if err != nil {
			panic(fmt.Sprintf("failed to connect in rabbitmq %v\n", err))
		}
		r.ch = ch
	}

	return r.connAMQP, r.ch
}

func (r *RabbitMQ) createConnection(connName string) (*amqp.Connection, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", r.user, r.password, r.host, r.port, r.vhost)

	props := make(map[string]interface{})
	if connName != "" {
		props["connection_name"] = connName
	}

	conn, err := amqp.DialConfig(url, amqp.Config{
		Heartbeat:  time.Second * time.Duration(10),
		Properties: amqp.Table(props),
	})

	return conn, err
}

func (r *RabbitMQ) Close() {
	if r.ch != nil && !r.ch.IsClosed() {
		r.ch.Close()
	}

	if r.connAMQP != nil {
		r.connAMQP.Close()
	}
}
