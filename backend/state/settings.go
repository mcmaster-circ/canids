// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package settings provides the simplified settings interface.
package state

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ainsleyclark/go-mail/drivers"
	"github.com/ainsleyclark/go-mail/mail"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/refresh"
	log "github.com/sirupsen/logrus"
)

const (
	indexConfiguration = "configuration"
)

// Settings is the mutable settings for the backend.
type Settings struct {
	EmailService string // EmailService is the email service used for sending emails
	// EmailAPIKey      string // EmailAPIKey is the service API key for email authentication
	// EmailFromAddress string // EmailFromAddress is the address emails are sent from
	// EmailName        string // EmailName is the name of the account sending emails
	// EmailDomain      string // EmailDomain is the domain used for password resets
	EmailConfig mail.Config

	MiddlewareDisable bool // MiddlewareDisable indicates if middleware is to be disabled
	HTTPSEnabled      bool // HTTPSEnabled indicates if the site is accessible over HTTPS

	UserRegistration bool // UserRegistration indicates if registration link is shown on login page
	UserActivated    bool // UserActivated indicates if user is automatically activated after registration

	DebugLogging bool // DebugLogging indicates if debug logging should be performed
}

type Service int64

const (
	None Service = iota
	Mailgun
	SendGrid
	Postal
	Postmark
	SparkPost
)

// Initialize settings in configuration index
func (settings *Settings) init(s *State) error {
	exists, err := s.Elastic.Indices.Exists(indexConfiguration).Do(s.ElasticCtx)
	if err != nil {
		return err
	}

	if !exists {
		defaultSettings := []DocumentSetting{
			{"MAIL_SERVICE", "NONE", false},
			{"MAIL_URL", "", false},
			{"MAIL_API_KEY", "", false},
			{"MAIL_FROM_ADDRESS", "", false},
			{"MAIL_FROM_NAME", "", false},
			{"MAIL_DOMAIN", "", false},
			{"MIDDLEWARE_DISABLE", "false", true},
			{"HTTPS_ENABLED", "false", true},
			{"USER_REGISTRATION", "true", false},
			{"USER_ACTIVATED", "false", false},
			{"DEBUG_LOGGING", "true", true},
		}

		s.Elastic.Indices.Create(indexConfiguration).Do(s.ElasticCtx)
		for _, setting := range defaultSettings {
			_, err := setting.index(s)
			if err != nil {
				s.Log.Error("error indexing new setting ", err)
				return nil
			}
		}
	}

	err = settings.load(s)
	if err != nil {
		return err
	}

	// set settings
	s.Settings = settings

	return nil
}

func (settings *Settings) load(s *State) error {
	allSettings, err := AllSettings(s)
	if err != nil {
		return err
	}

	for _, setting := range allSettings {
		refreshSetting(s, setting.Name, setting.Value)
	}

	return nil
}

// AllSettings will attempt to query the "configuration" index and return all settings in the
// system. It may return an error if the query cannot be completed.
func AllSettings(s *State) ([]DocumentSetting, error) {
	var out []DocumentSetting
	client, ctx := s.Elastic, s.ElasticCtx

	// perform query for all documents
	results, err := client.Search().Index(indexConfiguration).Query(&types.Query{
		MatchAll: types.NewMatchAllQuery(),
	}).Size(1000).Do(ctx)
	if err != nil {
		return nil, err
	}
	// parse settings into DocumentSetting, append to out
	for _, setting := range results.Hits.Hits {
		var d DocumentSetting
		err := json.Unmarshal(setting.Source_, &d)
		if err != nil {
			return nil, err
		}

		out = append(out, d)
	}
	return out, nil
}

// QueryDashboardByName will attempt to query the "dashboard" index for the
// dashboard with the specified UUID. It may return an error if the query cannot
// be completed.
func querySettingByName(s *State, name string) (DocumentSetting, string, error) {
	var d DocumentSetting
	client, ctx := s.Elastic, s.ElasticCtx

	// perform query for dashboard with provided uuid
	result, err := client.Search().Index(indexConfiguration).Query(&types.Query{
		Term: map[string]types.TermQuery{
			"name.keyword": {Value: name},
		},
	}).Do(ctx)
	if err != nil {
		return d, "", err
	}
	// ensure dashboard was returned
	if result.Hits.Total.Value == 0 {
		return d, "", errors.New("dashboard: no document with uuid found")
	}
	// select + parse dashboard into DocumentDashboard
	dashboard := result.Hits.Hits[0]
	err = json.Unmarshal(dashboard.Source_, &d)
	if err != nil {
		return d, "", err
	}
	// successful query
	return d, dashboard.Id_, nil
}

