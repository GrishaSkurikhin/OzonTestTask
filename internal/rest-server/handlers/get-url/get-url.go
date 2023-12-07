package geturl

import (
	"context"
	"errors"
	"io"
	"net/http"

	customerrors "github.com/GrishaSkurikhin/OzonTestTask/internal/custom-errors"
	resp "github.com/GrishaSkurikhin/OzonTestTask/internal/rest-server/api/response"
	"github.com/GrishaSkurikhin/OzonTestTask/internal/service/shortlinks"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
)

type Request struct {
	ShortURL string `json:"shortURL"`
}

type Response struct {
	resp.Response
	Url string `json:"url"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=ServiceURLGetter
type ServiceURLGetter interface {
	GetURL(ctx context.Context, shortURL string, getter shortlinks.URLGetter) (string, error)
}

func New(log zerolog.Logger, getter shortlinks.URLGetter, service ServiceURLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.geturl.New"

		log = log.With().
			Str("op", op).
			Str("request_id", middleware.GetReqID(r.Context())).
			Logger()

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error().Str("request body is empty", err.Error())
			render.JSON(w, r, resp.Error("empty request"))
			return
		}
		if err != nil {
			log.Error().Str("failed to decode request body", err.Error())
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}
		log.Info().Msg("request body decoded")

		shortURL := req.ShortURL
		if shortURL == "" {
			render.JSON(w, r, resp.Error("url is required"))
			log.Error().Msg("url is empty")
			return
		}

		longURL, err := service.GetURL(context.Background(), shortURL, getter)

		if err != nil {
			switch err.(type) {
			case customerrors.URLNotFound:
				render.JSON(w, r, resp.Error("url not found"))
			case customerrors.WrongURL:
				render.JSON(w, r, resp.Error("wrong url"))
			default:
				render.JSON(w, r, resp.Error("internal error"))
			}
			log.Error().Str("failed to get url", err.Error())
			return
		}

		render.JSON(w, r, ResponseOK(longURL))
		log.Info().Msg("url found and submitted")
	}
}

func ResponseOK(longURL string) Response {
	return Response{
		Response: resp.OK(),
		Url:      longURL,
	}
}
