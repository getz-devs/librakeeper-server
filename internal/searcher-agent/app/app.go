package app

import (
	app_rabbit "github.com/getz-devs/librakeeper-server/internal/searcher-agent/app/rabbit"
	"log/slog"
)

type App struct {
	AppRabbit *app_rabbit.RabbitApp
}

func New(rabbitUrl string, queueName string, log *slog.Logger) *App {
	appRabbit := app_rabbit.New(rabbitUrl, queueName, log)
	return &App{
		AppRabbit: appRabbit,
	}
}
