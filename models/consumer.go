package models

type Mode string

var (
	Nack Mode = "nack"
	Ack  Mode = "ack"
)

type GetMessagesRequest struct {
	Path     string
	Queue    string
	Quantity int
	Mode     Mode
}

type GetMessagesResponse struct {
	Type    string
	Payload []byte
}
