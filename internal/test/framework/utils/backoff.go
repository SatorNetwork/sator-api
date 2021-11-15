package utils

import (
	"testing"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/stretchr/testify/require"
)

func BackoffRetry(t *testing.T, o backoff.Operation) {
	const (
		maxRetries      = 300
		constantBackOff = 200 * time.Millisecond
	)
	err := backoff.Retry(o, backoff.WithMaxRetries(backoff.NewConstantBackOff(constantBackOff), maxRetries))
	require.NoError(t, err)
}
