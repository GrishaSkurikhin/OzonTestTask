package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/GrishaSkurikhin/OzonTestTask/internal/config"
	grpcserver "github.com/GrishaSkurikhin/OzonTestTask/internal/grpc-server"
	"github.com/GrishaSkurikhin/OzonTestTask/internal/logger"
	restserver "github.com/GrishaSkurikhin/OzonTestTask/internal/rest-server"
	inmemory "github.com/GrishaSkurikhin/OzonTestTask/internal/storage/in-memory"
	"github.com/GrishaSkurikhin/OzonTestTask/internal/storage/postgresql"
	"github.com/GrishaSkurikhin/OzonTestTask/pkg/closer"
)

const (
	shutdownTimeout = 5 * time.Second
)

func main() {
	cfg := config.MustLoad()

	zlog, err := logger.New(cfg.Env)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to initialize logger: %v", err))
	}
	zlog.Debug().Msg("debug mode is on")

	c := &closer.Closer{}

	var storage restserver.Storage
	if cfg.InMemory {
		zlog.Info().Msg("using in-memory storage")
		storage = inmemory.New()
	} else {
		zlog.Info().Msg("using postgresql storage")

		postgres, err := postgresql.New(cfg.DB.Source)
		if err != nil {
			log.Fatal(fmt.Sprintf("failed to connect postgresql: %v", err))
		}

		c.Add(postgres.Disconnect)
		storage = postgres
	}

	rserver := restserver.New(cfg, zlog, storage)
	c.Add(rserver.Close)

	gserver := grpcserver.New(zlog, cfg.ShortURLHost, storage)
	

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		zlog.Info().Msg(fmt.Sprintf("starting rest server on port %d", cfg.RestPort))
		if err := rserver.Start(); err != nil {
			log.Fatal(fmt.Sprintf("failed to start rest server: %v", err))
		}
	}()

	go func() {
		zlog.Info().Msg(fmt.Sprintf("starting grpc server on port %d", cfg.GRPCPort))
		if err := grpcserver.Run(gserver, cfg.GRPCPort); err != nil {
			log.Fatal(fmt.Sprintf("failed to start grpc server: %v", err))
		}
	}()

	zlog.Info().Msg("service started")

	<-ctx.Done()
	zlog.Info().Msg("stopping service")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	
	gserver.GracefulStop()
	if err := c.Close(shutdownCtx); err != nil {
		zlog.Error().Str("closer error", err.Error())
	}

	zlog.Info().Msg("service stopped")
}
