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
	log.Info("starting librakeeper server", slog.String("env", cfg.Env), slog.Int("port", cfg.Server.Port))

	// Create and run the server
	if err := server.NewServer(cfg, log).Run(); err != nil {
		log.Error("server error", slog.Any("error", err))
		os.Exit(1)
	}
}
