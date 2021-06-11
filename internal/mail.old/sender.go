package mail

import (
	"net/mail"
)

// Address is alias of net/mail.Address.
type Address = mail.Address

// Sender sends emails.
type Sender interface {
	FromAddress() Address
	SendEmail(msg *Message) error
}
