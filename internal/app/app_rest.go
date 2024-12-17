package app

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"message-service-kata/internal/app/infra"
	"message-service-kata/pkg/di"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

var exitSigs = []os.Signal{syscall.SIGTERM, syscall.SIGINT}

// StartRestServer - function to serve application with graceful shutdown
func StartRestServer() {
	exitCh := make(chan os.Signal, 1)
	shutdownCh := make(chan struct{})
	errCh := make(chan error) // error channel

	signal.Notify(exitCh, exitSigs...)

	go func() {
		defer func() { exitCh <- syscall.SIGTERM }()
		if err := di.Invoke(startRestApp); err != nil {
			log.Error().Msgf("Invoke: %s", err.Error())
		}
	}()

	select {
	case <-exitCh:
		log.Info().Msg("exit signal received")
	case err := <-errCh:
		log.Error().Msgf("error received: %s", err.Error())
	}

	close(shutdownCh)

	if err := di.Invoke(gracefulRestShutdown); err != nil {
		log.Error().Msgf("Invoke: %s", err.Error())
	}
}

func startRestApp(
	e *echo.Echo,
	eCfg *infra.AppCfg,
) error {
	if err := di.Invoke(setRoute); err != nil {
		return err
	}

	return e.StartServer(&http.Server{
		Addr:         eCfg.Address,
		ReadTimeout:  eCfg.ReadTimeout,
		WriteTimeout: eCfg.WriteTimeout,
	})
}

// StartConsumerServer - function to serve application with graceful shutdown
func StartConsumerServer(topic string) {
	exitCh := make(chan os.Signal, 1)
	shutdownCh := make(chan struct{})
	errCh := make(chan error) // error channel

	signal.Notify(exitCh, exitSigs...)

	go func() {
		defer func() { exitCh <- syscall.SIGTERM }()
		if err := di.Invoke(func(args ConsumerHandlerParams) {
			startConsumer(args, shutdownCh, errCh, topic)
		}); err != nil {
			log.Error().Msgf("Invoke: %s", err.Error())
		}
	}()

	select {
	case <-exitCh:
		log.Info().Msg("exit signal received")
	case err := <-errCh:
		log.Error().Msgf("error received: %s", err.Error())
	}

	close(shutdownCh)

	if err := di.Invoke(gracefulConsumerShutdown); err != nil {
		log.Error().Msgf("Invoke: %s", err.Error())
	}
}

func gracefulRestShutdown(
	e *echo.Echo,
	pg *sql.DB,
) {
	timeOutTime := 60 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeOutTime)
	defer cancel()

	log.Info().Msg("shutting down rest server")

	if err := pg.Close(); err != nil {
		log.Error().Msgf("postgres close: %s", err.Error())
	}

	if err := e.Shutdown(ctx); err != nil {
		log.Error().Msgf("echo shutdown: %s", err.Error())
	}

	log.Info().Msg("rest server gracefully stopped")
}

func gracefulConsumerShutdown(
	pg *sql.DB,
	consumer *kafka.Consumer,
) {
	log.Info().Msg("shutting down consumer server")

	if err := consumer.Close(); err != nil {
		log.Error().Msgf("consumer.Close: %s", err.Error())
	}

	if err := pg.Close(); err != nil {
		log.Error().Msgf("pg.Close: %s", err.Error())
	}

	log.Info().Msg("consumer server gracefully stopped")
}
