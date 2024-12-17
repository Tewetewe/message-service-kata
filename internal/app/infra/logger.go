package infra

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// InitLogger initializes the global logger with settings based on environment variables.
func InitLogger() {
	// Set default global log level to Info.
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// If APP_DEBUG is true, set log level to Debug.
	if os.Getenv("APP_DEBUG") == "true" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// For local environment, use a human-readable console writer.
	if os.Getenv("APP_BUILD_ENV") == "local" {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: "2006-01-02 15:04:05", // Set a human-readable time format
		})
	}

	log.Info().Msg("Logger initialized")
}
