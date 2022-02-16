package message

import (
	"time"

	"github.com/pkg/errors"
)

const (
	defaultTTLInMilliseconds   = 5000
	countdownTTLInMilliseconds = 3000
)

type MessageType uint8

const (
	PlayerIsJoinedMessageType = iota
	CountdownMessageType
	QuestionMessageType
	AnswerMessageType
	AnswerReplyMessageType
	WinnersTableMessageType
	PlayerIsActiveMessageType
	PlayerIsDisconnectedMessageType
)

func NewMessageTypeFromString(messageType string) (MessageType, error) {
	switch messageType {
	case "player_is_joined":
		return PlayerIsJoinedMessageType, nil
	case "countdown":
		return CountdownMessageType, nil
	case "question":
		return QuestionMessageType, nil
	case "answer":
		return AnswerMessageType, nil
	case "answer_reply":
		return AnswerReplyMessageType, nil
	case "winners_table":
		return WinnersTableMessageType, nil
	case "player_is_active":
		return PlayerIsActiveMessageType, nil
	case "player_is_disconnected":
		return PlayerIsDisconnectedMessageType, nil
	default:
		return 0, errors.Errorf("unknown message type: %v", messageType)
	}
}

func (m MessageType) String() string {
	switch m {
	case PlayerIsJoinedMessageType:
		return "player_is_joined"
	case CountdownMessageType:
		return "countdown"
	case QuestionMessageType:
		return "question"
	case AnswerMessageType:
		return "answer"
	case AnswerReplyMessageType:
		return "answer_reply"
	case WinnersTableMessageType:
		return "winners_table"
	case PlayerIsActiveMessageType:
		return "player_is_active"
	case PlayerIsDisconnectedMessageType:
		return "player_is_disconnected"
	default:
		return "<unknown message type>"
	}
}

type Message struct {
	MessageType                 MessageType                  `json:"message_type"`
	PlayerIsJoinedMessage       *PlayerIsJoinedMessage       `json:"player_is_joined_message,omitempty"`
	CountdownMessage            *CountdownMessage            `json:"countdown_message,omitempty"`
	QuestionMessage             *QuestionMessage             `json:"question_message,omitempty"`
	AnswerMessage               *AnswerMessage               `json:"answer_message,omitempty"`
	AnswerReplyMessage          *AnswerReplyMessage          `json:"answer_reply_message,omitempty"`
	WinnersTableMessage         *WinnersTableMessage         `json:"winners_table_message,omitempty"`
	PlayerIsActiveMessage       *PlayerIsActiveMessage       `json:"player_is_active_message,omitempty"`
	PlayerIsDisconnectedMessage *PlayerIsDisconnectedMessage `json:"player_is_disconnected_message,omitempty"`
	Date                        string                       `json:"date,omitempty"`
	TTL                         int                          `json:"ttl,omitempty"`
}

func (m *Message) GetAnswerMessage() (*AnswerMessage, error) {
	if err := m.CheckConsistency(); err != nil {
		return nil, err
	}

	return m.AnswerMessage, nil
}

// MustGetAnswerMessage may return potentially inconsistent message. It's better to use GetAnswerMessage.
func (m *Message) MustGetAnswerMessage() *AnswerMessage {
	return m.AnswerMessage
}

func (m *Message) CheckConsistency() error {
	if !m.isConsistent() {
		return NewErrInconsistentMessage(m)
	}

	return nil
}

func (m *Message) isConsistent() bool {
	switch m.MessageType {
	case PlayerIsJoinedMessageType:
		return m.PlayerIsJoinedMessage != nil
	case CountdownMessageType:
		return m.CountdownMessage != nil
	case QuestionMessageType:
		return m.QuestionMessage != nil
	case AnswerMessageType:
		ok := m.AnswerMessage != nil &&
			m.AnswerMessage.QuestionID != "" &&
			m.AnswerMessage.AnswerID != ""
		return ok
	case AnswerReplyMessageType:
		return m.AnswerReplyMessage != nil
	case WinnersTableMessageType:
		return m.WinnersTableMessage != nil
	case PlayerIsActiveMessageType:
		return m.PlayerIsActiveMessage != nil
	case PlayerIsDisconnectedMessageType:
		return m.PlayerIsDisconnectedMessage != nil
	default:
		return false
	}
}

