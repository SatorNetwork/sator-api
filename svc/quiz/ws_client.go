package quiz

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type (
	WsClient struct {
		conn    *websocket.Conn
		answers chan MessageAnswer
		send    chan Message
	}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func NewWsClient(w http.ResponseWriter, r *http.Request) (*WsClient, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, fmt.Errorf("could not establish websocket connection: %w", err)
	}
	return &WsClient{
		conn:    conn,
		answers: make(chan MessageAnswer, 100),
		send:    make(chan Message, 100),
	}, nil
}

func (c *WsClient) ReadAnswers() <-chan MessageAnswer {
	return c.answers
}

func (c *WsClient) Send(m Message) {
	c.send <- m
}

func (c *WsClient) Read() error {
	defer func() {
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		answer := MessageAnswer{}
		if err := json.Unmarshal(message, &answer); err != nil {
			log.Printf("could not decode response message: %v\nmessage: %+v", err, string(message))
			continue
		}
		c.answers <- answer
	}
	return nil
}

func (c *WsClient) Write() error {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return nil
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return fmt.Errorf("next writer error: %w", err)
			}
			b, err := json.Marshal(message)
			if err != nil {
				return fmt.Errorf("could not encode message: %w\nmessage: %+v", err, message)
			}
			w.Write(b)

			if err := w.Close(); err != nil {
				return fmt.Errorf("close connection: %w", err)
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return fmt.Errorf("ping error: %w", err)
			}
		}
	}
}
