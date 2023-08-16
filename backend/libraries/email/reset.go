// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package email provides email service.
package email

import (
	"strings"

	"github.com/ainsleyclark/go-mail/mail"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

// SendPasswordReset is for sending out password reset emails. It accepts a
// state, the recepient's name, recepient's email, and a reset URL. It may
// return an error if the reset email cannot be sent.
func SendPasswordReset(s *state.State, name, email, url string) error {
	// use reset email template
	htmlContent := resetEmail

	// update content with specific information
	htmlContent = strings.Replace(htmlContent, "#NAME", name, -1)
	htmlContent = strings.Replace(htmlContent, "#APPLICATION", s.Settings.EmailConfig.FromName, -1)
	htmlContent = strings.Replace(htmlContent, "#EMAIL", email, -1)
	htmlContent = strings.Replace(htmlContent, "#URL", url, -1)

	// plain text email template
	plainContent := resetEmailPlain

	// update plain text content with specific information
	plainContent = strings.Replace(plainContent, "#NAME", name, -1)
	plainContent = strings.Replace(plainContent, "#APPLICATION", s.Settings.EmailConfig.FromName, -1)
	plainContent = strings.Replace(plainContent, "#EMAIL", email, -1)
	plainContent = strings.Replace(plainContent, "#URL", url, -1)

	// prepare message
	tx := &mail.Transmission{
		Recipients: []string{email},
		Subject:    "Password Reset",
		HTML:       htmlContent,
	}

	// send message
	_, err := s.Mailer.Send(tx)
	return err
}
