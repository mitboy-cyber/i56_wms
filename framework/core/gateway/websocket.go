package gateway

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WSClient struct {
	conn     *websocket.Conn
	send     chan []byte
	channels []string
	closed   bool
}

type WSHub struct {
	clients map[string]map[*WSClient]bool
}

func NewWSHub() *WSHub {
	return &WSHub{clients: make(map[string]map[*WSClient]bool)}
}

func (h *WSHub) Run() { log.Println("[WS] Hub started") }

func (h *WSHub) Broadcast(channel string, msg interface{}) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	for c := range h.clients[channel] {
		select {
		case c.send <- data:
		default:
		}
	}
}

func (h *WSHub) Subscribe(channel string, client *WSClient) {
	if h.clients[channel] == nil {
		h.clients[channel] = make(map[*WSClient]bool)
	}
	h.clients[channel][client] = true
	client.channels = append(client.channels, channel)
}

func (h *WSHub) Unsubscribe(client *WSClient) {
	client.closed = true
	close(client.send)
	for _, ch := range client.channels {
		delete(h.clients[ch], client)
	}
}

func (gw *Gateway) wsHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	client := &WSClient{conn: conn, send: make(chan []byte, 64)}
	gw.wsHub.Subscribe("parcels", client)
	gw.wsHub.Subscribe("orders", client)
	go func() {
		defer func() { gw.wsHub.Unsubscribe(client); conn.Close() }()
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				break
			}
			var cmd map[string]interface{}
			if json.Unmarshal(msg, &cmd) == nil {
				if ch, ok := cmd["subscribe"].(string); ok {
					gw.wsHub.Subscribe(ch, client)
				}
			}
		}
	}()
	go func() {
		for msg := range client.send {
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				break
			}
		}
	}()
}
