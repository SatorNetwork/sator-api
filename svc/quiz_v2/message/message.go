package message

import (
	"encoding/json"
)

type MessageType uint8

const (
	PlayerIsJoinedMessageType = iota
	CountdownMessageType
	QuestionMessageType
	AnswerMessageType
	AnswerReplyMessageType
	WinnersTableMessageType
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
	case WinnersTableMessageType:
		return "winners_table_message_type"
	default:
		return "<unknown message type>"
	}
}

type Message struct {
	MessageType           MessageType            `json:"message_type"`
	PlayerIsJoinedMessage *PlayerIsJoinedMessage `json:"player_is_joined_message,omitempty"`
	CountdownMessage      *CountdownMessage      `json:"countdown_message,omitempty"`
	QuestionMessage       *QuestionMessage       `json:"question_message,omitempty"`
	AnswerMessage         *AnswerMessage         `json:"answer_message,omitempty"`
	AnswerReplyMessage    *AnswerReplyMessage    `json:"answer_reply_message,omitempty"`
	WinnersTableMessage   *WinnersTableMessage   `json:"winners_table_message,omitempty"`
}

func (m *Message) String() string {
	data, _ := json.Marshal(m)
	return string(data)
}

type PlayerIsJoinedMessage struct {
	PlayerID string `json:"player_id"`
	Username string `json:"username"`
}

func NewPlayerIsJoinedMessage(payload *PlayerIsJoinedMessage) *Message {
	return &Message{
		MessageType:           PlayerIsJoinedMessageType,
		PlayerIsJoinedMessage: payload,
	}
}

type CountdownMessage struct {
	SecondsLeft int
}

func NewCountdownMessage(payload *CountdownMessage) *Message {
	return &Message{
		MessageType:      CountdownMessageType,
		CountdownMessage: payload,
	}
}

type QuestionMessage struct {
	QuestionID     string         `json:"question_id"`
	QuestionText   string         `json:"question_text"`
	TimeForAnswer  int            `json:"time_for_answer"`
	QuestionNumber int            `json:"question_number"`
	AnswerOptions  []AnswerOption `json:"answer_options"`
}

type AnswerOption struct {
	AnswerID   string `json:"answer_id"`
	AnswerText string `json:"answer_text"`
}

func NewQuestionMessage(payload *QuestionMessage) *Message {
	return &Message{
		MessageType:     QuestionMessageType,
		QuestionMessage: payload,
	}
}

type AnswerMessage struct {
	UserID string `json:"user_id"`
	// TODO(evg): check that QuestionID is correct
	QuestionID string `json:"question_id"`
	AnswerID   string `json:"answer_id"`
}

func NewAnswerMessage(payload *AnswerMessage) *Message {
	return &Message{
		MessageType:   AnswerMessageType,
		AnswerMessage: payload,
	}
}

type AnswerReplyMessage struct {
	Success         bool `json:"success"`
	SegmentNum      int  `json:"segment_num"`
	IsFastestAnswer bool `json:"is_fastest_answer"`
}

func NewAnswerReplyMessage(payload *AnswerReplyMessage) *Message {
	return &Message{
		MessageType:        AnswerReplyMessageType,
		AnswerReplyMessage: payload,
	}
}

type WinnersTableMessage struct {
	PrizePoolDistribution map[string]float64 `json:"prize_pool_distribution"`
}

func NewWinnersTableMessage(payload *WinnersTableMessage) *Message {
	return &Message{
		MessageType:         WinnersTableMessageType,
		WinnersTableMessage: payload,
	}
}
