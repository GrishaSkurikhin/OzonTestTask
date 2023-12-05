package restserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/GrishaSkurikhin/OzonTestTask/internal/config"
	saveurl "github.com/GrishaSkurikhin/OzonTestTask/internal/rest-server/handlers/save-url"
	mwLogger "github.com/GrishaSkurikhin/OzonTestTask/internal/rest-server/middleware/logger"
	"github.com/GrishaSkurikhin/OzonTestTask/internal/url"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
)

const (
	ReadTimeout  = 5
	WriteTimeout = 5
	IdleTimeout  = 5
)

type restServer struct {
	*http.Server
}

type Storage interface {
	url.URLSaver
	url.URLGetter
}

func New(cfg config.Server, log *zerolog.Logger, strg Storage) *restServer {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(*log))
	router.Use(middleware.Recoverer)

	router.Route("url", func(r chi.Router) {
		r.Post("/", saveurl.New(*log, strg))
		r.Get("/{shortURL}", saveurl.New(*log, strg))
	})

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  ReadTimeout,
		WriteTimeout: WriteTimeout,
		IdleTimeout:  IdleTimeout,
	}

	return &restServer{srv}
}

func (srv *restServer) Start() error {
	const op = "restserver.Start"

	err := srv.ListenAndServe()
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}
	return nil
}

func (srv *restServer) Close(ctx context.Context) error {
	const op = "restserver.Close"

	err := srv.Shutdown(ctx)
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
