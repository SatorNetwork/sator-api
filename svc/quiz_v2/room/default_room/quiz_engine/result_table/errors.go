package result_table

import (
	"fmt"
	"github.com/google/uuid"
)

type ErrRowNotFound struct {
	userID uuid.UUID
}

func NewErrRowNotFound(userID uuid.UUID) *ErrRowNotFound {
	return &ErrRowNotFound{
		userID: userID,
	}
}

func (e *ErrRowNotFound) Error() string {
	return fmt.Sprintf("row not found for user with %v UID", e.userID)
}

type ErrCellNotFound struct {
	userID uuid.UUID
	qNum   int
	rowLen int
}

func NewErrCellNotFound(userID uuid.UUID, qNum, rowLen int) *ErrCellNotFound {
	return &ErrCellNotFound{
		userID: userID,
		qNum:   qNum,
		rowLen: rowLen,
	}
}

func (e *ErrCellNotFound) Error() string {
	return fmt.Sprintf("for user with %v UID qNum must be within [0..%v] internal, but got: %v", e.userID, e.rowLen-1, e.qNum)
}

type ErrIndexOutOfRange struct {
	sliceLen int
	indexNum int
}

func NewErrIndexOutOfRange(sliceLen, indexNum int) *ErrIndexOutOfRange {
	return &ErrIndexOutOfRange{
		sliceLen: sliceLen,
		indexNum: indexNum,
	}
}

func (e *ErrIndexOutOfRange) Error() string {
	return fmt.Sprintf("index must be within [0..%v] internal, but got: %v", e.sliceLen-1, e.indexNum)
}
