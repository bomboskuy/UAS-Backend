package utils

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var Log zerolog.Logger


func InitLogger() {
	zerolog.TimeFieldFormat = time.RFC3339

	log.Logger = zerolog.New(os.Stdout).
		With().
		Timestamp().
		Logger()
}
