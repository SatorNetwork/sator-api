package status_transactor

import "log"

type Status uint8

const (
	GatheringPlayersStatus Status = iota
	RoomIsFullStatus
	CountdownFinishedStatus
	QuestionAreSentStatus
	WinnersTableAreSent
	RoomIsFinished
	RoomIsClosed
)

var allowedTransitionsMap = map[Status][]Status{
	GatheringPlayersStatus:  {RoomIsFullStatus},
	RoomIsFullStatus:        {CountdownFinishedStatus},
	CountdownFinishedStatus: {QuestionAreSentStatus},
	QuestionAreSentStatus:   {WinnersTableAreSent},
	WinnersTableAreSent:     {RoomIsFinished},
	RoomIsFinished:          {RoomIsClosed},
	RoomIsClosed:            {},
}

func isTransitionAllowed(from, to Status) bool {
	allowedTransitions, ok := allowedTransitionsMap[from]
	if !ok {
		return false
	}
	for _, allowedTransition := range allowedTransitions {
		if to == allowedTransition {
			return true
		}
	}

	return false
}

type StatusTransactor struct {
	status Status

	notifyChan chan struct{}
}

func New(notifyChan chan struct{}) *StatusTransactor {
	return &StatusTransactor{
		status:     GatheringPlayersStatus,
		notifyChan: notifyChan,
	}
}

func (st *StatusTransactor) SetStatus(newStatus Status) {
	if st.status == newStatus {
		return
	}

	if !isTransitionAllowed(st.status, newStatus) {
		log.Printf("transition from %v to %v isn't allowed\n", st.status, newStatus)
		return
	}

	st.status = newStatus
	if st.status != RoomIsClosed {
		st.notifyChan <- struct{}{}
	}
}

func (st *StatusTransactor) GetStatus() Status {
	return st.status
}
