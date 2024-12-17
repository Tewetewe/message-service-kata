package app

import (
	"net/http"
	"os"

	controller "message-service-kata/internal/app/controller/rest"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

const (
	// ContextPath - Application API base path
	ContextPath = "/v1/message/"

	// PostMessage - Create New messages api path
	PostMessage = ContextPath + "post"

	// HealthPath - Application health check api path
	HealthPath = ContextPath + "health"
)

// setRoute - registering route to the application
func setRoute(
	e *echo.Echo,
	messageCtrl controller.MessageCtrl,
) {
	// Public API
	e.POST(PostMessage, messageCtrl.PostMessage)

	e.GET(HealthPath, messageCtrl.Health)
}

// health - application health api
func health(ec echo.Context) error {
	type resp struct {
		Code          int    `json:"code"`
		Version       string `json:"version"`
		Environment   string `json:"environment"`
		BuildCommitID string `json:"build_commit_id"`
		BuidTimeStamp string `json:"build_timestamp"`
	}

	var (
		ctx = ec.Request().Context()
	)

	response := resp{
		Code:          http.StatusOK,
		Version:       os.Getenv("APP_VERSION"),
		Environment:   os.Getenv("APP_BUILD_ENV"),
		BuildCommitID: os.Getenv("APP_BUILD_COMMIT_ID"),
		BuidTimeStamp: os.Getenv("APP_BUILD_TIMESTAMP"),
	}

	log.Info().Ctx(ctx).Msgf("health: %v", response)

	return ec.JSON(http.StatusOK, response)
}
