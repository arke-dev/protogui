package service

import (
	"context"

	"github.com/arke-dev/protogui/infra"
	"github.com/arke-dev/protogui/models"
)

type GRPC interface {
	Invoke(ctx context.Context, req models.GRPCRequest) (string, error)
}

type grpcSvc struct {
	protoCompiler ProtoCompiler
	connGRPC      *infra.GRPC
}

func NewGRPC(protoCompiler ProtoCompiler, connGRPC *infra.GRPC) *grpcSvc {
	return &grpcSvc{protoCompiler: protoCompiler, connGRPC: connGRPC}
}

func (g *grpcSvc) Invoke(ctx context.Context, reqGRPC models.GRPCRequest) (string, error) {
	conn, err := g.connGRPC.GetConn(reqGRPC.Address)
	if err != nil {
		return "", err
	}

	reqTypename, resTypename, err := g.protoCompiler.GetRequestResponseFromMethod(reqGRPC.Path, reqGRPC.Method)
	if err != nil {
		return "", err
	}

	req, err := g.protoCompiler.JSONToProto("", reqTypename, reqGRPC.RequestJsonMsg)
	if err != nil {
		return "", err
	}

	res, err := g.protoCompiler.JSONToProto("", resTypename, "{}")
	if err != nil {
		return "", err
	}

	err = conn.Invoke(ctx, reqGRPC.Method, req, res)
	if err != nil {
		return "", err
	}

	return g.protoCompiler.ProtoToJSON(res)
}
