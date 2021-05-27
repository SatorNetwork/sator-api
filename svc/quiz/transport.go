package quiz

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/SatorNetwork/sator-api/internal/httpencoder"
	"github.com/SatorNetwork/sator-api/internal/jwt"
	"github.com/dustin/go-broadcast"
	"github.com/go-chi/chi"
	jwtkit "github.com/go-kit/kit/auth/jwt"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"syreclabs.com/go/faker"
)

// Predefined message types
const (
	MsgTypeUserConnected   = "player_connected"
	MsgTypeUserDisonnected = "player_disconnected"
	MsgTypeCountdown       = "countdown"
	MsgTypeQuestion        = "question"
	MsgTypeQuestionResult  = "question_result"
	MsgTypeQuizResult      = "quiz_result"
	MsgTypeChallengeResult = "challenge_result"
	MsgTypeLose            = "user_lose"
	MsgTypeWin             = "user_win"
)

type (
	logger interface {
		Log(keyvals ...interface{}) error
	}

	// WsMsg ...
	WsMsg struct {
		Type    string      `json:"type"`
		Payload interface{} `json:"payload"`
		SentAt  time.Time   `json:"sent_at"`
	}

	// MsgUserConnected ...
	MsgUserConnected struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	}

	// MsgBeReady ...
	MsgBeReady struct {
		Countdown int `json:"countdown"`
	}

	// MsgQuestion ...
	MsgQuestion struct {
		ID            string    `json:"id"`
		Question      string    `json:"question"`
		Options       []QOption `json:"options"`
		TimeForAnswer int       `json:"time_for_answer"`
		correctID     string    `json:"-"`
	}

	// QOption ...
	QOption struct {
		ID     string `json:"id"`
		Option string `json:"option"`
	}

	// MsgQuestionResult ...
	MsgQuestionResult struct {
		ID              string `json:"id"`
		Result          int    `json:"result"`
		Rate            int    `json:"rate"`
		CorrectAnswerID string `json:"correct_answer_id"`
	}

	// MsgQuizResult ...
	MsgQuizResult struct {
		Result  int    `json:"result"`
		Message string `json:"message"`
	}

	// MsgChallengeResult ...
	MsgChallengeResult struct {
		ChallengeID string   `json:"challenge_id"`
		PrizePool   float64  `json:"prize_pool"`
		Winners     []Winner `json:"winners"`
	}

	// Winner ...
	Winner struct {
		ID       string  `json:"id"`
		Username string  `json:"username"`
		Prize    float64 `json:"prize"`
	}

	// ReqAnswer ...
	ReqAnswer struct {
		QuestionID string `json:"question_id"`
		AnswerID   string `json:"answer_id"`
	}
)

// MakeHTTPHandler ...
func MakeHTTPHandler(_ Endpoints, log logger) http.Handler {
	r := chi.NewRouter()

	r.Get("/{challenge_id}", quizHandler(
		jwtkit.HTTPToContext(),
		httpencoder.EncodeError(log, codeAndMessageFrom),
	))

	return r
}

// returns http error code by error type
func codeAndMessageFrom(err error) (int, interface{}) {
	return httpencoder.CodeAndMessageFrom(err)
}

