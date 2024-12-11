package infra

import (
	"net/http"
	"time"
)

type HTTPCli struct {
	*http.Client
}

func NewHTTPCli() *HTTPCli {
	httpCli := &http.Client{
		Timeout: time.Second * 5,
	}

	return &HTTPCli{httpCli}
}
