package server

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// upgrader to konwertuje HTTP na WebSocket
var upgrader = websocket.Upgrader{}

// clientMsg to wiadomość od klienta
type clientMsg struct {
	Type string `json:"type"` // "start" | "pause" | "reset"
}

// obsługa WebSocket
func handleWS(s *Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("ws upgrade:", err)
			return
		}
		defer conn.Close()

		done := make(chan struct{})
		go func() {
			defer close(done)
			for {
				var msg clientMsg
				if err := conn.ReadJSON(&msg); err != nil {
					log.Println("ws read:", err)
					return
				}
				switch msg.Type {
				case "start":
					s.Start()
				case "pause":
					s.Pause()
				case "reset":
					s.Reset()
				}
			}
		}()
		ticker := time.NewTicker(time.Second / 60)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				if err := conn.WriteJSON(s.State()); err != nil {
					log.Println("ws write:", err)
					return
				}
			}
		}
	}
}

// uruchamia HTTP na addr i serwuje pliki z webRoot (np. "web").
func ListenAndServe(addr, webRoot string) error {
	mux := http.NewServeMux() // tablica tras: ścieżka → handler
	s := NewSession()
	go s.Run(0.005)
	mux.HandleFunc("/ws", handleWS(s))
	// wszytkie pliki z katalogu web
	mux.Handle("/", http.FileServer(http.Dir(webRoot)))

	return http.ListenAndServe(addr, mux)
}
