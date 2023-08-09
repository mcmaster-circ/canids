package websocket

import (
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
}

type Frame struct {
	Header   Header   `json:"header,omitempty"`    // Header
	AssetID  string   `json:"asset_id,omitempty"`  // Asset identifier
	FileName string   `json:"file_name,omitempty"` // Name of file payload is from
	Payload  [][]byte `json:"payload,omitempty"`   // Multiple JSON byte lines from Zeek
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

	go handleQueue(s)

	for {
		var frame Frame
		err := wsjson.Read(r.Context(), conn, &frame)
		if err != nil {
			log.Println("Error reading WebSocket message:", err)
			break
		}
		err = Validate(&frame.Header)
		if err != nil {
			log.Println("Invalid header: ", err)
		}
		log.Printf("[ws] Frame recieved\n")
		server.queue <- &frame
	}
}

func handleQueue(s *state.State) {
	for {
		chunk := <-server.queue

		ingest(chunk, s, maxIndexSize)
	}
}
