package app

import (
	grpcapp "github.com/getz-devs/librakeeper-server/internal/searcher/app/grpc"
	"github.com/getz-devs/librakeeper-server/internal/searcher/services/searcher"
	mongostorage "github.com/getz-devs/librakeeper-server/internal/searcher/storage/mongo"
	"log/slog"
)

type App struct {
	GRPCSrv *grpcapp.App
	Storage *mongostorage.Storage
}

func New(
	log *slog.Logger,
	grpcPort int,
	databaseMongoConfig mongostorage.DatabaseMongoConfig,
) *App {
	storage := mongostorage.New(databaseMongoConfig)
	searcherService := searcher_service.New(log, storage)
	grpcApp := grpcapp.New(log, searcherService, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
		Storage: storage,
	}

}
