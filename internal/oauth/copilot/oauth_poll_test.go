package copilot

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCopilotPollDelay(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name     string
		interval int
		want     time.Duration
	}{
		{name: "minimum interval", interval: 3, want: 7 * time.Second},
		{name: "server interval", interval: 5, want: 7 * time.Second},
		{name: "slow down interval", interval: 10, want: 12 * time.Second},
		{name: "second slow down interval", interval: 15, want: 17 * time.Second},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tt.want, copilotPollDelay(tt.interval))
		})
	}
}

func TestPollForTokenCanceledBeforePoll(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	start := time.Now()
	_, err := PollForToken(ctx, &DeviceCode{ExpiresIn: 300, Interval: 5})
	require.ErrorIs(t, err, context.Canceled)
	require.Less(t, time.Since(start), time.Second)
}

func TestPollForTokenExpiredBeforePoll(t *testing.T) {
	t.Parallel()

	_, err := PollForToken(context.Background(), &DeviceCode{ExpiresIn: 0, Interval: 5})
	require.ErrorContains(t, err, "authorization timed out")
}
