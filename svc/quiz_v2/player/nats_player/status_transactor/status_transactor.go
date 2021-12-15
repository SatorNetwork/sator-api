package status_transactor

type Status uint8

const (
	PlayerDisconnectedStatus Status = iota
	PlayerConnectedStatus
)

type StatusTransactor struct {
	status Status

	notifyChan chan struct{}
}

func New(notifyChan chan struct{}) *StatusTransactor {
	return &StatusTransactor{
		status:     PlayerDisconnectedStatus,
		notifyChan: notifyChan,
	}
}

func (st *StatusTransactor) SetStatus(newStatus Status) {
	if st.status == newStatus {
		return
	}
	st.status = newStatus

	st.notifyChan <- struct{}{}
}

func (st *StatusTransactor) GetStatus() Status {
	return st.status
}
