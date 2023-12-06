package restserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/GrishaSkurikhin/OzonTestTask/internal/config"
	geturl "github.com/GrishaSkurikhin/OzonTestTask/internal/rest-server/handlers/get-url"
	saveurl "github.com/GrishaSkurikhin/OzonTestTask/internal/rest-server/handlers/save-url"
	mwLogger "github.com/GrishaSkurikhin/OzonTestTask/internal/rest-server/middleware/logger"
	"github.com/GrishaSkurikhin/OzonTestTask/internal/service/shortlinks"
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
	shortlinks.URLSaver
	shortlinks.URLGetter
}

type Shortlinks interface {
	geturl.ServiceURLGetter
}

func New(cfg config.Server, log *zerolog.Logger, strg Storage) *restServer {
	shortlinks := &shortlinks.ShortlinksService{}
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(*log))
	router.Use(middleware.Recoverer)

	router.Route("/url", func(r chi.Router) {
		r.Post("/", saveurl.New(*log, strg, cfg.Address, shortlinks))
		r.Get("/{shortURL}", geturl.New(*log, strg, shortlinks))
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
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}
