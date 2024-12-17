package controller

import (
	"net/http"
	"os"

	"message-service-kata/internal/app/service"
	"message-service-kata/pkg/domain/entities"
	"message-service-kata/pkg/domain/response"
	"message-service-kata/pkg/validator"

	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

type (
	// MessageCtrl - controller interfacing for Message
	MessageCtrl interface {
		PostMessage(c echo.Context) error
		Health(c echo.Context) error
	}

	// MessageCtrlImpl - Implement service / usecase in Message controller
	MessageCtrlImpl struct {
		dig.In
		MessageSvc service.MessageSvc
	}
)

// NewMessageCtrl - Message controller instance
func NewMessageCtrl(impl MessageCtrlImpl) MessageCtrl {
	return &impl
}

// PostMessage handler to post message
func (r *MessageCtrlImpl) PostMessage(c echo.Context) error {
	var (
		req entities.CreateMessageRequest
		ctx = c.Request().Context()
	)

	err := c.Bind(&req)
	if err != nil {
		return response.ErrUnprocessableEntity.WithInternal(err)
	}

	err = validator.Validate(req)
	if err != nil {
		return response.ErrBadRequest.WithInternal(err)
	}

	err = r.MessageSvc.PostMessage(ctx, &req)
	if err != nil {
		return response.ErrInternalServerError.WithInternal(err)
	}

	return c.JSON(http.StatusOK, response.HTTPResponse{
		Status:  http.StatusOK,
		Message: response.DefaultMessage,
	})
}

// Health handler to health svc
func (r *MessageCtrlImpl) Health(c echo.Context) error {
	type resp struct {
		Code          int    `json:"code"`
		Version       string `json:"version"`
		Environment   string `json:"environment"`
		BuildCommitID string `json:"build_commit_id"`
		BuidTimeStamp string `json:"build_timestamp"`
	}

	response := resp{
		Code:          http.StatusOK,
		Version:       os.Getenv("APP_VERSION"),
		Environment:   os.Getenv("APP_BUILD_ENV"),
		BuildCommitID: os.Getenv("APP_BUILD_COMMIT_ID"),
		BuidTimeStamp: os.Getenv("APP_BUILD_TIMESTAMP"),
	}

	return c.JSON(http.StatusOK, response)
}
