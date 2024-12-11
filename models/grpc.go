package models

type GRPCRequest struct {
	Address        string
	Path           string
	Method         string
	RequestJsonMsg string
}