func quizHandler(jwtToContextFn httptransport.RequestFunc, errEnc httptransport.ErrorEncoder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := jwtToContextFn(r.Context(), r)
		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			errEnc(ctx, err, w)
			log.Println("error:", err)
			return
		}
		username, err := jwt.UsernameFromContext(ctx)
		if err != nil {
			errEnc(ctx, err, w)
			log.Println("error:", err)
			return
		}

		channelID := uid.String()

		listener := openListener(channelID)
		defer closeListener(channelID, listener)

		finished := make(chan bool, 5)
		defer close(finished)

		upgrader := websocket.Upgrader{}
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("upgrade error:", err)
			return
		}
		defer c.Close()

		questions := getQuestions(10)

		wg := &sync.WaitGroup{}
		wg.Add(2)

		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					log.Printf("[client_disconnected] context done: %s", channelID)
					return
				case <-finished:
					log.Printf("quiz finished: %s", channelID)
					return
				default:
					mt, message, err := c.ReadMessage()
					if err != nil {
						log.Println("read:", err)
						break
					}
					log.Printf("recv: [%d] %s", mt, message)

					if mt != websocket.TextMessage {
						break
					}

					msg := &ReqAnswer{}
					if err := json.Unmarshal(message, msg); err != nil {
						log.Printf("could not decode received message: %v", err)
						break
					}

					q, ok := questions[msg.QuestionID]
					if !ok {
						log.Printf("could not find question with id=%s", msg.QuestionID)
						channel(channelID).Submit(MsgQuizResult{
							Result:  0,
							Message: "could not find the question you replied",
						})
						finished <- true
						break
					}

					if q.correctID != msg.AnswerID {
						channel(channelID).Submit(MsgQuestionResult{
							ID:              q.ID,
							Result:          0,
							CorrectAnswerID: q.correctID,
						})
						channel(channelID).Submit(MsgQuizResult{
							Result:  0,
							Message: "you lose",
						})
						finished <- true
						break
					}

					channel(channelID).Submit(MsgQuestionResult{
						ID:              q.ID,
						Result:          1,
						CorrectAnswerID: q.correctID,
					})
				}
			}
		}()

		go func() {
			defer wg.Done()
			for {
				select {
				case <-r.Context().Done():
					log.Printf("[client_disconnected] client closed connection: %s", channelID)
					return
				case <-ctx.Done():
					log.Printf("[client_disconnected] context done: %s", channelID)
					return
				case <-finished:
					log.Printf("quiz finished: %s", channelID)
					return
				case event := <-listener:
					var msg []byte
					switch e := event.(type) {
					case WsMsg:
						msg, _ = json.Marshal(e)
					case MsgUserConnected:
						msg, _ = json.Marshal(WsMsg{
							Type:    MsgTypeUserConnected,
							Payload: e,
							SentAt:  time.Now(),
						})
					case MsgBeReady:
						msg, _ = json.Marshal(WsMsg{
							Type:    MsgTypeCountdown,
							Payload: e,
							SentAt:  time.Now(),
						})
					case MsgQuestion:
						msg, _ = json.Marshal(WsMsg{
							Type:    MsgTypeQuestion,
							Payload: e,
							SentAt:  time.Now(),
						})
					case MsgQuestionResult:
						msg, _ = json.Marshal(WsMsg{
							Type:    MsgTypeQuestionResult,
							Payload: e,
							SentAt:  time.Now(),
						})
					case MsgQuizResult:
						msg, _ = json.Marshal(WsMsg{
							Type:    MsgTypeQuizResult,
							Payload: e,
							SentAt:  time.Now(),
						})
					case MsgChallengeResult:
						msg, _ = json.Marshal(WsMsg{
							Type:    MsgTypeChallengeResult,
							Payload: e,
							SentAt:  time.Now(),
						})
					}

					if err := c.WriteMessage(websocket.TextMessage, msg); err != nil {
						log.Println("write:", err)
					}
				}
			}
		}()

		channel(channelID).Submit(MsgUserConnected{
			ID:       uid.String(),
			Username: username,
		})

		for i := 0; i < 9; i++ {
			channel(channelID).Submit(MsgUserConnected{
				ID:       uuid.New().String(),
				Username: faker.Internet().UserName(),
			})
			time.Sleep(time.Second)
		}

		for i := 3; i > 0; i-- {
			channel(channelID).Submit(MsgBeReady{
				Countdown: i,
			})
			time.Sleep(time.Millisecond * 1500)
		}

		for _, q := range questions {
			channel(channelID).Submit(q)
			time.Sleep(time.Second * 8)
		}

		channel(channelID).Submit(MsgQuizResult{
			Result:  1,
			Message: "Congrats!",
		})

		wg.Wait()
	}
}

func getQuestions(n int) map[string]MsgQuestion {
	questions := make(map[string]MsgQuestion)

	for i := 0; i < n; i++ {
		correctID := uuid.New().String()
		qid := uuid.New().String()

		questions[qid] = MsgQuestion{
			ID:       uuid.New().String(),
			Question: faker.Lorem().Sentence(7),
			Options: []QOption{
				{
					ID:     uuid.New().String(),
					Option: faker.Lorem().Word(),
				},
				{
					ID:     correctID,
					Option: faker.Lorem().Word(),
				},
				{
					ID:     uuid.New().String(),
					Option: faker.Lorem().Word(),
				},
				{
					ID:     uuid.New().String(),
					Option: faker.Lorem().Word(),
				},
			},
			correctID: correctID,
		}
	}

	return questions
}

var channels = make(map[string]broadcast.Broadcaster)

func openListener(channelID string) chan interface{} {
	listener := make(chan interface{})
	channel(channelID).Register(listener)
	return listener
}

func closeListener(channelID string, listener chan interface{}) {
	channel(channelID).Unregister(listener)
	close(listener)
}

func channel(channelID string) broadcast.Broadcaster {
	channelID = strings.ToLower(channelID)
	b, ok := channels[channelID]
	if !ok {
		b = broadcast.NewBroadcaster(10)
		channels[channelID] = b
	}
	return b
}

func deleteBroadcast(channelID string) {
	channelID = strings.ToLower(channelID)
	b, ok := channels[channelID]
	if ok {
		b.Close()
		delete(channels, channelID)
	}
}
