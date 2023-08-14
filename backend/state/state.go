// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package state provides application state for the backend.
package state

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/mcmaster-circ/canids-v2/backend/libraries/ipsetmgr"
	"github.com/olivere/elastic/v7"
	"github.com/oschwald/geoip2-golang"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	log "github.com/sirupsen/logrus"
)

const (
	esIdleConns   = 32               // max number of idle connections
	esIdleTimeout = 90 * time.Second // idle connection timeout
	esDialTimeout = 5 * time.Second  // connection establishment timeout
	esKeepAlive   = 30 * time.Second // keep alive duration
)

var (
	// esIndexes are indexes to create
	esIndexes = []string{
		"dashboard",
		"view",
	}
)

// State is the global state for the backend.
type State struct {
	Hash       string           // Hash is the hash of the latest Git commit
	Log        *log.Logger      // Log is a structured event logger
	Start      time.Time        // Start is the start time of the backend
	IsDocker   bool             // IsDocker indicates if running inside Docker
	Config     *Config          // Config contains global configuration
	Elastic    *elastic.Client  // Elastic is the Elasticsearch client state
	ElasticCtx context.Context  // ElasticCtx is the Elasticsearch client context
	AuthReady  bool             // AuthReady is if the "auth" index exists
	SendGrid   *sendgrid.Client // SendGrid is the email client

	GeoIPASN     *geoip2.Reader // GeoIPASN is the GeoIP ASN database
	GeoIPCity    *geoip2.Reader // GeoIPCity is the geoIP City database
	GeoIPCountry *geoip2.Reader // GeoIPCountry is the GeoIP country database

	AlarmManager *ipsetmgr.IPSetsManager // AlarmManager contains the ip lists that will trigger an alarm
}

// Provision will attempt to generate and return a new State. It will return an
// error if State fails to fully initialize.
func Provision(gitHash string) (*State, error) {
	// empty state
	var s State

	// generate Hash in State
	s.Hash = gitHash
	if s.Hash == "" {
		s.Hash = "unknown"
	}

	// generate Log in State
	s.Log = log.New()
	s.Log.Infof("[state] initializing backend, build %s", s.Hash)

	// generate Start in State
	s.Start = time.Now().UTC()

	// generate Config in State
	s.Log.Info("[state] initializing config in state")
	err := s.config()
	if err != nil {
		return &s, err
	}

	/*
		// generate Elasticsearch in State
		s.Log.Info("[state] initializing elasticsearch in state")
		err = s.elasticsearch()
		if err != nil {
			return &s, err
		}

		// generate AuthReady in State
		if err := checkAuthReady(&s); err != nil {
			s.Log.Error("[state] failed to check if Elasticsearch 'auth' index exists")
			return &s, err
		}
	*/

	// generate SendGrid in State
	if s.Config.SendGridToken != "" {
		s.SendGrid = sendgrid.NewSendClient(s.Config.SendGridToken)
	}

	// generate GeoIP entryes in state
	s.Log.Info("[state] initializing GeoIPASN, GeoIPCity, GeoIPCountry in state")
	err = s.geoIP()
	if err != nil {
		return &s, err
	}

	// put alarm manager in state
	s.AlarmManager = ipsetmgr.NewIPSetsManager()

	s.Log.Info("[state] backend initialized, no errors")
	return &s, nil
}

// config attempts to populate the Config field in State.
func (s *State) config() error {
	// if hostname is not populated, not running in Docker container
	s.IsDocker = os.Getenv("HOSTNAME") != ""

	// .env need to be manually loaded outside of Docker
	if !s.IsDocker {
		s.Log.Info("[state] backend not connected to Docker network, loading config files manually")
		configPath, err := filepath.Abs("../config/config.env")
		if err != nil {
			return err
		}
		secretPath, err := filepath.Abs("../config/secret.env")
		if err != nil {
			return err
		}
		err = godotenv.Load(configPath, secretPath)
		if err != nil {
			return err
		}
	} else {
		s.Log.Info("[state] backend connected to Docker network, autoloading config")
	}

	// generate configuration from loaded environment
	s.Config = &Config{}
	return s.Config.load()
}

