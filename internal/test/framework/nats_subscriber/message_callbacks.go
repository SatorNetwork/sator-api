package nats_subscriber

import (
	"github.com/stretchr/testify/require"

	"github.com/SatorNetwork/sator-api/svc/quiz_v2/message"
)

func ReplyWithTrueCallback(s *natsSubscriber, msg *message.Message) {
	payload := message.AnswerMessage{
		AnswerFlag: true,
		UserID:     s.userID,
	}
	respMsg := message.NewAnswerMessage(&payload)
	err := s.SendMessage(respMsg)
	require.NoError(s.t, err)
}
