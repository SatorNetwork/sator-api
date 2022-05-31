package gapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAmountMultiplier(t *testing.T) {
	amount := 12.45
	mltpl := 1e9
	res := uint64(amount * float64(mltpl))

	assert.Equal(t, uint64(12450000000), res)
}
