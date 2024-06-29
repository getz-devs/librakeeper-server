package main

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

func main() {
	cfg := config.MustLoad()

	log := prettylog.SetupLogger(cfg.Env)

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

	// --------------------------- Start Application server -----------------------
	application := app.New(log, cfg.GRPC.Port, databaseMongoConfig)
	go application.GRPCSrv.MustRun()

	// --------------------------- Register stop signal ---------------------------
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	// --------------------------- Wait for stop signal ---------------------------
	sign := <-stop

	log.Info("shutting down ...",
		slog.String("signal", sign.String()),
	)

	application.GRPCSrv.Stop()

	application.Storage.Close()

	log.Info("application fully stopped")
}
