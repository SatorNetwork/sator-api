package message_verifier

import (
	"github.com/pkg/errors"

	"github.com/SatorNetwork/sator-api/svc/quiz_v2/message"
)

type MessageVerifier struct {
	expectedMessages []*message.Message
	recvMessageChan  <-chan *message.Message
	receivedMessages []*message.Message

	done chan struct{}
}

func New(expectedMessages []*message.Message, recvMessageChan <-chan *message.Message) *MessageVerifier {
	return &MessageVerifier{
		expectedMessages: expectedMessages,
		recvMessageChan:  recvMessageChan,
		receivedMessages: make([]*message.Message, 0),

		done: make(chan struct{}),
	}
}

func (v *MessageVerifier) SetExpectedMessages(expectedMessages []*message.Message) {
	v.expectedMessages = expectedMessages
}

func (v *MessageVerifier) Start() {
LOOP:
	for {
		select {
		case msg := <-v.recvMessageChan:
			v.receivedMessages = append(v.receivedMessages, msg)

		case <-v.done:
			break LOOP
		}
	}
}

func (v *MessageVerifier) Close() {
	v.done <- struct{}{}
}

func (v *MessageVerifier) Verify() error {
	if len(v.expectedMessages) != len(v.receivedMessages) {
		return errors.Errorf("expected %v messages, got: %v", len(v.expectedMessages), len(v.receivedMessages))
	}

	messagesNum := len(v.expectedMessages)
	for i := 0; i < messagesNum; i++ {
		expectedMsg := v.expectedMessages[i]
		receivedMsg := v.receivedMessages[i]

		if err := v.compareMessages(expectedMsg, receivedMsg); err != nil {
			return err
		}
	}

	return nil
}

func (v *MessageVerifier) compareMessages(expectedMsg, receivedMsg *message.Message) error {
	if expectedMsg.MessageType != receivedMsg.MessageType {
		return errors.Errorf(
			"expected %v message type, got: %v",
			expectedMsg.MessageType.String(),
			receivedMsg.MessageType.String(),
		)
	}

	messageType := expectedMsg.MessageType
	switch messageType {
	case message.PlayerIsJoinedMessageType:
		err := v.comparePlayerIsJoinedMessages(expectedMsg, receivedMsg)
		if err != nil {
			return errors.Wrap(err, "error during comparing player is joined messages")
		}
		return nil
	case message.CountdownMessageType:
		err := v.compareCountdownMessages(expectedMsg, receivedMsg)
		if err != nil {
			return errors.Wrap(err, "error during comparing countdown messages")
		}
		return nil
	case message.QuestionMessageType:
		err := v.compareQuestionMessages(expectedMsg, receivedMsg)
		if err != nil {
			return errors.Wrap(err, "error during comparing question messages")
		}
		return nil
	case message.AnswerMessageType:
		return nil
	case message.AnswerReplyMessageType:
		err := v.compareAnswerReplyMessages(expectedMsg, receivedMsg)
		if err != nil {
			return errors.Wrap(err, "error during comparing answer reply messages")
		}
		return nil
	default:
		return errors.Errorf("<unknown message type>")
	}
}

func (v *MessageVerifier) comparePlayerIsJoinedMessages(expectedMsg, receivedMsg *message.Message) error {
	if expectedMsg.PlayerIsJoinedMessage == nil {
		return errors.Errorf("player_is_joined_message shouldn't be nil")
	}
	if receivedMsg.PlayerIsJoinedMessage == nil {
		return errors.Errorf("player_is_joined_message shouldn't be nil")
	}

	if expectedMsg.PlayerIsJoinedMessage.Username != receivedMsg.PlayerIsJoinedMessage.Username {
		return errors.Errorf(
			"expected %v username, got: %v",
			expectedMsg.PlayerIsJoinedMessage.Username,
			receivedMsg.PlayerIsJoinedMessage.Username,
		)
	}

	return nil
}

func (v *MessageVerifier) compareCountdownMessages(expectedMsg, receivedMsg *message.Message) error {
	if expectedMsg.CountdownMessage == nil {
		return errors.Errorf("countdown_message shouldn't be nil")
	}
	if receivedMsg.CountdownMessage == nil {
		return errors.Errorf("countdown_message shouldn't be nil")
	}

	if expectedMsg.CountdownMessage.SecondsLeft != receivedMsg.CountdownMessage.SecondsLeft {
		return errors.Errorf(
			"expected %v seconds left, got: %v",
			expectedMsg.CountdownMessage.SecondsLeft,
			receivedMsg.CountdownMessage.SecondsLeft,
		)
	}

	return nil
}

func (v *MessageVerifier) compareQuestionMessages(expectedMsg, receivedMsg *message.Message) error {
	if expectedMsg.QuestionMessage == nil {
		return errors.Errorf("question_message shouldn't be nil")
	}
	if receivedMsg.QuestionMessage == nil {
		return errors.Errorf("question_message shouldn't be nil")
	}

	if expectedMsg.QuestionMessage.Text != receivedMsg.QuestionMessage.Text {
		return errors.Errorf(
			"expected %v question text, got: %v",
			expectedMsg.QuestionMessage.Text,
			receivedMsg.QuestionMessage.Text,
		)
	}

	return nil
}

func (v *MessageVerifier) compareAnswerReplyMessages(expectedMsg, receivedMsg *message.Message) error {
	if expectedMsg.AnswerReplyMessage == nil {
		return errors.Errorf("answer_reply_message shouldn't be nil")
	}
	if receivedMsg.AnswerReplyMessage == nil {
		return errors.Errorf("answer_reply_message shouldn't be nil")
	}

	if expectedMsg.AnswerReplyMessage.Success != receivedMsg.AnswerReplyMessage.Success {
		return errors.Errorf(
			"expected %v answer reply, got: %v",
			expectedMsg.AnswerReplyMessage.Success,
			receivedMsg.AnswerReplyMessage.Success,
		)
	}

	return nil
}
