package message

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type LegacyMessage struct {
	MessageType string          `json:"type"`
	Payload     json.RawMessage `json:"payload"`
	Date        string          `json:"date"`
	TTL         int             `json:"ttl"`
}

func (m *Message) getPayload() (interface{}, error) {
	if err := m.CheckConsistency(); err != nil {
		return nil, err
	}

	switch m.MessageType {
	case PlayerIsJoinedMessageType:
		return m.PlayerIsJoinedMessage, nil
	case CountdownMessageType:
		return m.CountdownMessage, nil
	case QuestionMessageType:
		return m.QuestionMessage, nil
	case AnswerMessageType:
		return m.AnswerMessage, nil
	case AnswerReplyMessageType:
		return m.AnswerReplyMessage, nil
	case WinnersTableMessageType:
		return m.WinnersTableMessage, nil
	case PlayerIsActiveMessageType:
		return m.PlayerIsActiveMessage, nil
	case PlayerIsDisconnectedMessageType:
		return m.PlayerIsDisconnectedMessage, nil
	case TimeOutMessageType:
		return m.TimeOutMessage, nil
	default:
		return nil, NewErrInconsistentMessage(m)
	}
}

func (m *Message) MarshalJSON() ([]byte, error) {
	payload, err := m.getPayload()
	if err != nil {
		return nil, err
	}

	return json.Marshal(&struct {
		MessageType string      `json:"type"`
		Payload     interface{} `json:"payload"`
		Date        string      `json:"date"`
		TTL         int         `json:"ttl"`
	}{
		MessageType: m.MessageType.String(),
		Payload:     payload,
		Date:        m.Date,
		TTL:         m.TTL,
	})
}

func (m *Message) UnmarshalJSON(data []byte) error {
	var legacyMsg LegacyMessage
	if err := json.Unmarshal(data, &legacyMsg); err != nil {
		return err
	}
	var err error
	m.MessageType, err = NewMessageTypeFromString(legacyMsg.MessageType)
	if err != nil {
		return errors.Wrap(err, "can't create new message type from string")
	}
	m.Date = legacyMsg.Date
	m.TTL = legacyMsg.TTL

	switch m.MessageType {
	case PlayerIsJoinedMessageType:
		if err := json.Unmarshal(legacyMsg.Payload, &m.PlayerIsJoinedMessage); err != nil {
			return err
		}
	case CountdownMessageType:
		if err := json.Unmarshal(legacyMsg.Payload, &m.CountdownMessage); err != nil {
			return err
		}
	case QuestionMessageType:
		if err := json.Unmarshal(legacyMsg.Payload, &m.QuestionMessage); err != nil {
			return err
		}
	case AnswerMessageType:
		if err := json.Unmarshal(legacyMsg.Payload, &m.AnswerMessage); err != nil {
			return err
		}
	case AnswerReplyMessageType:
		if err := json.Unmarshal(legacyMsg.Payload, &m.AnswerReplyMessage); err != nil {
			return err
		}
	case WinnersTableMessageType:
		if err := json.Unmarshal(legacyMsg.Payload, &m.WinnersTableMessage); err != nil {
			return err
		}
	case PlayerIsActiveMessageType:
		if err := json.Unmarshal(legacyMsg.Payload, &m.PlayerIsActiveMessage); err != nil {
			return err
		}
	case PlayerIsDisconnectedMessageType:
		if err := json.Unmarshal(legacyMsg.Payload, &m.PlayerIsDisconnectedMessage); err != nil {
			return err
		}
	case TimeOutMessageType:
		if err := json.Unmarshal(legacyMsg.Payload, &m.TimeOutMessage); err != nil {
			return err
		}
	default:
		return errors.Errorf("unknown message type")
	}

	return m.CheckConsistency()
}
