package status_transactor

import (
	"log"
	"sync"
)

type Status uint8

const (
	GatheringPlayersStatus Status = iota
	RoomIsFullStatus
	CountdownFinishedStatus
	QuestionAreSentStatus
	WinnersTableAreSent
	RewardsAreSent
	RoomIsFinished
	RoomIsClosed
)

func (s Status) String() string {
	switch s {
	case GatheringPlayersStatus:
		return "gathering_players_status"
	case RoomIsFullStatus:
		return "room_is_full_status"
	case CountdownFinishedStatus:
		return "countdown_finished_status"
	case QuestionAreSentStatus:
		return "question_are_sent_status"
	case WinnersTableAreSent:
		return "winners_table_are_sent"
	case RewardsAreSent:
		return "rewards_are_sent"
	case RoomIsFinished:
		return "room_is_finished"
	case RoomIsClosed:
		return "room_is_closed"
	default:
		return "gathering_players_status"
	}
}

var allowedTransitionsMap = map[Status][]Status{
	GatheringPlayersStatus:  {RoomIsFullStatus},
	RoomIsFullStatus:        {CountdownFinishedStatus},
	CountdownFinishedStatus: {QuestionAreSentStatus},
	QuestionAreSentStatus:   {WinnersTableAreSent},
	WinnersTableAreSent:     {RewardsAreSent},
	RewardsAreSent:          {RoomIsFinished},
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
	status      Status
	statusMutex *sync.Mutex

	notifyChan chan struct{}
}

func New(notifyChan chan struct{}) *StatusTransactor {
	return &StatusTransactor{
		status:      GatheringPlayersStatus,
		statusMutex: &sync.Mutex{},
		notifyChan:  notifyChan,
	}
}

func (st *StatusTransactor) SetStatus(newStatus Status) {
	st.statusMutex.Lock()
	defer st.statusMutex.Unlock()

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
	st.statusMutex.Lock()
	defer st.statusMutex.Unlock()

	return st.status
}
