package websocket

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

const (
	WSPort     = 50000
	bufferSize = 4096
)

// Frame is a file upload with payloads (lines).
type Frame struct {
	AssetID  string   // AssetID is asset identifier
	FileName string   // FileName is upload file name
	Payload  [][]byte // Payload is a list of JSON byte payloads (lines)
}

// IngestServer handles WebSocket connections.
type WebSocketServer struct {
	queue chan *Frame
}

// HandleWebSocket handles incoming WebSocket connections.
func (s *WebSocketServer) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close(websocket.StatusInternalError, "WebSocket closed")

	for {
		var frame Frame
		err := wsjson.Read(r.Context(), conn, &frame)
		if err != nil {
			log.Println("Error reading WebSocket message:", err)
			break
		}
		s.queue <- &frame
	}
}

// Provision registers and starts the WebSocket service. Returns error if the
// WebSocket server fails to start.
func Provision(ctx context.Context) error {
	server := &WebSocketServer{
		queue: make(chan *Frame, bufferSize),
	}

	http.HandleFunc("/websocket", server.HandleWebSocket)

	serverAddr := fmt.Sprintf(":%d", WSPort)
	serverHandler := http.Server{
		Addr:    serverAddr,
		Handler: nil,
	}

	go func() {
		if err := serverHandler.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start WebSocket server: %v", err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := serverHandler.Shutdown(shutdownCtx); err != nil {
		log.Printf("Failed to gracefully shut down WebSocket server: %v", err)
	}

	return nil
}