// elasticsearch attempts to populate the Elastic and ElasticCtx fields in
// State. It also creates required Elasticsearch indexes.
func (s *State) elasticsearch() error {
	// determine if localhost is required
	host := s.Config.ElasticHost
	if !s.IsDocker {
		s.Log.Info("[state] backend not connected to Docker network, using loopback interface for Elasticsearch")
		host = "127.0.0.1"
	}
	s.Config.ElasticHost = host

	// generate context
	s.ElasticCtx = context.Background()

	// elasticsearch transport mechanism
	httpTransport := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:       esIdleConns,
			IdleConnTimeout:    esIdleTimeout,
			DisableCompression: true,
			Dial: (&net.Dialer{
				Timeout:   esDialTimeout,
				KeepAlive: esKeepAlive,
			}).Dial,
		},
	}

	// connect to Elasticsearch
	esURI := fmt.Sprintf("http://%s:%s", s.Config.ElasticHost, s.Config.ElasticPort)
	client, err := elastic.NewSimpleClient(
		elastic.SetURL(esURI),
		elastic.SetHttpClient(httpTransport),
	)
	if err != nil {
		return err
	}

	flag := true
	for flag {
		// ping to ensure connectivity
		_, _, err = client.Ping(esURI).Do(s.ElasticCtx)
		if err != nil {
			// give warning, Elastcsearch is probably still starting up
			s.Log.Warn("[state] elasticsearch ping failed, database is starting or incorrect parameters, sleeping 30 seconds...")
			time.Sleep(30 * time.Second)
		} else {
			flag = false
		}
		s.Elastic = client
	}

	// create indexes
	for _, index := range esIndexes {
		s.Elastic.CreateIndex(index).Do(s.ElasticCtx)
	}

	return nil
}

// checkAuthReady attempts to populate the AuthReady field in State. It checks
// if the "auth" index is created. If the existence of the index can not be
// checked, an error will be returned.
func checkAuthReady(s *State) error {
	s.Log.Info("[state] checking if elasticsearch 'auth' index exists")
	exists, err := s.Elastic.IndexExists("auth").Do(s.ElasticCtx)
	if err != nil {
		return err
	}
	s.AuthReady = exists
	if !s.AuthReady {
		s.Log.Warn("elasticsearch 'auth' index does not exist, system not initialized")
	}
	return nil
}

// geoIP attempts to populate the GeoIPASN, GeoIPCity, and GeoIPCountry
// databases. It will return an error if the databases cannot be loaded.
func (s *State) geoIP() error {
	// ASN database
	fullPath, err := filepath.Abs("geoip/GeoLite2-ASN.mmdb")
	if err != nil {
		return err
	}
	db, err := geoip2.Open(fullPath)
	if err != nil {
		return err
	}
	s.GeoIPASN = db

	// city database
	fullPath, err = filepath.Abs("geoip/GeoLite2-City.mmdb")
	if err != nil {
		return err
	}
	db, err = geoip2.Open(fullPath)
	if err != nil {
		return err
	}
	s.GeoIPCity = db

	// country database
	fullPath, err = filepath.Abs("geoip/GeoLite2-Country.mmdb")
	if err != nil {
		return err
	}
	db, err = geoip2.Open(fullPath)
	if err != nil {
		return err
	}
	s.GeoIPCountry = db

	return nil
}

// SendEmail sends an email using sendgrid if a sendgrid token was provided,
// otherwise no-ops
func (s *State) SendEmail(message *mail.SGMailV3) error {
	if s.SendGrid != nil {
		_, err := s.SendGrid.Send(message)
		return err
	}

	return nil
}
