package infra

import (
	"time"

	"message-service-kata/pkg/domain/response"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

type (
	// AppCfg is echo application configs
	AppCfg struct {
		Address        string        `envconfig:"ADDRESS" default:":8089" required:"true"`
		ReadTimeout    time.Duration `envconfig:"READ_TIMEOUT" default:"5s"`
		WriteTimeout   time.Duration `envconfig:"WRITE_TIMEOUT" default:"10s"`
		Debug          bool          `envconfig:"DEBUG" default:"true"`
		Version        string        `envconfig:"VERSION" default:"v1"`
		BuildEnv       string        `envconfig:"BUILD_ENV" default:"local"`
		BuildCommitID  string        `envconfig:"BUILD_COMMIT_ID" default:"local"`
		BuildTimestamp string        `envconfig:"BUILD_TIMESTAMP" default:"local"`
	}
)

// NewEcho used to create echo instance
func NewEcho(cfg *AppCfg) *echo.Echo {
	e := echo.New()
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORS())
	e.Use(echoMiddleware.Gzip())
	e.Use(echoMiddleware.RequestID())

	e.HideBanner = true
	e.Debug = cfg.Debug
	e.HTTPErrorHandler = response.DefaultHTTPErrorHandler
	return e
}
