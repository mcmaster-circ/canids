package websocket

import (
	"context"
	"log"
	"net/http"
	"time"

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
}

type Frame struct {
	Header   Header   `json:"header,omitempty"`    // Header
	AssetID  string   `json:"asset_id,omitempty"`  // Asset identifier
	FileName string   `json:"file_name,omitempty"` // Name of file payload is from
	Payload  [][]byte `json:"payload,omitempty"`   // Multiple JSON byte lines from Zeek
}

type SentMsg struct {
	MsgType int `json:"type,omitempty"` // Message type: 0 - Misc, 1 - Ping
}

// IngestServer handles WebSocket connections.
type WebSocketServer struct {
	queue chan *Frame
}

var server = &WebSocketServer{
	queue: make(chan *Frame, bufferSize),
}

var allowedKeys = []string{"hello", "there"}
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
	token := r.Header.Get("Authorization")
	allowed := false
	for _, item := range allowedKeys {
		if token == item {
			allowed = true
			log.Println("Valid auth token")
		}
	}

	if !allowed {
		log.Println("Invalid auth token")
		return
	}

	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close(websocket.StatusInternalError, "WebSocket closed")

	var timeLastPong = time.Now()
	var timeLastPing = time.Now()

	for {

		if timeLastPing.Add(time.Second * 5).Before(time.Now()) {
			timeLastPing = time.Now()
			var pingMessage SentMsg
			pingMessage.MsgType = 1
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
			err := wsjson.Write(ctx, conn, pingMessage)
			cancel()
			if err != nil {
				log.Println("Failed to send ping message: ", err)
				conn.Close(websocket.StatusBadGateway, "Failed to send ping message")
				break
			}
			log.Println("Successfully sent ping message")
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

		if frame.Header.MsgType == 1 {
			log.Printf("Pong recieved")
			timeLastPong = time.Now()
			continue
		}

		if timeLastPong.Add(time.Second * 15).Before(time.Now()) {
			log.Println("No pong recieved for 15 seconds")
			conn.Close(websocket.StatusBadGateway, "No pong received")
			break
		}

		log.Printf("[ws] Frame recieved\n")
		server.queue <- &frame
	}
}

func HandleQueue(s *state.State) {
	for {
		chunk := <-server.queue

		ingest(chunk, s, maxIndexSize)
	}
}
