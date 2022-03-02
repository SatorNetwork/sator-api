package message

import (
	"fmt"
)

type ErrInconsistentMessage struct {
	msg *Message
}

func NewErrInconsistentMessage(msg *Message) *ErrInconsistentMessage {
	return &ErrInconsistentMessage{
		msg: msg,
	}
}

func (err *ErrInconsistentMessage) Error() string {
	return fmt.Sprintf("message isn't consistent, message: %+v\n", err.msg)
}
