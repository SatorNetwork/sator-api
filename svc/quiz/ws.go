package quiz

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dustin/go-broadcast"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"
	"syreclabs.com/go/faker"
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

type (
	// connection token parse function
	tokenParser func(ctx context.Context, token string) (*TokenPayload, error)
)

// QuizWsHandler handles websocket connections
func QuizWsHandler(tpfn tokenParser) http.HandlerFunc {
	// channels := make(map[string]broadcast.Broadcaster)

	return func(w http.ResponseWriter, r *http.Request) {
		// Current connection variables
		challengeID := chi.URLParam(r, "challenge_id")
		token := chi.URLParam(r, "token")
		tokenPayload, err := tpfn(r.Context(), token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		uid := tokenPayload.UserID
		username := tokenPayload.Username
		roomID := tokenPayload.ChallengeRoomID

		// ws connection
		upgrader := websocket.Upgrader{}
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("upgrade error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer c.Close()

		broadcaster := broadcast.NewBroadcaster(100)
		defer broadcaster.Close()

		// run listeners
		g, ctx := errgroup.WithContext(r.Context())
		g.Go(readMsg(ctx, broadcaster, c))
		g.Go(writeMsg(ctx, broadcaster, c))
		if err := g.Wait(); err != nil {
			log.Printf("unexpected stop service: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func readMsg(ctx context.Context, b broadcast.Broadcaster, conn *websocket.Conn) func() error {
	return func() error {
		// write msg channel
		ch := make(chan interface{}, 100)
		b.Register(ch)

		defer func() {
			conn.Close()
			b.Unregister(ch)
			close(ch)
		}()

		conn.SetReadLimit(maxMessageSize)
		conn.SetReadDeadline(time.Now().Add(pongWait))
		conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

		for {
			select {
			case <-ctx.Done():
				fmt.Printf("readMsg: connection closed")
				return nil
			default:
				mt, messageSrc, err := conn.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						log.Printf("error: %v", err)
					}
					break
				}
				if mt != websocket.TextMessage {
					return fmt.Errorf("unexpected message type: %v", mt)
				}

				message := &Answer{}
				if err := json.Unmarshal(messageSrc, message); err != nil {
					return fmt.Errorf("could not read message: %v", mt)
				}

				ch <- message
			}
		}

		return nil
	}
}

func writeMsg(ctx context.Context, b broadcast.Broadcaster, conn *websocket.Conn) func() error {
	return func() error {
		// Ping ticker
		ticker := time.NewTicker(pingPeriod)
		// write msg channel
		ch := make(chan interface{}, 100)
		b.Register(ch)

		defer func() {
			ticker.Stop()
			conn.Close()
			b.Unregister(ch)
			close(ch)
		}()

		select {
		case <-ctx.Done():
			fmt.Println("writeMsg: finished")
			return nil
		case msg := <-ch:
			fmt.Println("msg:", msg)
			return nil
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return fmt.Errorf("ping: %w", err)
			}
		}

		return nil
	}
}

func getQuestions(n int) map[string]Question {
	questions := make(map[string]Question)

	for i := 0; i < n; i++ {
		correctID := uuid.New().String()
		qid := uuid.New().String()

		questions[qid] = Question{
			QuestionID:     uuid.New().String(),
			QuestionText:   faker.Lorem().Sentence(7),
			TimeForAnswer:  8,
			TotalQuestions: n,
			QuestionNumber: n + 1,
			AnswerOptions: []AnswerOption{
				{
					AnswerID:   uuid.New().String(),
					AnswerText: faker.Lorem().Word(),
				},
				{
					AnswerID:   correctID,
					AnswerText: faker.Lorem().Word(),
				},
				{
					AnswerID:   uuid.New().String(),
					AnswerText: faker.Lorem().Word(),
				},
				{
					AnswerID:   uuid.New().String(),
					AnswerText: faker.Lorem().Word(),
				},
			},
			correctID: correctID,
		}
	}

	return questions
}
