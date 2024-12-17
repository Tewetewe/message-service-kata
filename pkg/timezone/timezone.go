package timezone

import (
	"time"
)

type (
	// Cfg used to load timezone config from .env
	Cfg struct {
		Location string `envconfig:"LOCATION" required:"true" default:"Asia/Jakarta"`
	}
)

// NewTimezone initialize application timezone
func NewTimezone(cfg *Cfg) error {
	local, err := time.LoadLocation(cfg.Location)
	if err != nil {
		return err
	}
	time.Local = local

	return nil
}
