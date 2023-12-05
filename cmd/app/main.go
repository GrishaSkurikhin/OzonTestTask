package main

import (
	"fmt"
	"log"

	"github.com/GrishaSkurikhin/OzonTestTask/internal/config"
	"github.com/GrishaSkurikhin/OzonTestTask/internal/logger"
	restserver "github.com/GrishaSkurikhin/OzonTestTask/internal/rest-server"
)

func main() {
	cfg := config.MustLoad()

	zlog, err := logger.New(cfg.Env)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to initialize logger: %v", err))
	}

	rserver := restserver.New(cfg, zlog)
	if err := rserver.Start(); err != nil {
		log.Fatal(fmt.Sprintf("failed to start rest server: %v", err))
	}
}
