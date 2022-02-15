// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package grpcservice provides gRPC streaming ingestion services.
package grpcservice

import (
	"errors"

	"github.com/mcmaster-circ/canids-v2/protocol"
)

// validateHeader reutrns an error if the header is not defined, or required
// fields within the header are not defined. validateSession indicates if header
// validates session field.
func validateHeader(header *protocol.Header, validateSesion bool) error {
	if header == nil {
		return errors.New("invalid header provided")
	}
	if header.MsgUuid == "" {
		return errors.New("invalid header UUID")
	}
	if header.MsgTimestamp.AsTime().IsZero() {
		return errors.New("invalid header timestamp")
	}
	if header.Status == protocol.Status_NULL_STATUS {
		return errors.New("invalid header status")
	}
	if validateSesion && header.Session == "" {
		return errors.New("invalid session")
	}
	return nil
}
