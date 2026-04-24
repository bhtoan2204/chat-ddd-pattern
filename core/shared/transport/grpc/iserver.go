package grpc

import "google.golang.org/grpc"

type GRPCServer interface {
	Register(server grpc.ServiceRegistrar)
}
