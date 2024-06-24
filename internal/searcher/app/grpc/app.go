package grpcapp

import (
	"fmt"
	searcherrpc "github.com/getz-devs/librakeeper-server/internal/searcher/grpc"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

// New creates a new App
func New(
	log *slog.Logger,
	port int,
) *App {
	gRPCServer := grpc.NewServer()
	searcherrpc.Register(gRPCServer)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

// MustRun
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run method
func (a *App) Run() error {
	const op = "grpcapp.App.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s ,error when creating listener: %w", op, err)
	}

	log.Info("gRPC server is running",
		slog.String("address", listener.Addr().String()),
	)

	if err := a.gRPCServer.Serve(listener); err != nil {
		return fmt.Errorf("%s ,error when starting gRPC server: %w", op, err)
	}

	return nil
}

// Stop method
func (a *App) Stop() {
	const op = "grpcapp.App.Stop"

	a.log.With(slog.String("op", op)).Info("stopping gRPC server")
	a.gRPCServer.GracefulStop()
}
