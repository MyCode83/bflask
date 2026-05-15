package logging

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

func New(verbose, quiet, json bool) zerolog.Logger {
	level := zerolog.InfoLevel
	if verbose {
		level = zerolog.DebugLevel
	}
	if quiet {
		level = zerolog.ErrorLevel
	}

	zerolog.SetGlobalLevel(level)

	var w io.Writer = os.Stderr
	if !json {
		w = zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
			NoColor:    quiet,
			FormatLevel: func(i interface{}) string {
				if i == nil {
					return "[???]"
				}
				level := i.(string)
				switch level {
				case "info":
					return "[INF]"
				case "debug":
					return "[DBG]"
				case "error":
					return "[ERR]"
				case "warn":
					return "[WRN]"
				default:
					return "[" + level + "]"
				}
			},
			PartsExclude: []string{zerolog.TimestampFieldName},
		}
	}

	return zerolog.New(w).With().Timestamp().Logger()
}