type PlayerIsJoinedMessage struct {
	PlayerID string `json:"user_id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

func NewPlayerIsJoinedMessage(payload *PlayerIsJoinedMessage) (*Message, error) {
	msg := &Message{
		MessageType:           PlayerIsJoinedMessageType,
		PlayerIsJoinedMessage: payload,
		Date:                  time.Now().Format(time.RFC3339),
		TTL:                   defaultTTLInMilliseconds,
	}
	if err := msg.CheckConsistency(); err != nil {
		return nil, err
	}

	return msg, nil
}

type CountdownMessage struct {
	SecondsLeft int `json:"countdown"`
}

func NewCountdownMessage(payload *CountdownMessage) (*Message, error) {
	msg := &Message{
		MessageType:      CountdownMessageType,
		CountdownMessage: payload,
		Date:             time.Now().Format(time.RFC3339),
		TTL:              countdownTTLInMilliseconds,
	}
	if err := msg.CheckConsistency(); err != nil {
		return nil, err
	}

	return msg, nil
}

type QuestionMessage struct {
	QuestionID     string         `json:"question_id"`
	QuestionText   string         `json:"question_text"`
	TimeForAnswer  int            `json:"time_for_answer"`
	QuestionNumber int            `json:"question_number"`
	TotalQuestions int            `json:"total_questions"`
	AnswerOptions  []AnswerOption `json:"answer_options"`
}

type AnswerOption struct {
	AnswerID   string `json:"answer_id"`
	AnswerText string `json:"answer_text"`
}

func NewQuestionMessage(payload *QuestionMessage, ttl int) (*Message, error) {
	msg := &Message{
		MessageType:     QuestionMessageType,
		QuestionMessage: payload,
		Date:            time.Now().Format(time.RFC3339),
		TTL:             ttl,
	}
	if err := msg.CheckConsistency(); err != nil {
		return nil, err
	}

	return msg, nil
}

type AnswerMessage struct {
	// TODO(evg): check that QuestionID is correct
	QuestionID string `json:"question_id"`
	AnswerID   string `json:"answer_id"`
}

func NewAnswerMessage(payload *AnswerMessage) (*Message, error) {
	msg := &Message{
		MessageType:   AnswerMessageType,
		AnswerMessage: payload,
		Date:          time.Now().Format(time.RFC3339),
		TTL:           defaultTTLInMilliseconds,
	}
	if err := msg.CheckConsistency(); err != nil {
		return nil, err
	}

	return msg, nil
}

type AnswerReplyMessage struct {
	QuestionID      string `json:"question_id"`
	Success         bool   `json:"result"`
	Rate            int    `json:"rate"`
	CorrectAnswerID string `json:"correct_answer_id"`
	QuestionsLeft   int    `json:"questions_left"`
	AdditionalPTS   int    `json:"additional_pts"`
	SegmentNum      int    `json:"segment_num"`
	IsFastestAnswer bool   `json:"is_fastest_answer"`
}

func NewAnswerReplyMessage(payload *AnswerReplyMessage) (*Message, error) {
	msg := &Message{
		MessageType:        AnswerReplyMessageType,
		AnswerReplyMessage: payload,
		Date:               time.Now().Format(time.RFC3339),
		TTL:                defaultTTLInMilliseconds,
	}
	if err := msg.CheckConsistency(); err != nil {
		return nil, err
	}

	return msg, nil
}

type WinnersTableMessage struct {
	ChallengeID           string             `json:"challenge_id"`
	PrizePool             string             `json:"prize_pool"`
	ShowTransactionURL    string             `json:"show_transaction_url"`
	Winners               []*Winner          `json:"winners"`
	Losers                []*Loser           `json:"losers"`
	PrizePoolDistribution map[string]float64 `json:"prize_pool_distribution"`
}

type Winner struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Prize    string `json:"prize"`
	Bonus    string `json:"bonus"`
	Avatar   string `json:"avatar"`
}

type Loser struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	PTS      uint32 `json:"pts"`
	Avatar   string `json:"avatar"`
}

func NewWinnersTableMessage(payload *WinnersTableMessage) (*Message, error) {
	msg := &Message{
		MessageType:         WinnersTableMessageType,
		WinnersTableMessage: payload,
		Date:                time.Now().Format(time.RFC3339),
		TTL:                 defaultTTLInMilliseconds,
	}
	if err := msg.CheckConsistency(); err != nil {
		return nil, err
	}

	return msg, nil
}

type PlayerIsActiveMessage struct{}

func NewPlayerIsActiveMessage(payload *PlayerIsActiveMessage) (*Message, error) {
	msg := &Message{
		MessageType:           PlayerIsActiveMessageType,
		PlayerIsActiveMessage: payload,
		Date:                  time.Now().Format(time.RFC3339),
		TTL:                   defaultTTLInMilliseconds,
	}
	if err := msg.CheckConsistency(); err != nil {
		return nil, err
	}

	return msg, nil
}

type PlayerIsDisconnectedMessage struct {
	PlayerID string `json:"user_id"`
	Username string `json:"username"`
}

func NewPlayerIsDisconnectedMessage(payload *PlayerIsDisconnectedMessage) (*Message, error) {
	msg := &Message{
		MessageType:                 PlayerIsDisconnectedMessageType,
		PlayerIsDisconnectedMessage: payload,
		Date:                        time.Now().Format(time.RFC3339),
		TTL:                         defaultTTLInMilliseconds,
	}
	if err := msg.CheckConsistency(); err != nil {
		return nil, err
	}

	return msg, nil
}
