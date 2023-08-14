package websocket

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/mcmaster-circ/canids-v2/backend/libraries/elasticsearch"
	"github.com/mcmaster-circ/canids-v2/backend/state"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

const (
	WSPort     = 50000
	bufferSize = 4096
)

type Header struct {
	MsgUuid      string    `json:"msg_uuid,omitempty"`      // Unique message identifier
	MsgTimestamp time.Time `json:"msg_timestamp,omitempty"` // Message timestamp
	ErrorMsg     string    `json:"error_msg,omitempty"`     // Request error message(s) (use with NACK)
	Session      string    `json:"session,omitempty"`       // Connection session UUID
	MsgType      int       `json:"type,omitempty"`          // Message type: 0 - data, 1 - pong
	Encrypted    bool      `json:"encrypted,omitempty"`     // Whether the payload is encrypted (true) or not (false)
}

type Frame struct {
	Header   Header   `json:"header,omitempty"`    // Header
	AssetID  string   `json:"asset_id,omitempty"`  // Asset identifier
	FileName string   `json:"file_name,omitempty"` // Name of file payload is from
	Payload  [][]byte `json:"payload,omitempty"`   // Multiple JSON byte lines from Zeek
	Key      []byte   // For storing associated key
}

type Message struct {
	MsgType int    `json:"type,omitempty"` // Message type: 0 - Misc, 1 - Ping
	Msg     string `json:"msg,omitempty"`
}

// IngestServer handles WebSocket connections.
type WebSocketServer struct {
	queue chan *Frame
}

var server = &WebSocketServer{
	queue: make(chan *Frame, bufferSize),
}

var maxIndexSize = 1000000

func SetMaxElasticIndexSize(newSize int) {
	maxIndexSize = newSize
}

func GetMaxElasticIndexSize() int {
	return maxIndexSize
}

// HandleWebSocket handles incoming WebSocket connections.
func HandleWebSocket(s *state.State, w http.ResponseWriter, r *http.Request) {

	log.Println("Recieved connection request")
	uuid := r.Header.Get("Authorization")

	ingestion, err := elasticsearch.QueryIngestionByUUID(s, uuid)
	if err != nil {
		log.Println("Unable to get specified ingestion from elasticsearch: ", err)
		return
	}

	// Accept (temporarily) ws connection
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close(websocket.StatusInternalError, "WebSocket closed")

	// Generate random data for symmetrical encryption key check
	randData := make([]byte, 32)
	_, err = rand.Read(randData)
	if err != nil {
		log.Println("Failed to generate random data: ", err)
		return
	}

	encodedData := base64.StdEncoding.EncodeToString(randData)

	dataMessage := Message{
		MsgType: 0,
		Msg:     encodedData,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	wsjson.Write(ctx, conn, dataMessage)
	cancel()

	// Get encrypted message
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*1)
	var msg Message
	err = wsjson.Read(ctx, conn, &msg)
	if err != nil {
		log.Println("Did not receive response: ", err)
		cancel()
		return
	}
	cancel()

	decoded, err := base64.StdEncoding.DecodeString(msg.Msg)
	if err != nil {
		log.Println("Failed to decode message", err)
	}

	decodedKey, err := base64.StdEncoding.DecodeString(ingestion.Key)
	if err != nil {
		log.Println("Failed to decode key: ", err)
	}

	decrypted, err := Decrypt(decoded, decodedKey)
	if err != nil {
		log.Println("Failed to decrypt received text: ", err)
	}

	if !bytes.Equal(randData, decrypted) {
		log.Println("Data received does not equal original data")
		return
	}

	log.Println("Successful connection with: ", ingestion.UUID)

	var timeLastPong = time.Now()
	var timeLastPing = time.Now()

	for {

		if timeLastPing.Add(time.Second * 5).Before(time.Now()) {
			timeLastPing = time.Now()
			var pingMessage Message
			pingMessage.MsgType = 1
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
			err := wsjson.Write(ctx, conn, pingMessage)
			cancel()
			if err != nil {
				log.Println("Failed to send ping message: ", err)
				conn.Close(websocket.StatusBadGateway, "Failed to send ping message")
				break
			}
		}

		var frame Frame
		err := wsjson.Read(r.Context(), conn, &frame)
		if err != nil {
			log.Println("Error reading WebSocket message:", err)
			conn.Close(websocket.StatusInvalidFramePayloadData, "Could not read websocket message")
			break
		}
		err = Validate(&frame.Header)
		if err != nil {
			log.Println("Invalid header: ", err)
			continue
		}

		if frame.Header.Encrypted {
			frame.Key = decodedKey
		}

		if frame.Header.MsgType == 1 {
			timeLastPong = time.Now()
			continue
		}

		if timeLastPong.Add(time.Second * 15).Before(time.Now()) {
			log.Println("No pong recieved for 15 seconds")
			conn.Close(websocket.StatusBadGateway, "No pong received")
			break
		}
		server.queue <- &frame
	}
}

func HandleQueue(s *state.State) {
	for {
		chunk := <-server.queue

		ingest(chunk, s, maxIndexSize)
	}
}

func Decrypt(text []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(text) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, text := text[:nonceSize], text[nonceSize:]
	return gcm.Open(nil, nonce, text, nil)
}
