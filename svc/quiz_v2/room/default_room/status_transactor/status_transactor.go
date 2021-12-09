package status_transactor

//type status uint8
//
//const (
//	gatheringPlayersStatus = iota
//	countdownStatus
//	sendingQuestionStatus
//)

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
	st.status = newStatus

	if st.status != RoomIsClosed {
		st.notifyChan <- struct{}{}
	}
}

func (st *StatusTransactor) GetStatus() Status {
	return st.status
}
