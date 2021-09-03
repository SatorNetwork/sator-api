package quiz

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/oklog/run"
	"syreclabs.com/go/faker"
)

type (
	quizService interface {
		ParseQuizToken(ctx context.Context, token string) (*TokenPayload, error)
		Play(ctx context.Context, quizID, uid uuid.UUID, username string) error
		SetupNewQuizHub(ctx context.Context, qid uuid.UUID) (*Hub, error)
		StoreAnswer(ctx context.Context, userID, quizId, questionID, answerID uuid.UUID) error
	}

	challengesClient interface {
		StoreChallengeAttempt(ctx context.Context, challengeID, userID uuid.UUID) error
	}
)

// QuizWsHandler handles websocket connections
func QuizWsHandler(s quizService, callback func(uid, qid uuid.UUID), c challengesClient, botsTimeout time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := chi.URLParam(r, "token")
		tokenPayload, err := s.ParseQuizToken(r.Context(), token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		uid := uuid.MustParse(tokenPayload.UserID)
		username := tokenPayload.Username
		quizID := uuid.MustParse(tokenPayload.QuizID)

		ctx, cancel := context.WithCancel(context.Background())
		defer func() {
			cancel()
		}()

		quizHub, err := s.SetupNewQuizHub(ctx, quizID)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		client, err := NewWsClient(w, r)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := quizHub.AddPlayer(uid, username); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer func() {
			callback(uid, quizID)
			// log.Println("Callback called ////////////////////// ////////////////////// ////////////////////// //////////////////////")                   // TODO: Remove it!
			// log.Printf("Defer: %v, QuizID: %v ////////////////////// ////////////////////// ////////////////////// //////////////////////", uid, quizID) // TODO: Remove it!
		}()

		if err := c.StoreChallengeAttempt(ctx, quizHub.ChallengeID, uid); err != nil {
			log.Printf("could not store challenge attempt: user_id=%s, challenge_id=%s, error: %v",
				uid.String(), quizHub.ChallengeID.String(), err)
		}

		defer func() {
			quizHub.RemovePlayer(uid)
		}()

		fakePlayers := make([]struct {
			ID       uuid.UUID
			Username string
		}, 0, 10)
		for i := 0; i < 10; i++ {
			fakePlayers = append(fakePlayers, struct {
				ID       uuid.UUID
				Username string
			}{
				ID:       uuid.New(),
				Username: faker.Internet().UserName(),
			})
		}
// 		callback(uid, quizID)
// 		log.Println("Callback called ////////////////////// ////////////////////// ////////////////////// //////////////////////")                 // TODO: Remove it!
// 		log.Printf("uid: %v, QuizID: %v ////////////////////// ////////////////////// ////////////////////// //////////////////////", uid, quizID) // TODO: Remove it!

		var g run.Group
		{
			// read messages from websocket connection
			g.Add(client.Read, func(err error) {
				if err != nil {
					log.Printf("messages reading has been stopped with error: %v", err)
				}
			})

			// write messages to websocket connection
			g.Add(client.Write, func(err error) {
				if err != nil {
					log.Printf("messages reading has been stopped with error: %v", err)
				}
			})

			// connect to quiz hub and listen events
			g.Add(func() error {
				return s.Play(ctx, quizID, uid, username)
			}, func(err error) {
				if err != nil {
					log.Printf("messages reading has been stopped with error: %v", err)
				}
			})

			send := make(chan interface{}, 100)
			quit := make(chan interface{}, 2)
			quizHub.ListenPlayerQuitEvent(uid, quit)
			quizHub.ListenMessageToSend(uid, send)
			defer func() {
				quizHub.UnsubscribeMessageToSend(uid, send)
				quizHub.UnsubscribePlayerQuitEvent(uid, quit)
				close(send)
				close(quit)
			}()

			g.Add(func() error {
				for {
					select {
					case <-ctx.Done():
						return nil
					case answer := <-client.ReadAnswers():
						log.Printf("answer: %+v", answer)

						if err := s.StoreAnswer(
							ctx, uid, quizID,
							uuid.MustParse(answer.Payload.QuestionID),
							uuid.MustParse(answer.Payload.AnswerID),
						); err != nil {
							return fmt.Errorf("could not store answer: %v", err)
						}
					case msg := <-send:
						m := Message{}
						if err := json.Unmarshal(msg.([]byte), &m); err != nil {
							return fmt.Errorf("could not decode message: %v", err)
						}
						log.Printf("send message: %+v", m)
						client.Send(m)
					case e := <-quit:
						log.Printf("quit: %+v", e)
						return client.conn.Close()
					}
				}
			}, func(err error) {
				if err != nil {
					log.Printf("messages reading has been stopped with error: %v", err)
				}
			})
		}

		go func() {
			if botsTimeout > 0 {
				time.Sleep(botsTimeout)
			} else {
				time.Sleep(time.Second * 5)
			}

			for _, u := range fakePlayers {
				time.Sleep(time.Second)
				if err := quizHub.AddPlayer(u.ID, u.Username); err != nil {
					log.Printf("add dummy player %s: %v", u.Username, err)
					continue
				}

				if err := quizHub.Connect(u.ID); err != nil {
					log.Printf("connect dummy player %s: %v", u.Username, err)
				}
			}
		}()

		if err := quizHub.sendConnectedPlayersList(uid); err != nil {
			log.Printf("previous players %s: %v", username, err)
		}

		if err := quizHub.Connect(uid); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		g.Run()
	}
}