func UpdateSettings(s *State, changedSettings []DocumentSetting) error {
	for _, setting := range changedSettings {
		fmt.Println(setting.Name, setting.Value)
		_, esDocID, err := querySettingByName(s, setting.Name)
		if err != nil {
			return err
		}
		setting.update(s, esDocID)
		refreshSetting(s, setting.Name, setting.Value)
	}

	return nil
}

func setMailDriver(s *State) error {
	var err error
	switch s.Settings.EmailService {
	case "NONE":
		s.Mailer = nil
	case "MAILGUN":
		s.Mailer, err = drivers.NewMailgun(s.Settings.EmailConfig)
	case "SENDGRID":
		s.Mailer, err = drivers.NewSendGrid(s.Settings.EmailConfig)
	case "POSTAL":
		s.Mailer, err = drivers.NewPostal(s.Settings.EmailConfig)
	case "POSTMARK":
		s.Mailer, err = drivers.NewPostmark(s.Settings.EmailConfig)
	case "SPARKPOST":
		s.Mailer, err = drivers.NewSparkPost(s.Settings.EmailConfig)
	default:
		s.Log.Debug("Service not available")
	}
	return err
}

func refreshSetting(s *State, name string, value string) error {
	switch name {
	case "MAIL_SERVICE":
		s.Settings.EmailService = value
		setMailDriver(s)
	case "MAIL_URL":
		s.Settings.EmailConfig.URL = value
	case "MAIL_API_KEY":
		s.Settings.EmailConfig.APIKey = value
	case "MAIL_FROM_ADDRESS":
		s.Settings.EmailConfig.FromAddress = value
	case "MAIL_FROM_NAME":
		s.Settings.EmailConfig.FromName = value
	case "MAIL_DOMAIN":
		s.Settings.EmailConfig.Domain = value
	case "MIDDLEWARE_DISABLE":
		s.Settings.MiddlewareDisable = value == "true"
	case "HTTPS_ENABLED":
		s.Settings.HTTPSEnabled = value == "true"
	case "USER_REGISTRATION":
		s.Settings.UserRegistration = value == "true"
	case "USER_ACTIVATED":
		s.Settings.UserActivated = value == "true"
	case "DEBUG_LOGGING":
		s.Settings.DebugLogging = value == "true"
		if s.Settings.DebugLogging {
			s.Log.SetLevel(log.DebugLevel)
		} else {
			s.Log.SetLevel(log.InfoLevel)
		}
	}
	return nil
}

// DocumentSetting represents a document from the "configuration" index.
type DocumentSetting struct {
	// UUID       string `json:"uuid"`       // UUID is the unique setting identifier
	Name       string `json:"name"`       // Name is the setting display name
	Value      string `json:"value"`      // Value is the setting value
	IsAdvanced bool   `json:"isAdvanced"` // IsAdvanced is true for an advanced setting
}

// Index will attempt to index the document to the "configuration" index. It will return
// the newly created document ID or an error.
func (d *DocumentSetting) index(s *State) (string, error) {
	client, ctx := s.Elastic, s.ElasticCtx
	result, err := client.Index(indexConfiguration).Document(d).Refresh(refresh.True).Do(ctx)
	return result.Id_, err
}

// Update will attempt to update the document in the "configuration" with the provided
// Elasticsearch document ID. It will return an error if the transaction can not
// be performed.
func (d *DocumentSetting) update(s *State, esDocID string) error {
	client, ctx := s.Elastic, s.ElasticCtx

	_, err := client.Update(indexConfiguration, esDocID).Doc(
		map[string]interface{}{
			// "uuid":  d.UUID,
			"value": d.Value,
		}).DetectNoop(true).Refresh(refresh.True).Do(ctx)

	return err
}
