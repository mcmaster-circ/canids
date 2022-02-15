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

// Upload implements protocol.Upload interface. If the incoming request is
// valid, it adds the upload request to the internal queue and responds with an
// ACK. Returns error if the upload request is invalid.
func (s *IngestServer) Upload(ctx context.Context, in *protocol.UploadRequest) (*protocol.UploadResponse, error) {
	// validate header, validate session field
	err := validateHeader(in.Header, true)
	if err != nil {
		return nil, err
	}

	// define log with session field
	l := s.state.Log.WithField("session", in.Header.Session)

	// validate asset ID
	if in.AssetId == "" {
		l.Warn("invalid asset ID provided")
		return &protocol.UploadResponse{
			Header: &protocol.Header{
				MsgUuid:      uuid.Generate(),
				MsgTimestamp: timestamppb.Now(),
				Status:       protocol.Status_NACK,
				ErrorMsg:     "Invalid asset ID provided.",
				Session:      in.Header.Session,
			},
		}, nil
	}

	// validate file name
	if in.FileName == "" {
		l.Warn("invalid file name provided")
		return &protocol.UploadResponse{
			Header: &protocol.Header{
				MsgUuid:      uuid.Generate(),
				MsgTimestamp: timestamppb.Now(),
				Status:       protocol.Status_NACK,
				ErrorMsg:     "Invalid file name provided.",
				Session:      in.Header.Session,
			},
		}, nil
	}

	// validate payloads were provided
	if len(in.Payload) == 0 {
		l.Warn("empty payload provided")
		return &protocol.UploadResponse{
			Header: &protocol.Header{
				MsgUuid:      uuid.Generate(),
				MsgTimestamp: timestamppb.Now(),
				Status:       protocol.Status_NACK,
				ErrorMsg:     "Invalid payload provided.",
				Session:      in.Header.Session,
			},
		}, nil
	}

	// add data into ingestion queue
	s.queue <- &Frame{
		AssetID:  in.AssetId,
		FileName: in.FileName,
		Payload:  in.Payload,
	}

	// ACK message
	return &protocol.UploadResponse{
		Header: &protocol.Header{
			MsgUuid:      uuid.Generate(),
			MsgTimestamp: timestamppb.Now(),
			Status:       protocol.Status_ACK,
			ErrorMsg:     "",
			Session:      in.Header.Session,
		},
	}, nil
}
