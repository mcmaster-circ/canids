package engine

import (
	"context"
	"log"
	"net/http"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

const key = "hello"

type RecievedMsg struct {
	MsgType int `json:"type,omitempty"` // Message type: 0 - Misc, 1 - Ping
}

type MessageChannels struct {
	pingQueue chan *RecievedMsg
}

var queues = &MessageChannels{
	pingQueue: make(chan *RecievedMsg, 10000),
}

func ConnectWebsocketServer(s *state, db *database, endpoint string) error {
	// Attempt connection to server
	url := "http://host.docker.internal:6060/websocket/"
	log.Printf("[CanIDS] attempting connection to %s\n", url)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	dialOptions := websocket.DialOptions{
		HTTPHeader: http.Header{},
	}

	dialOptions.HTTPHeader.Set("Authorization", key)

	conn, _, err := websocket.Dial(ctx, url, &dialOptions)
	if err != nil {
		log.Printf("[CanIDS] failed to establish connection. %s. retrying in %s\n", err, s.RetryDelay)
		return err
	}
	// Defer closure for client exit
	defer conn.Close(websocket.StatusInternalError, "WebSocket closed")

	log.Println("[CanIDS] successful connection")

	// Start period poll of file system for new files and stale files
	go fsPollingLoop(s, db)

	go wsReader(s, conn)

	// Start file scanner
	for {

		var frame *UploadRequest
		select {
		case <-queues.pingQueue:
			log.Println("Item in ping queue")

			frame = generatePongFrame(s)
		default:
			log.Println("No ping")
			// Get next frame, generate JSON payload
			frame, err = scannerGetFrame(s, db)
			if err != nil {
				log.Println("[CanIDS] failed to generate frame", err)
				continue
			}
		}

		ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
		// Send frame to WebSocket server
		err = wsjson.Write(ctx, conn, frame)
		if err != nil {
			log.Println("[CanIDS] failed to send frame over WebSocket", err)
			log.Println("[CanIDS] retrying in", s.RetryDelay)
			close(s.PollingAbort)
			conn.Close(websocket.StatusInternalError, "WebSocket closed")
			cancel()
			return err
		}
		cancel()

		//log.Printf("[CanIDS] successful frame sent")
		// if s.Debug {
		// 	log.Printf("[CanIDS] successful frame sent: %+v\n", frame)
		// }

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

func wsReader(s *state, conn *websocket.Conn) {
	for {
		var msg RecievedMsg
		err := wsjson.Read(context.Background(), conn, &msg)
		if err != nil {
			log.Println("Invalid message reived")
			break
		}

		if msg.MsgType == 1 {
			queues.pingQueue <- &msg
		}
	}
}
