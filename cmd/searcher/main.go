package searcher_cmd

import (
	"github.com/getz-devs/librakeeper-server/internal/searcher/app"
	"github.com/getz-devs/librakeeper-server/internal/searcher/config"
	mongostorage "github.com/getz-devs/librakeeper-server/internal/searcher/storage/mongo"
	"github.com/getz-devs/librakeeper-server/lib/prettylog"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("startingg ...",
		slog.String("env", cfg.Env),
		slog.Any("config", cfg),
		slog.Int("port", cfg.GRPC.Port),
	)

	databaseMongoConfig := mongostorage.DatabaseMongoConfig{
		ConnectUrl: cfg.DatabaseMongo.ConnectURL,
		Database:   cfg.DatabaseMongo.DatabaseName,
		Collection: cfg.DatabaseMongo.CollectionName,
	}

	application := app.New(log, cfg.GRPC.Port, databaseMongoConfig)
	go application.GRPCSrv.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Info("shutting down ...",
		slog.String("signal", sign.String()),
	)

	application.GRPCSrv.Stop()

	application.Storage.Close()

	log.Info("application fully stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(prettylog.NewHandler(&slog.HandlerOptions{
			Level:       slog.LevelInfo,
			AddSource:   false,
			ReplaceAttr: nil,
		}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
