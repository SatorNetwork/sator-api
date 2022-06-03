package firebase

import (
	"context"
	"fmt"

	"firebase.google.com/go/messaging"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"google.golang.org/api/option"

	lib_google_firebase "github.com/SatorNetwork/sator-api/lib/google_firebase"
	google_firebase_client "github.com/SatorNetwork/sator-api/lib/google_firebase/client"
	firebase_repository "github.com/SatorNetwork/sator-api/svc/firebase/repository"
)

type (
	Service struct {
		fr              firebaseRepository
		app             lib_google_firebase.AppInterface
		messagingClient lib_google_firebase.MessagingClientInterface
	}

	firebaseRepository interface {
		GetRegistrationToken(
			ctx context.Context,
			arg firebase_repository.GetRegistrationTokenParams,
		) (firebase_repository.FirebaseRegistrationToken, error)
		UpsertRegistrationToken(
			ctx context.Context,
			arg firebase_repository.UpsertRegistrationTokenParams,
		) error

		IsNotificationEnabled(ctx context.Context, arg firebase_repository.IsNotificationEnabledParams) (bool, error)
		IsNotificationDisabled(ctx context.Context, arg firebase_repository.IsNotificationDisabledParams) (bool, error)
	}

	Empty struct{}

	RegisterTokenRequest struct {
		DeviceId string `json:"device_id"`
		Token    string `json:"token"`
	}
)

func NewService(
	fr firebaseRepository,
	credsInJSON []byte,
) (*Service, error) {
	creds := option.WithCredentialsJSON(credsInJSON)
	app, err := google_firebase_client.NewApp(context.Background(), nil, creds)
	if err != nil {
		return nil, errors.Wrapf(err, "can't initialize firebase app")
	}
	messagingClient, err := app.Messaging(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "can't initialize firebase messaging client")
	}

	s := &Service{
		fr:              fr,
		app:             app,
		messagingClient: messagingClient,
	}

	return s, nil
}

func (s *Service) RegisterToken(ctx context.Context, userID uuid.UUID, req *RegisterTokenRequest) (*Empty, error) {
	err := s.fr.UpsertRegistrationToken(ctx, firebase_repository.UpsertRegistrationTokenParams{
		UserID:            userID,
		DeviceID:          req.DeviceId,
		RegistrationToken: req.Token,
	})
	if err != nil {
		return nil, errors.Wrap(err, "can't register token")
	}

	topics, err := s.getUserSubscribedTopics(ctx, userID)
	if err != nil {
		return nil, err
	}
	if err := s.subscribeTokenToTopics(ctx, req.Token, topics); err != nil {
		return nil, err
	}

	return &Empty{}, nil
}

func (s *Service) getUserSubscribedTopics(ctx context.Context, userID uuid.UUID) ([]string, error) {
	userSubscribedTopics := make([]string, 0)
	for _, topic := range allTopics {
		enabled, err := s.fr.IsNotificationEnabled(ctx, firebase_repository.IsNotificationEnabledParams{
			UserID: userID,
			Topic:  topic,
		})
		if err != nil {
			return nil, err
		}
		if enabled {
			userSubscribedTopics = append(userSubscribedTopics, topic)
		}
	}

	return userSubscribedTopics, nil
}

func (s *Service) subscribeTokenToTopics(ctx context.Context, token string, topics []string) error {
	for _, topic := range topics {
		_, err := s.messagingClient.SubscribeToTopic(ctx, []string{token}, topic)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) SendNotificationToTopic(ctx context.Context, topicName, title, body string, extraData map[string]string) error {
	data := map[string]string{
		"type":  topicName,
		"title": title,
		"body":  body,
	}
	for k, v := range extraData {
		data[k] = v
	}
	_, err := s.messagingClient.Send(ctx, &messaging.Message{
		Data: data,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		APNS: &messaging.APNSConfig{
			Headers: map[string]string{
				"apns-priority": "10",
			},
		},
		FCMOptions: &messaging.FCMOptions{
			AnalyticsLabel: topicName,
		},
		Topic: topicName,
	})
	if err != nil {
		return errors.Wrapf(err, "can't send notification to %v topic", topicName)
	}

	return nil
}

func (s *Service) SendNewShowNotification(ctx context.Context, showTitle string, showID uuid.UUID) error {
	err := s.SendNotificationToTopic(
		ctx,
		NewShowTopicName,
		fmt.Sprintf("New Arrival"),
		fmt.Sprintf("%v is now on Sator.", showTitle),
		map[string]string{
			"show_id": showID.String(),
		},
	)
	if err != nil {
		return errors.Wrap(err, "can't send notification to topic")
	}

	return nil
}

func (s *Service) SendNewEpisodeNotification(ctx context.Context, showTitle, episodeTitle string, showID, seasonID, episodeID uuid.UUID) error {
	err := s.SendNotificationToTopic(
		ctx,
		NewEpisodeTopicName,
		fmt.Sprintf("New Arrival"),
		fmt.Sprintf("%s: %s is now on Sator.", showTitle, episodeTitle),
		map[string]string{
			"show_id":    showID.String(),
			"seasonID":   seasonID.String(),
			"episode_id": episodeID.String(),
		},
	)
	if err != nil {
		return errors.Wrap(err, "can't send notification to topic")
	}

	return nil
}
