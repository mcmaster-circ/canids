// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package grpcservice provides gRPC streaming ingestion services.
package grpcservice

import (
	"context"

	"github.com/mcmaster-circ/canids-v2/backend/libraries/uuid"
	"github.com/mcmaster-circ/canids-v2/protocol"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Register implements protocol.Register interface. If the incoming request is
// valid, it responds with an ACK and session identifer to start data stream.
// Returns error if register request is invalid.
func (s *IngestServer) Register(ctx context.Context, in *protocol.RegisterRequest) (*protocol.RegisterResponse, error) {
	// validate header, do not validate session field
	err := validateHeader(in.Header, false)
	if err != nil {
		return nil, err
	}

	// generate session identifier
	session := uuid.Generate()

	// respond with timeout data
	return &protocol.RegisterResponse{
		Header: &protocol.Header{
			MsgUuid:      uuid.Generate(),
			MsgTimestamp: timestamppb.Now(),
			Status:       protocol.Status_ACK,
			ErrorMsg:     "",
			Session:      session,
		},
		Timeout: int32(connTimeout.Seconds()),
	}, nil
}
