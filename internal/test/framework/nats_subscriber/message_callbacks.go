package nats_subscriber

import (
	"fmt"

	"github.com/stretchr/testify/require"

	"github.com/SatorNetwork/sator-api/svc/quiz_v2/message"
)

var defaultAnswerMap = map[string]string{
	"Joey played Dr. Drake Ramoray on which soap opera show?": "Days of Our Lives",
	"What store does Phoebe hate?":                            "Pottery Barn",
	"Rachel got a job with which company in Paris?":           "Louis Vuitton",
	"Phoebe’s scientist boyfriend David worked in what city?": "Minsk",
	"Monica dated an ophthalmologist named?":                  "Richard",
}

func ReplyWithCorrectAnswerCallback(s *natsSubscriber, msg *message.Message) {
	payload := message.AnswerMessage{
		UserID:     s.userID,
		QuestionID: msg.QuestionMessage.QuestionID,
		AnswerID:   findCorrectAnswerID(defaultAnswerMap, msg),
	}
	respMsg, err := message.NewAnswerMessage(&payload)
	require.NoError(s.t, err)
	err = s.SendMessage(respMsg)
	require.NoError(s.t, err)

	if s.IsDebugModeEnabled() {
		fmt.Printf("Send a message: %v\n", respMsg)
	}
}

func ReplyWithWrongAnswerCallback(s *natsSubscriber, msg *message.Message) {
	payload := message.AnswerMessage{
		UserID:     s.userID,
		QuestionID: msg.QuestionMessage.QuestionID,
		AnswerID:   findWrongAnswerID(defaultAnswerMap, msg),
	}
	respMsg, err := message.NewAnswerMessage(&payload)
	require.NoError(s.t, err)
	err = s.SendMessage(respMsg)
	require.NoError(s.t, err)
}

func NoAnswerCallback(s *natsSubscriber, msg *message.Message) {}

func findCorrectAnswerID(answerMap map[string]string, msg *message.Message) string {
	correctAnswerText := answerMap[msg.QuestionMessage.QuestionText]

	for _, option := range msg.QuestionMessage.AnswerOptions {
		if correctAnswerText == option.AnswerText {
			return option.AnswerID
		}
	}

	return ""
}

func findWrongAnswerID(answerMap map[string]string, msg *message.Message) string {
	correctAnswerText := answerMap[msg.QuestionMessage.QuestionText]

	for _, option := range msg.QuestionMessage.AnswerOptions {
		if correctAnswerText != option.AnswerText {
			return option.AnswerID
		}
	}

	return ""
}
