package quiz

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"syreclabs.com/go/faker"
)

type (
	// connection token parse function
	tokenParser func(ctx context.Context, token string) (*TokenPayload, error)
)

// QuizWsHandler handles websocket connections
func QuizWsHandler(tpfn tokenParser) http.HandlerFunc {
	// rooms := make(map[string]broadcast.Broadcaster)

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
		// roomID := tokenPayload.ChallengeRoomID

		client, err := NewWsClient(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		go client.Read()
		go client.Write()

		questions := getQuestions(5)
		answers := make(map[string]QuestionResult)

		ctx, cancel := context.WithCancel(context.Background())
		defer func() {
			cancel()
		}()

		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case answer := <-client.ReadAnswers():
					log.Printf("answer: %+v", answer)
					log.Printf("question: %+v", questions[answer.Payload.QuestionID])
					log.Printf("is correct answer: %v", questions[answer.Payload.QuestionID].correctID == answer.Payload.AnswerID)
					if q, ok := questions[answer.Payload.QuestionID]; ok && q.correctID == answer.Payload.AnswerID {
						answers[answer.Payload.QuestionID] = QuestionResult{
							QuestionID:      q.QuestionID,
							Result:          true,
							Rate:            0,
							CorrectAnswerID: q.correctID,
							QuestionsLeft:   len(questions) - len(answers) + 1,
							AdditionalPts:   0,
						}
					}
				}
			}
		}()

		client.Send(Message{
			Type:   UserConnectedMessage,
			SentAt: time.Now(),
			Payload: User{
				UserID:   uid,
				Username: username,
			},
		})

		for i := 0; i < 9; i++ {
			client.Send(Message{
				Type:   UserConnectedMessage,
				SentAt: time.Now(),
				Payload: User{
					UserID:   uuid.New().String(),
					Username: faker.Internet().UserName(),
				},
			})
			time.Sleep(time.Second)
		}

		for i := 3; i > 0; i-- {
			client.Send(Message{
				Type:   CountdownMessage,
				SentAt: time.Now(),
				Payload: Countdown{
					Countdown: i,
				},
			})
			time.Sleep(time.Second)
		}

		for _, q := range questions {
			client.Send(Message{
				Type:    QuestionMessage,
				SentAt:  time.Now(),
				Payload: q,
			})
			time.Sleep(time.Second * time.Duration(q.TimeForAnswer))

			if res, ok := answers[q.QuestionID]; ok {
				client.Send(Message{
					Type:    QuestionResultMessage,
					SentAt:  time.Now(),
					Payload: res,
				})
				time.Sleep(time.Second * 3)
				continue
			}

			client.Send(Message{
				Type:   QuestionResultMessage,
				SentAt: time.Now(),
				Payload: QuestionResult{
					QuestionID:      q.QuestionID,
					Result:          false,
					CorrectAnswerID: q.correctID,
				},
			})
			time.Sleep(time.Second * 3)
			break
		}

		if len(questions) == len(answers) {
			client.Send(Message{
				Type:   ChallengeResultMessage,
				SentAt: time.Now(),
				Payload: ChallengeResult{
					ChallengeID: challengeID,
					PrizePool:   "250 SAO",
					Winners: []Winner{
						{
							UserID:   uid,
							Username: username,
							Prize:    "125 SAO",
						},
						{
							UserID:   uuid.New().String(),
							Username: faker.Internet().UserName(),
							Prize:    "100 SAO",
						},
						{
							UserID:   uuid.New().String(),
							Username: faker.Internet().UserName(),
							Prize:    "25 SAO",
						},
					},
				},
			})
			time.Sleep(time.Second * 5)
		}
	}
}

func getQuestions(n int) map[string]Question {
	questions := make(map[string]Question)

	for i := 0; i < n; i++ {
		correctID := uuid.New().String()
		qid := uuid.New().String()

		questions[qid] = Question{
			QuestionID:     qid,
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
