// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package state provides application state for the backend.
package state

import (
	"errors"
	"log"
	"os"
)

// Config is the environment variable configuration for the backend.
type Config struct {
	ElasticHost string // ElasticHost is the hostname of Elasticsearch
	ElasticPort string // ElasticPort is the port of Elasticsearch

	SendGridToken  string // SendGridToken is the SendGrid API token for email authentication
	SendGridEmail  string // SendGridEmail is the address emails are sent from
	SendGridName   string // SendGridName is the name of the account sending emails
	SendGridDomain string // SendGridDomain is the domain used for password resets

	MiddlewareDisable bool // MiddlewareDisable indicates if middleware is to be disabled
	HTTPSEnabled      bool // HTTPSEnabled indicates if the site is accessible over HTTPS

	UserRegistration bool // UserRegistration indicates if registration link is shown on login page
	UserActivated    bool // UserActivated indicates if user is automatically activated after registration

	DebugLogging bool // DebugLogging indicates if debug logging should be performed
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

	// SendGrid parameters
	c.SendGridToken = os.Getenv("SENDGRID_TOKEN")
	if c.SendGridToken == "" {
		log.Printf("env SENDGRID_TOKEN not defined, disabling email")
	}
	c.SendGridEmail = os.Getenv("SENDGRID_EMAIL")
	if c.SendGridEmail == "" {
		return errors.New("env SENDGRID_EMAIL not defined")
	}
	c.SendGridName = os.Getenv("SENDGRID_NAME")
	if c.SendGridName == "" {
		return errors.New("env SENDGRID_NAME not defined")
	}
	c.SendGridDomain = os.Getenv("SENDGRID_DOMAIN")
	if c.SendGridDomain == "" {
		return errors.New("env SENDGRID_DOMAIN not defined")
	}

	// middleware parameter
	c.MiddlewareDisable = os.Getenv("MIDDLEWARE_DISABLE") == "true"

	// HTTPS parameter
	c.HTTPSEnabled = os.Getenv("HTTPS_ENABLED") == "true"

	// user registration parameters
	c.UserRegistration = os.Getenv("USER_REGISTRATION") == "true"
	c.UserActivated = os.Getenv("USER_ACTIVATED") == "true"

	// debug parameter
	c.DebugLogging = os.Getenv("DEBUG_LOGGING") == "true"

	return nil
}
