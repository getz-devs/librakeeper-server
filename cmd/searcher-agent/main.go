package main

import (
	"github.com/getz-devs/librakeeper-server/internal/searcher-agent/app"
	"github.com/getz-devs/librakeeper-server/internal/searcher-agent/config"
	mongostorage "github.com/getz-devs/librakeeper-server/internal/searcher-agent/storage/mongo"
	"github.com/getz-devs/librakeeper-server/lib/prettylog"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//err := scrapISBNFindBook("9785206000344")
	//if err != nil {
	//	panic(err)
	//}

	//const rabbitUrl = "amqp://guest:guest@192.168.1.161:5672/"

	cfg := config.MustLoad()

	log := prettylog.SetupLogger(cfg.Env)

	log.Info("starting ...",
		slog.String("env", cfg.Env),
		slog.Any("config", cfg),
	)

	databaseMongoConfig := mongostorage.DatabaseMongoConfig{
		ConnectUrl: cfg.DatabaseMongo.ConnectURL,
		Database:   cfg.DatabaseMongo.DatabaseName,
		Collection: cfg.DatabaseMongo.CollectionName,
	}

	application := app.New(cfg.ConnectUrl, cfg.QueueName, databaseMongoConfig, log)
	go application.AppRabbit.MustRun()

	// --------------------------- Register stop signal ---------------------------
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	// --------------------------- Wait for stop signal ---------------------------
	sign := <-stop

	log.Info("shutting down ...",
		slog.String("signal", sign.String()),
	)

	application.AppRabbit.Close()
	application.Storage.Close()

	//application.

	log.Info("application fully stopped")
}
