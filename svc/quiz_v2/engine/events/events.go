package events

type Event interface{}

type ForgetRoomEvent struct {
	ChallengeID string
}

func NewForgetRoomEvent(challengeID string) *ForgetRoomEvent {
	return &ForgetRoomEvent{
		ChallengeID: challengeID,
	}
}
