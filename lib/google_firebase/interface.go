package google_firebase

import (
	"context"

	"firebase.google.com/go/messaging"
)

//go:generate mockgen -destination=mock_client.go -package=google_firebase github.com/SatorNetwork/sator-api/lib/google_firebase AppInterface
type AppInterface interface {
	Messaging(ctx context.Context) (MessagingClientInterface, error)
}

//go:generate mockgen -destination=mock_messaging_client.go -package=google_firebase github.com/SatorNetwork/sator-api/lib/google_firebase MessagingClientInterface
type MessagingClientInterface interface {
	SubscribeToTopic(ctx context.Context, tokens []string, topic string) (*messaging.TopicManagementResponse, error)
	Send(ctx context.Context, message *messaging.Message) (string, error)
}
