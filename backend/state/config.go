// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package state provides application state for the backend.
package state

import (
	"errors"
	"os"
)

// Config is the environment variable configuration for the backend.
type Config struct {
	ElasticHost string // ElasticHost is the hostname of Elasticsearch
	ElasticPort string // ElasticPort is the port of Elasticsearch
}

// load will attempt to load the required environment variables into the Config
// struct. An error will be returned if a required variable is not defined.
func (c *Config) load() error {
	// Elasticsearch parameters
	c.ElasticHost = os.Getenv("ELASTIC_HOST")
	if c.ElasticHost == "" {
		return errors.New("env ELASTIC_HOST not defined")
	}
	c.ElasticPort = os.Getenv("ELASTIC_PORT")
	if c.ElasticHost == "" {
		return errors.New("env ELASTIC_PORT not defined")
	}

	return nil
}
