package main

import "github.com/kelseyhightower/envconfig"

type config struct {
	ServeRESTAddress string `envconfig:"service_rest_address" default:":8000"`
	DatabaseUrl      string `envconfig:"database_url"`
}

func ParseEnv() (*config, error) {
	c := new(config)
	if err := envconfig.Process("", c); err != nil {
		return nil, err
	}

	return c, nil
}
