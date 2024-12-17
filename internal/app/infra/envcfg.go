package infra

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

const (
	// ServiceConsumerKafka variable service consumer kafka
	ServiceConsumerKafka = "consumer"
	// ServiceRestAPI variable service rest api
	ServiceRestAPI = "rest"
)

// AvailableServices initiate available server on this service
var AvailableServices = []string{
	ServiceRestAPI,
	ServiceConsumerKafka,
}

// LoadPgDatabaseCfg loading postgres database config using envconfig library
func LoadPgDatabaseCfg() (*DatabaseCfg, error) {
	var cfg DatabaseCfg
	prefix := "PG"
	if err := envconfig.Process(prefix, &cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", prefix, err)
	}

	return &cfg, nil
}

// LoadAppCfg loading echo application config using envconfig library
func LoadAppCfg() (*AppCfg, error) {
	var cfg AppCfg
	prefix := "APP"
	if err := envconfig.Process(prefix, &cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", prefix, err)
	}
	return &cfg, nil
}

// LoadKafkaCfg loading kafka config using envconfig library
func LoadKafkaCfg() (*KafkaCfg, error) {
	var cfg KafkaCfg
	prefix := "KAFKA"
	if err := envconfig.Process(prefix, &cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", prefix, err)
	}

	return &cfg, nil
}
