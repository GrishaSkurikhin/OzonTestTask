package main

import (
	"fmt"
	"log"

	"github.com/GrishaSkurikhin/OzonTestTask/internal/config"
	"github.com/GrishaSkurikhin/OzonTestTask/internal/logger"
	restserver "github.com/GrishaSkurikhin/OzonTestTask/internal/rest-server"
	"github.com/GrishaSkurikhin/OzonTestTask/internal/storage/postgresql"
)

func main() {
	cfg := config.MustLoad()

	zlog, err := logger.New(cfg.Env)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to initialize logger: %v", err))
	}

	zlog.Info().Msg(fmt.Sprintf("starting service on %s", cfg.Address))
	zlog.Debug().Msg("debug mode is on")


	storage, err := postgresql.New(cfg.DB.Source)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to connect postgresql: %v", err))
	}

	zlog.Info().Msg("starting rest server")
	rserver := restserver.New(cfg, zlog, storage)
	if err := rserver.Start(); err != nil {
		log.Fatal(fmt.Sprintf("failed to start rest server: %v", err))
	}
}
