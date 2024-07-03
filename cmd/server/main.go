package main

import (
	"github.com/getz-devs/librakeeper-server/internal/server"
	"github.com/getz-devs/librakeeper-server/internal/server/config"
	"github.com/getz-devs/librakeeper-server/lib/prettylog"
	"log/slog"
	"os"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("starting librakeeper srv", slog.String("env", cfg.Env), slog.Int("port", cfg.Server.Port))

	// Create and initialize the srv
	srv := server.NewServer(cfg, log)
	if err := srv.Initialize(); err != nil {
		log.Error("failed to initialize srv", slog.Any("error", err))
		os.Exit(1)
	}

	// Run the srv
	if err := srv.Run(); err != nil {
		log.Error("failed to run srv", slog.Any("error", err))
		os.Exit(1)
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "local":
		log = slog.New(prettylog.NewHandler(&slog.HandlerOptions{
			Level:       slog.LevelDebug,
			AddSource:   false,
			ReplaceAttr: nil,
		}))
	case "development":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "production":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
