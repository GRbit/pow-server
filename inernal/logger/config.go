package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Config(logLevel string, console bool) {
	log.Logger = log.
		Level(zerolog.DebugLevel).
		With().Timestamp().
		Logger()

	switch logLevel {
	case "debug":
		log.Logger = log.Level(zerolog.DebugLevel)
	case "info":
		log.Logger = log.Level(zerolog.InfoLevel)
	case "warn":
		log.Logger = log.Level(zerolog.WarnLevel)
	case "error":
		log.Logger = log.Level(zerolog.ErrorLevel)
	}

	if console {
		log.Logger = log.
			Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.StampMicro}).
			With().Caller().Timestamp().Logger()
		zerolog.TimeFieldFormat = time.StampMicro
	}

	log.Debug().Msg("Logger configured")
}
