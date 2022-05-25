package google_firebase

import (
	"firebase.google.com/go/messaging"
	"github.com/golang/mock/gomock"
)

func (m *MockAppInterface) ExpectMessagingAny() *gomock.Call {
	return m.EXPECT().
		Messaging(gomock.Any()).
		Return(&messaging.Client{}, nil).
		AnyTimes()
}
