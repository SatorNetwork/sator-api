//go:build !mock_google_firebase

package client

import (
	"context"

	firebase "firebase.google.com/go"
	"github.com/pkg/errors"
	"google.golang.org/api/option"

	lib_google_firebase "github.com/SatorNetwork/sator-api/lib/google_firebase"
)

type (
	FirebaseApp struct {
		client *firebase.App
	}
)

func NewApp(ctx context.Context, config *firebase.Config, opts ...option.ClientOption) (lib_google_firebase.AppInterface, error) {
	app, err := firebase.NewApp(ctx, config, opts...)
	if err != nil {
		return nil, errors.Wrapf(err, "can't initialize firebase app")
	}

	return &FirebaseApp{
		client: app,
	}, nil
}

func (a *FirebaseApp) Messaging(ctx context.Context) (lib_google_firebase.MessagingClientInterface, error) {
	return a.client.Messaging(ctx)
}
