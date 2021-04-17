package main

import "github.com/kelseyhightower/envconfig"

const appId = "orderservice"

type config struct {
	ServeRESTAddress string `envconfig:"orderservice_service_rest_address" default:":8000"`
	DatabaseUrl string `envconfig:"database_url" default:"orderservice:1234@/orderservice"`
}

func ParseEnv() (*config, error) {
	c := new(config)
	if err := envconfig.Process(appId, c); err != nil {
		return nil, err
	}

	return c, nil
}