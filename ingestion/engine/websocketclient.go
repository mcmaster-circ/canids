package engine

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
	"net/http"

	// "net/http"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type Message struct {
	MsgType int    `json:"type,omitempty"` // Message type: 0 - Misc, 1 - Ping, 2 - connection success
	Msg     string `json:"msg,omitempty"`
}

type MessageChannels struct {
	pingQueue chan *Message
}

var queues = &MessageChannels{
	pingQueue: make(chan *Message, 10000),
}

func ConnectWebsocketServer(s *state, db *database, endpoint string) error {
	// Attempt connection to server

	log.Printf("[CanIDS] attempting connection to %s\n", endpoint)

	dialOptions := websocket.DialOptions{
		HTTPHeader: http.Header{},
	}

	dialOptions.HTTPHeader.Set("Authorization", s.AssetID)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	conn, _, err := websocket.Dial(ctx, endpoint, &dialOptions)
	if err != nil {
		log.Printf("[CanIDS] failed to establish connection. %s. retrying in %s\n", err, s.RetryDelay)
		return err
	}
	// Defer closure for client exit
	defer conn.Close(websocket.StatusInternalError, "WebSocket closed")
	cancel()

	// Get random message for handshake
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*1)
	var msg Message
	err = wsjson.Read(ctx, conn, &msg)
	if err != nil {
		log.Printf("[CanIDS] failed to establish connection. %s. retrying in %s\n", err, s.RetryDelay)
		return err
	}
	cancel()

	// Decode message
	text, err := base64.StdEncoding.DecodeString(msg.Msg)
	if err != nil {
		log.Printf("[CanIDS] failed to establish connection. %s. retrying in %s\n", err, s.RetryDelay)
		return err
	}

	// Decode encryption key
	key, err := base64.StdEncoding.DecodeString(s.EncryptionKey)
	if err != nil {
		log.Printf("[CanIDS] failed to establish connection. %s. retrying in %s\n", err, s.RetryDelay)
		return err
	}

	// Encrypt received message
	encrypted, err := Encrypt(text, key)
	if err != nil {
		log.Printf("[CanIDS] failed to establish connection. %s. retrying in %s\n", err, s.RetryDelay)
		return err
	}

	// Encode and write message
	encodedString := base64.StdEncoding.EncodeToString(encrypted)

	msg.Msg = encodedString
	msg.MsgType = 0
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*1)
	wsjson.Write(ctx, conn, msg)
	cancel()

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*1)
	err = wsjson.Read(ctx, conn, &msg)
	if err != nil {
		log.Printf("[CanIDS] failed to establish connection. %s. retrying in %s\n", err, s.RetryDelay)
		return err
	}
	cancel()

	if msg.MsgType != 2 {
		err = errNoSuccess
		log.Printf("[CanIDS] failed to establish connection. %s. retrying in %s\n", err, s.RetryDelay)
		return err
	}

	//Success message

	log.Println("Successful connection")

	// Start period poll of file system for new files and stale files
	go fsPollingLoop(s, db)

	go wsReader(s, conn)

	// Start file scanner
	for {

		var frame *UploadRequest
		select {
		case <-queues.pingQueue:
			frame = generatePongFrame(s)
		default:
			// Get next frame, generate JSON payload
			frame, err = scannerGetFrame(s, db, key)
			if err != nil {
				log.Println("[CanIDS] failed to generate frame", err)
				continue
			}
		}

		if s.Encryption {
			frame.Header.Encrypted = true
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
		var msg Message
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

func Encrypt(text []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		log.Println("Failed to generate cipher")
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Println("Failed to generate gcm")
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, text, nil), nil
}
