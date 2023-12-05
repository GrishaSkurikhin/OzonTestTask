package geturl

import (
	"net/http"

	customerrors "github.com/GrishaSkurikhin/OzonTestTask/internal/custom-errors"
	resp "github.com/GrishaSkurikhin/OzonTestTask/internal/rest-server/api/response"
	"github.com/GrishaSkurikhin/OzonTestTask/internal/url"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
)

type Response struct {
	resp.Response
	Url string `json:"url"`
}

func New(log zerolog.Logger, getter url.URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.geturl.New"

		log = log.With().
			Str("op", op).
			Str("request_id", middleware.GetReqID(r.Context())).
			Logger()

		shortURL := r.URL.Query().Get("shortURL")
		longURL, err := url.GetURL(shortURL, getter)

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
