//go:build mock_google_firebase

package client

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go"
	"github.com/SatorNetwork/sator-api/test/mock"
	"github.com/golang/mock/gomock"
	"google.golang.org/api/option"

	lib_google_firebase "github.com/SatorNetwork/sator-api/lib/google_firebase"
)

type stub struct{}

func (l *stub) Errorf(format string, args ...interface{}) {
	panic(fmt.Sprintf("unexpected mock Errorf stub is called!"+format, args...))
}

func (l *stub) Fatalf(format string, args ...interface{}) {
	panic(fmt.Sprintf("unexpected mock Fatalf stub is called!"+format, args...))
}

func NewApp(ctx context.Context, config *firebase.Config, opts ...option.ClientOption) (lib_google_firebase.AppInterface, error) {
	m := mock.GetMockObject(mock.GoogleFirebaseProvider)
	if m == nil {
		m = lib_google_firebase.NewMockAppInterface(gomock.NewController(&stub{}))
		m.(*lib_google_firebase.MockAppInterface).ExpectMessagingAny()
	}
	return m.(lib_google_firebase.AppInterface), nil
}
