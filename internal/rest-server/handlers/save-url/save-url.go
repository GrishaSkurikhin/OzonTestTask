package saveurl

import (
	"context"
	"errors"
	"fmt"
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
	LongURL string `json:"longURL"`
}

type Response struct {
	resp.Response
	Url string `json:"shortURL"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=ServiceURLSaver
type ServiceURLSaver interface {
	SaveURL(ctx context.Context, longURL string, host string, saver shortlinks.URLSaver) (string, error)
}

func New(log zerolog.Logger, saver shortlinks.URLSaver, host string, service ServiceURLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.saveurl.New"

		log = log.With().
			Str("op", op).
			Str("request_id", middleware.GetReqID(r.Context())).
			Logger()

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error().Str(fmt.Sprintf("%s: request body is empty", op), err.Error())
			render.JSON(w, r, resp.Error("empty request"))
			return
		}
		if err != nil {
			log.Error().Str(fmt.Sprintf("%s: failed to decode request body", op), err.Error())
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}
		log.Info().Msg(fmt.Sprintf("%s: request body decoded", op))

		longURL := req.LongURL
		if longURL == "" {
			render.JSON(w, r, resp.Error("url is required"))
			log.Error().Msg(fmt.Sprintf("%s: url is empty", op))
			return
		}

		shortURL, err := service.SaveURL(context.Background(), req.LongURL, host, saver)
		if err != nil {
			switch err.(type) {
			case customerrors.WrongURL:
				render.JSON(w, r, resp.Error("wrong url"))
			default:
				render.JSON(w, r, resp.Error("internal error"))
			}
			log.Error().Str(fmt.Sprintf("%s: failed to get url", op), err.Error())
			return
		}

		render.JSON(w, r, ResponseOK(shortURL))
		log.Info().Msg(fmt.Sprintf("%s: url saved successfully", op))
	}
}

func ResponseOK(shortURL string) Response {
	return Response{
		Response: resp.OK(),
		Url:      shortURL,
	}
}
