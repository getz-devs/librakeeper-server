package grpcapp

import (
	"context"
	"fmt"
	searcherrpc "github.com/getz-devs/librakeeper-server/internal/searcher/grpc"
	searcher_service "github.com/getz-devs/librakeeper-server/internal/searcher/services/searcher"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

// New creates a new App
func New(
	log *slog.Logger,
	searchService *searcher_service.SearcherService,
	port int,
) *App {
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			//logging.StartCall, logging.FinishCall,
			logging.PayloadReceived, logging.PayloadSent,
		),
		// Add any other option (check functions starting with logging.With).
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			log.Error("Recovered from panic", slog.Any("panic", p))

			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
		logging.UnaryServerInterceptor(InterceptorLogger(log), loggingOpts...),
	))
	searcherrpc.Register(gRPCServer, searchService, log)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

// InterceptorLogger adapts slog logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
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
