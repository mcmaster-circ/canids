// Copyright (c) 2020 Computing Infrastructure Research Centre (CIRC), McMaster
// University. All rights reserved.

package engine

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/mcmaster-circ/canids-v2/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	connectTimeout = 5 * time.Second
)

// Connect is responsible for establishing the grPC connection to the provided
// gRPC server. Debug indicates if debug logging should be used. It will return
// an error if the gRPC cannot be read.
func Connect(s *state, db *database, endpoint string) error {
	// attempt connection to server
	var conn *grpc.ClientConn
	var err error

	// dial with secure or insecure connection
	if s.Insecure {
		log.Printf("[CanIDS] attempting connection to %s (insecure)\n", endpoint)
		conn, err = grpc.Dial(endpoint,
			grpc.WithInsecure(),
			grpc.WithTimeout(connectTimeout),
		)
	} else {
		log.Printf("[CanIDS] attempting connection to %s\n", endpoint)

		var creds credentials.TransportCredentials
		creds, err = credentials.NewClientTLSFromFile("/run/secrets/canids_ca_cert", "")
		if err != nil {
			log.Printf("cannot load TLS credentials: ", err)
		}

		conn, err = grpc.Dial(endpoint,
			grpc.WithTransportCredentials(creds),
			grpc.WithTimeout(connectTimeout),
		)
	}

	if err != nil {
		log.Printf("[CanIDS] failed to establish connection. %s. retrying in %s\n", err, s.RetryDelay)
		return err
	}
	// defer closure for client exit
	defer conn.Close()

	log.Println("[CanIDS] successful connection, performing self-registration")

	client := protocol.NewWireServiceClient(conn)

	// send registration
	registerResp, err := client.Register(context.Background(), &protocol.RegisterRequest{
		Header: &protocol.Header{
			MsgUuid:      uuid.New().String(),
			MsgTimestamp: timestamppb.Now(),
			Status:       protocol.Status_REQUEST,
			ErrorMsg:     "",
			Session:      "",
		},
	})
	if err != nil {
		log.Println("[CanIDS] failed to perform self-registration", err.Error())
		log.Println("[CanIDS] retrying in", s.RetryDelay)
		return err
	}
	log.Printf("[CanIDS] successful registration %+v\n", registerResp)
	s.Session = registerResp.Header.Session

	// start period poll of file system for new files and stale files
	go fsPollingLoop(s, db)

	// start file scanner
	for {
		// get next frame, generate protobuf binary
		frame, err := scannerGetFrame(s, db)
		if err != nil {
			log.Println("[CanIDS] failed to generate frame", err)
			continue
		}

		// emitt frame
		s.NetworkMutex.Lock()
		uploadResponse, err := client.Upload(context.Background(), frame)
		s.NetworkMutex.Unlock()
		if err != nil {
			log.Println("[CanIDS] failed to upload frame", err)
			log.Println("[CanIDS] retrying in", s.RetryDelay)
			// error writing message, terminate fs polling loop + close
			// connection, terminate
			close(s.PollingAbort)
			conn.Close()
			return err
		}

		// ensure valid response
		if uploadResponse.Header.Status != protocol.Status_ACK {
			log.Println("[CanIDS] gRPC server error", uploadResponse.Header.ErrorMsg)
			log.Println("[CanIDS] retrying in", s.RetryDelay)
			// error writing message, terminate fs polling loop + close
			// connection, terminate
			close(s.PollingAbort)
			conn.Close()
			return err
		}

		if s.Debug {
			log.Printf("[CanIDS] successful upload %+v\n", uploadResponse)
		}
	}
}

// fsPollingLoop will perodically synchronize the local database for new/removed
// files in the specified directory.
func fsPollingLoop(s *state, db *database) {
	for {
		select {
		case <-s.PollingAbort:
			// received exit signal from event loop, terminate self
			return
		default:
			// sync the scanner to retreive latest database
			new, err := syncScanner(s)
			if err != nil {
				log.Println("[CanIDS] local database error:", err)
			}
			db.Next = new.Next
			db.Files = new.Files
			time.Sleep(s.FileScan)
		}
	}
}
