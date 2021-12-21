package message_container

import (
	"github.com/mohae/deepcopy"

	"github.com/SatorNetwork/sator-api/svc/quiz_v2/message"
)

type (
	predicateFunc func(idx int, msg *message.Message) bool
	modifierFunc  func(msg *message.Message)
)

type messageContainer struct {
	messages []*message.Message
}

func New(messages []*message.Message) *messageContainer {
	messagesCopy := deepcopy.Copy(messages).([]*message.Message)

	return &messageContainer{
		messages: messagesCopy,
	}
}

func (mc *messageContainer) Copy() *messageContainer {
	return New(mc.Messages())
}

func (mc *messageContainer) Messages() []*message.Message {
	return mc.messages
}

func (mc *messageContainer) Modify(p predicateFunc, m modifierFunc) *messageContainer {
	for idx, msg := range mc.messages {
		if p(idx, msg) {
			m(msg)
		}
	}

	return mc
}

func (mc *messageContainer) FilterOut(p predicateFunc) *messageContainer {
	filtered := make([]*message.Message, 0)
	for idx, msg := range mc.messages {
		if p(idx, msg) {
			continue
		}

		filtered = append(filtered, msg)
	}
	mc.messages = filtered

	return mc
}

func PFuncMessageType(messageType message.MessageType) predicateFunc {
	return func(idx int, msg *message.Message) bool {
		return msg.MessageType == messageType
	}
}

func PFuncIndex(pidx int) predicateFunc {
	return func(idx int, msg *message.Message) bool {
		return pidx == idx
	}
}
