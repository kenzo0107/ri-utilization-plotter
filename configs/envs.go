package configs

import (
	"log"

	"github.com/caarlos0/env"
	"github.com/pkg/errors"
)

// Envs : environment values
var Envs envParameters

type envParameters struct {
	DatadogAPIKeyName string `env:"DD_API_KEY_NAME" envDefault:"datadog_api_key"`
	DatadogAppKeyName string `env:"DD_APP_KEY_NAME" envDefault:"datadog_app_key"`
	Profile           string `env:"PROFILE"`
	AWSRegionID       string `env:"AWS_REGION"`
}

func init() {
	cfg := envParameters{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(errors.Wrap(err, "failed on env.Parse(&cfg)"))
	}
	Envs = cfg
}
