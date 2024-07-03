package app

import (
	app_rabbit "github.com/getz-devs/librakeeper-server/internal/searcher-agent/app/rabbit"
	mongostorage "github.com/getz-devs/librakeeper-server/internal/searcher-agent/storage/mongo"
	"log/slog"
)

type App struct {
	AppRabbit *app_rabbit.RabbitApp
	Storage   *mongostorage.Storage
}

func New(
	rabbitUrl string,
	queueName string,
	databaseMongoConfig mongostorage.DatabaseMongoConfig,
	log *slog.Logger,
) *App {

	storage := mongostorage.New(databaseMongoConfig)

	appRabbit := app_rabbit.New(rabbitUrl, queueName, log, storage)

	return &App{
		AppRabbit: appRabbit,
		Storage:   storage,
	}
}
