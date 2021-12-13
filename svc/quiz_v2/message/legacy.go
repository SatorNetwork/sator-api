package message

import (
	"encoding/json"
	"github.com/pkg/errors"
)

type LegacyMessage struct {
	MessageType MessageType     `json:"type"`
	Payload     json.RawMessage `json:"payload"`
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
		MessageType MessageType `json:"type"`
		Payload     interface{} `json:"payload"`
	}{
		MessageType: m.MessageType,
		Payload:     payload,
	})
}

func (m *Message) UnmarshalJSON(data []byte) error {
	var legacyMsg LegacyMessage
	if err := json.Unmarshal(data, &legacyMsg); err != nil {
		return err
	}
	m.MessageType = legacyMsg.MessageType

	switch legacyMsg.MessageType {
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
	default:
		return errors.Errorf("unknown message type")
	}

	return m.CheckConsistency()
}
