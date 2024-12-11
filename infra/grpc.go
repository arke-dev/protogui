package infra

import (
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type GRPC struct {
	conns map[string]*grpc.ClientConn
	mux   *sync.Mutex
}

func NewGRPC() *GRPC {
	return &GRPC{conns: make(map[string]*grpc.ClientConn), mux: &sync.Mutex{}}
}

func (g *GRPC) GetConn(address string) (*grpc.ClientConn, error) {
	g.mux.Lock()
	defer g.mux.Unlock()

	conn, ok := g.conns[address]
	if ok {
		return conn, nil
	}

	conn, err := createCommonClientConn(address)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc client %v", err)
	}

	g.conns[address] = conn
	return conn, nil
}

func (g *GRPC) Close() {
	for _, conn := range g.conns {
		conn.Close()
	}
}

func createCommonClientConn(serverAddr string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		serverAddr,

		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  time.Millisecond * 100,
				Multiplier: 1.1,
				Jitter:     0.3,
				MaxDelay:   time.Second * 120,
			},
			MinConnectTimeout: time.Second * 2,
		}),

		grpc.WithTransportCredentials(insecure.NewCredentials()),

		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                time.Second * 30,
			Timeout:             time.Second * 10,
			PermitWithoutStream: true,
		}),
	)
	return conn, err
}
