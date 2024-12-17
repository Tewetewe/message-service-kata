package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"

	"message-service-kata/internal/app"
	kafkaCtrl "message-service-kata/internal/app/controller/kafka"
	controller "message-service-kata/internal/app/controller/rest"
	"message-service-kata/internal/app/infra"
	"message-service-kata/internal/app/repo/kafka"
	"message-service-kata/internal/app/repo/postgres"
	"message-service-kata/internal/app/service"

	"github.com/joho/godotenv"

	"message-service-kata/pkg/di"
	"message-service-kata/pkg/timezone"
	"message-service-kata/pkg/validator"

	"golang.org/x/exp/slices"
)

var (
	topicFlag   = flag.String("topic", "", fmt.Sprintf("Kafka topic to bind. available topic: \n\t - %s", strings.Join(app.Topics, ",\n\t - ")))
	serviceFlag = flag.String("service", "", fmt.Sprintf("service to run. available services: \n\t - %s", strings.Join(infra.AvailableServices, ",\n\t - ")))
)

func main() {
	flag.Parse()

	if !slices.Contains(infra.AvailableServices, *serviceFlag) {
		fmt.Printf("unknown service: %s\n", *serviceFlag)
		fmt.Print("\n\n")
		flag.Usage()

		os.Exit(1)
	}

	if *serviceFlag == infra.ServiceConsumerKafka && *topicFlag != "" {
		if !slices.Contains(app.Topics, *topicFlag) {
			fmt.Printf("unknown topic: %s\n", *topicFlag)
			fmt.Print("\n\n")
			flag.Usage()

			os.Exit(1)
		}
	}

	serve(*serviceFlag, *topicFlag)
}

func serve(serviceNameFlag, topicNameFlag string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}

	infra.InitLogger()

	err = LoadApplicationConfig()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	var cfg timezone.Cfg
	if err = envconfig.Process("TIMEZONE", &cfg); err != nil {
		log.Fatal().Msg(err.Error())
	}

	if err = timezone.NewTimezone(&cfg); err != nil {
		log.Fatal().Msg(err.Error())
	}

	switch serviceNameFlag {
	case infra.ServiceRestAPI:
		err = LoadApplicationRestPackage()
	case infra.ServiceConsumerKafka:
		err = LoadApplicationKafkaPackage()
	}
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	err = LoadApplicationRepository()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	err = LoadApplicationService()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	err = LoadApplicationController()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	switch serviceNameFlag {
	case infra.ServiceRestAPI:
		app.StartRestServer()
	case infra.ServiceConsumerKafka:
		app.StartConsumerServer(topicNameFlag)
	}
}

// LoadApplicationConfig load application config
func LoadApplicationConfig() error {
	// infra
	err := di.Provide(infra.LoadAppCfg)
	if err != nil {
		return fmt.Errorf("LoadAppCfg: %s", err.Error())
	}

	err = di.Provide(infra.LoadPgDatabaseCfg)
	if err != nil {
		return fmt.Errorf("LoadPgDatabaseCfg: %s", err.Error())
	}

	err = di.Provide(infra.LoadKafkaCfg)
	if err != nil {
		return fmt.Errorf("LoadKafkaCfg: %s", err.Error())
	}

	return nil
}

// LoadApplicationRestPackage Load application package used by the rest application
//
//nolint:dupl
func LoadApplicationRestPackage() error {
	err := di.Provide(infra.NewEcho)
	if err != nil {
		return fmt.Errorf("NewEcho: %s", err.Error())
	}

	err = di.Provide(infra.NewDatabases)
	if err != nil {
		return fmt.Errorf("NewDatabases: %s", err.Error())
	}

	err = di.Provide(infra.NewProducer)
	if err != nil {
		return fmt.Errorf("NewProducer: %s", err.Error())
	}

	err = di.Invoke(validator.NewValidator)
	if err != nil {
		return fmt.Errorf("NewValidator: %s", err.Error())
	}

	return nil
}

// LoadApplicationKafkaPackage Load application package used by the rest application
//
//nolint:dupl
func LoadApplicationKafkaPackage() error {
	err := di.Provide(infra.NewDatabases)
	if err != nil {
		return fmt.Errorf("NewDatabases: %s", err.Error())
	}

	err = di.Provide(infra.NewConsumer)
	if err != nil {
		return fmt.Errorf("NewConsumer: %s", err.Error())
	}

	err = di.Provide(infra.NewProducer)
	if err != nil {
		return fmt.Errorf("NewProducer: %s", err.Error())
	}

	// controller
	err = di.Provide(kafkaCtrl.NewProcessor)
	if err != nil {
		return fmt.Errorf("NewKafkaController: %s", err.Error())
	}

	return nil
}

// LoadApplicationRepository load repository using ubed dig
//
//nolint:dupl
func LoadApplicationRepository() error {
	// repo kafka
	err := di.Provide(kafka.NewKafkaRepository)
	if err != nil {
		return fmt.Errorf("NewKafkaRepository: %s", err.Error())
	}

	// repo postgres
	err = di.Provide(postgres.NewMessageRepository)
	if err != nil {
		return fmt.Errorf("NewMessageRepository: %s", err.Error())
	}

	return nil
}

// LoadApplicationService Load service or usecase of the application using uber dig
func LoadApplicationService() error {
	// service
	err := di.Provide(service.NewMessageSvc)
	if err != nil {
		return fmt.Errorf("NewMessageSvc: %s", err.Error())
	}

	return nil
}

// LoadApplicationController load application controller using uber dig
func LoadApplicationController() error {
	// controller
	err := di.Provide(controller.NewMessageCtrl)
	if err != nil {
		return fmt.Errorf("NewMessageCtrl: %s", err.Error())
	}

	return nil
}
