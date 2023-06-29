// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

// Package grpcservice provides gRPC streaming ingestion services.
package grpcservice

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/mcmaster-circ/canids-v2/backend/state"
	"github.com/mcmaster-circ/canids-v2/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	// "io/ioutil"
)

const (
	gRPCPort    = 50000
	connTimeout = 15 * time.Second
	bufferSize  = 4096
)

// IngestServer is gRPC server.
type IngestServer struct {
	state *state.State
	queue chan *Frame
	protocol.UnimplementedWireServiceServer
}

// Frame is a file upload with payloads (lines).
type Frame struct {
	AssetID  string   // AssetID is asset identifier
	FileName string   // FileName is upload file name
	Payload  [][]byte // Payload is a list of JSON byte payloads (lines)
}

// Provision registers and starts the gRPC service. Returns error if gRPC fails
// to register or start service.
func Provision(ctx context.Context, s *state.State) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", gRPCPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var creds credentials.TransportCredentials
	creds, err = credentials.NewServerTLSFromFile("/run/secrets/canids_server_cert", "/run/secrets/canids_server_key")

	if err != nil {
		log.Fatalf("failed to retrieve TLS credentials ", err)
	}

	registrar := grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(unaryInterceptor),
	)

	// register handlers as gRPC wire service server
	server := &IngestServer{
		state: s,
		queue: make(chan *Frame, bufferSize),
	}
	protocol.RegisterWireServiceServer(registrar, server)

	// start the server pipeline to consume from queue and index
	go server.pipeline()

	return registrar.Serve(listener)
}

// pipeline consumes from the internal queue to ingest chunks.
func (s *IngestServer) pipeline() {
	for {
		// ingest chunks from queue
		chunk := <-s.queue
		s.ingest(chunk)
	}
}

// unaryInterceptor implements unary server interceptor interface.
func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	log.Printf("request - Method:%s\tDuration:%s\tError:%v\n", info.FullMethod, time.Since(start), err)
	return resp, err
}
