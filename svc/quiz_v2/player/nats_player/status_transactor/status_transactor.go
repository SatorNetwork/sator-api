package status_transactor

import "sync"

type Status uint8

const (
	UndefinedStatus Status = iota
	PlayerDisconnectedStatus
	PlayerConnectedStatus
)

type StatusTransactor struct {
	status      Status
	statusMutex *sync.Mutex

	notifyChan chan struct{}
}

func New(notifyChan chan struct{}) *StatusTransactor {
	return &StatusTransactor{
		status:      UndefinedStatus,
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
	st.status = newStatus

	st.notifyChan <- struct{}{}
}

func (st *StatusTransactor) GetStatus() Status {
	st.statusMutex.Lock()
	defer st.statusMutex.Unlock()

	return st.status
}
