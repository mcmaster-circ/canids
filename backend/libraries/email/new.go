// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package email provides email service.
package email

import (
	"strings"

	"github.com/ainsleyclark/go-mail/mail"
	"github.com/mcmaster-circ/canids-v2/backend/state"
)

// SendNewReset is for sending out a password set email to new users. It accepts
// a state, the recepient's name, recepient's email, set URL, sender's name,
// sender's email. It may return an error if the reset email cannot be sent.
func SendNewReset(s *state.State, name, email, url, senderName, senderEmail string) error {
	// use reset email template
	htmlContent := newEmail

	// update content with specific information
	htmlContent = strings.Replace(htmlContent, "#NAME", name, -1)
	htmlContent = strings.Replace(htmlContent, "#APPLICATION", s.Settings.EmailConfig.FromName, -1)
	htmlContent = strings.Replace(htmlContent, "#SENDERNAME", senderName, -1)
	htmlContent = strings.Replace(htmlContent, "#SENDEREMAIL", senderEmail, -1)
	htmlContent = strings.Replace(htmlContent, "#URL", url, -1)

	// plain text email template
	plainContent := newEmailPlain

	// update plain text content with specific information
	plainContent = strings.Replace(plainContent, "#NAME", name, -1)
	plainContent = strings.Replace(plainContent, "#APPLICATION", s.Settings.EmailConfig.FromName, -1)
	plainContent = strings.Replace(plainContent, "#SENDERNAME", senderName, -1)
	plainContent = strings.Replace(plainContent, "#SENDEREMAIL", senderEmail, -1)
	plainContent = strings.Replace(plainContent, "#URL", url, -1)

	// prepare message
	tx := &mail.Transmission{
		Recipients: []string{email},
		Subject:    "Activate Account",
		HTML:       htmlContent,
	}

	// send message
	_, err := s.Mailer.Send(tx)
	return err
}
