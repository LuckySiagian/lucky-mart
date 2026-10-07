package logger

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

func InitLogger(env string) {
	zerolog.TimeFieldFormat = time.RFC3339

	var level zerolog.Level
	switch strings.ToLower(env) {
	case "production":
		level = zerolog.WarnLevel
	case "staging":
		level = zerolog.InfoLevel
	default:
		level = zerolog.DebugLevel
	}

	zerolog.SetGlobalLevel(level)

	if strings.ToLower(env) == "development" {
		Log = zerolog.New(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "15:04:05",
		}).With().Timestamp().Caller().Logger()
	} else {
		Log = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}

	Log.Info().Msgf("Logger initialized with level: %s [Env: %s]", level.String(), env)
}
