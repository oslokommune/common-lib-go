package grpcruntime

import (
	"fmt"
	"net"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type grpcRuntime struct {
	grpcServer *grpc.Server
}

func NewServer(enableTracing bool, opt ...grpc.ServerOption) *grpcRuntime {
	if enableTracing {
		opt = append(opt, grpc.StatsHandler(otelgrpc.NewServerHandler()))
	}

	grpcServer := grpc.NewServer(opt...)
	return &grpcRuntime{grpcServer}
}

func (g *grpcRuntime) GetServer() *grpc.Server {
	return g.grpcServer
}

func (g *grpcRuntime) Start(port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to listen on port %d", port)
	}

	log.Info().Msgf("Server started on grpc://localhost:%d", port)

	g.grpcServer.Serve(listener)
}
