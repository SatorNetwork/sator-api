package message

import (
	"encoding/json"
	"fmt"
)

type MessageType uint8

const (
	PlayerIsJoinedMessageType = iota
	CountdownMessageType
	QuestionMessageType
	AnswerMessageType
	AnswerReplyMessageType
)

func (m MessageType) String() string {
	switch m {
	case PlayerIsJoinedMessageType:
		return "player_is_joined_message_type"
	case CountdownMessageType:
		return "countdown_message_type"
	case QuestionMessageType:
		return "question_message_type"
	case AnswerMessageType:
		return "answer_message_type"
	case AnswerReplyMessageType:
		return "answer_reply_message_type"
	default:
		return "<unknown message type>"
	}
}

type Message struct {
	MessageType MessageType `json:"message_type,omitempty"`
	//Payload     interface{}
	PlayerIsJoinedMessage *PlayerIsJoinedMessage `json:"player_is_joined_message,omitempty"`
	CountdownMessage      *CountdownMessage      `json:"countdown_message,omitempty"`
	QuestionMessage       *QuestionMessage       `json:"question_message,omitempty"`
	AnswerMessage         *AnswerMessage         `json:"answer_message,omitempty"`
	AnswerReplyMessage    *AnswerReplyMessage    `json:"answer_reply_message,omitempty"`
}

func (m *Message) String() string {
	data, _ := json.Marshal(m)
	return string(data)
}

type PlayerIsJoinedMessage struct {
	PlayerID string
	Username string
}

func NewPlayerIsJoinedMessage(payload *PlayerIsJoinedMessage) *Message {
	return &Message{
		MessageType:           PlayerIsJoinedMessageType,
		PlayerIsJoinedMessage: payload,
	}
}

func (m *PlayerIsJoinedMessage) String() string {
	tmpl := `
PlayerID: %v
`
	return fmt.Sprintf(tmpl, m.PlayerID)
}

type CountdownMessage struct{
	SecondsLeft int
}

func NewCountdownMessage(payload *CountdownMessage) *Message {
	return &Message{
		MessageType:      CountdownMessageType,
		CountdownMessage: payload,
	}
}

type QuestionMessage struct {
	Text string
}

func NewQuestionMessage(payload *QuestionMessage) *Message {
	return &Message{
		MessageType:     QuestionMessageType,
		QuestionMessage: payload,
	}
}

type AnswerMessage struct {
	AnswerFlag bool
	UserID     string
}

func NewAnswerMessage(payload *AnswerMessage) *Message {
	return &Message{
		MessageType:   AnswerMessageType,
		AnswerMessage: payload,
	}
}

type AnswerReplyMessage struct {
	Success bool
}

func NewAnswerReplyMessage(payload *AnswerReplyMessage) *Message {
	return &Message{
		MessageType:        AnswerReplyMessageType,
		AnswerReplyMessage: payload,
	}
}
