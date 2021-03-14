package main

import "github.com/kelseyhightower/envconfig"

const appId = "orderService"

type config struct {
	ServeRESTAddress string `envconfig:"SERVICE_REST_ADDRESS" default:"8000"`
	DataBaseUrl string `envconfig:"DATABASE_URL"`
}

func ParseEnv() (*config, error) {
	c := new(config)
	if err := envconfig.Process(appId, c); err != nil {
		return nil, err
	}

	return c, nil
}