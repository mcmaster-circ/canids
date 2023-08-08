package engine

import (
	"context"
	"log"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func ConnectWebsocketServer(s *state, db *database, endpoint string) error {
	// Attempt connection to server
	url := "ws://" + endpoint + "websocket"
	log.Printf("[CanIDS] attempting connection to %s\n", url)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		log.Printf("[CanIDS] failed to establish connection. %s. retrying in %s\n", err, s.RetryDelay)
		return err
	}
	// Defer closure for client exit
	defer conn.Close(websocket.StatusInternalError, "WebSocket closed")

	log.Println("[CanIDS] successful connection, performing self-registration")

	// Start period poll of file system for new files and stale files
	go fsPollingLoop(s, db)

	// Start file scanner
	for {
		// Get next frame, generate JSON payload
		frame, err := scannerGetFrame(s, db)
		if err != nil {
			log.Println("[CanIDS] failed to generate frame", err)
			continue
		}

		// Send frame to WebSocket server
		err = wsjson.Write(context.Background(), conn, frame)
		if err != nil {
			log.Println("[CanIDS] failed to send frame over WebSocket", err)
			log.Println("[CanIDS] retrying in", s.RetryDelay)
			close(s.PollingAbort)
			conn.Close(websocket.StatusInternalError, "WebSocket closed")
			return err
		}

		if s.Debug {
			log.Printf("[CanIDS] successful frame sent: %+v\n", frame)
		}
	}
}
