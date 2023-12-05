package logger

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

func New(env string) (*zerolog.Logger, error) {
	const op = "logger.New"

	var log zerolog.Logger
	switch env {
	case "local":
		log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).
			Level(zerolog.Level(zerolog.DebugLevel)).
			With().
			Logger()
	case "prod":
		log = zerolog.New(os.Stderr).
			Level(zerolog.Level(zerolog.InfoLevel)).
			With().
			Timestamp().
			Logger()
	default:
		return nil, fmt.Errorf("%s: %s", op, "unknown environment")
	}
	
	return &log, nil
}
