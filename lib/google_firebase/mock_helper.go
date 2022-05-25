package google_firebase

import (
	"firebase.google.com/go/messaging"
	"github.com/golang/mock/gomock"
)

func (m *MockAppInterface) ExpectMessagingAny(messagingClient MessagingClientInterface) *gomock.Call {
	return m.EXPECT().
		Messaging(gomock.Any()).
		Return(messagingClient, nil).
		AnyTimes()
}

func (m *MockMessagingClientInterface) ExpectSubscribeToTopicAny() *gomock.Call {
	return m.EXPECT().
		SubscribeToTopic(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(&messaging.TopicManagementResponse{}, nil).
		AnyTimes()
}

//Send(ctx context.Context, message *messaging.Message) (string, error)
func (m *MockMessagingClientInterface) ExpectSendAny() *gomock.Call {
	return m.EXPECT().
		Send(gomock.Any(), gomock.Any()).
		Return("", nil).
		AnyTimes()
}
