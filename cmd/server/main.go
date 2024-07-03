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

	log := prettylog.SetupLogger(cfg.Env)
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
