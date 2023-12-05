package saveurl

import (
	"errors"
	"io"
	"net/http"

	resp "github.com/GrishaSkurikhin/OzonTestTask/internal/rest-server/api/response"
	"github.com/GrishaSkurikhin/OzonTestTask/internal/url"
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

func New(log zerolog.Logger, saver url.URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.saveurl.New"

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

		shortURL, err := url.SaveURL(req.LongURL, saver)
		if err != nil {
			log.Error().Str("failed to save url", err.Error())
			render.JSON(w, r, resp.Error("failed to save url"))
			return
		}

		render.JSON(w, r, ResponseOK(shortURL))
		log.Info().Msg("url saved successfully")
	}
}

func ResponseOK(shortURL string) Response {
	return Response{
		Response: resp.OK(),
		Url:      shortURL,
	}
}
