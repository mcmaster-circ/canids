// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package email provides email service through SendGrid.
package email

import (
	"strings"

	"github.com/mcmaster-circ/canids-v2/backend/state"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendPasswordReset is for sending out password reset emails. It accepts a
// state, the recepient's name, recepient's email, and a reset URL. It may
// return an error if the reset email cannot be sent.
func SendPasswordReset(s *state.State, name, email, url string) error {
	// prepare message
	from := mail.NewEmail(s.Config.SendGridName, s.Config.SendGridEmail)
	subject := "Password Reset"
	to := mail.NewEmail(name, email)

	// use reset email template
	htmlContent := resetEmail

	// update content with specific information
	htmlContent = strings.Replace(htmlContent, "#NAME", name, -1)
	htmlContent = strings.Replace(htmlContent, "#APPLICATION", s.Config.SendGridName, -1)
	htmlContent = strings.Replace(htmlContent, "#EMAIL", email, -1)
	htmlContent = strings.Replace(htmlContent, "#URL", url, -1)

	// plain text email template
	plainContent := resetEmailPlain

	// update plain text content with specific information
	plainContent = strings.Replace(plainContent, "#NAME", name, -1)
	plainContent = strings.Replace(plainContent, "#APPLICATION", s.Config.SendGridName, -1)
	plainContent = strings.Replace(plainContent, "#EMAIL", email, -1)
	plainContent = strings.Replace(plainContent, "#URL", url, -1)

	// send message
	message := mail.NewSingleEmail(from, subject, to, plainContent, htmlContent)
	err := s.SendEmail(message)
	return err
}
