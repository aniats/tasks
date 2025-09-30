package channel

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChannelSemaphore_TryAcquire(t *testing.T) {
	tests := []struct {
		name        string
		capacity    int64
		permits     int64
		wantSuccess bool
		description string
	}{
		{
			name:        "acquire within capacity",
			capacity:    3,
			permits:     2,
			wantSuccess: true,
			description: "should succeed when permits <= capacity",
		},
		{
			name:        "acquire exactly capacity",
			capacity:    3,
			permits:     3,
			wantSuccess: true,
			description: "should succeed when permits == capacity",
		},
		{
			name:        "acquire beyond capacity",
			capacity:    3,
			permits:     4,
			wantSuccess: false,
			description: "should fail when permits > capacity",
		},
		{
			name:        "acquire zero permits",
			capacity:    3,
			permits:     0,
			wantSuccess: true,
			description: "should succeed when permits == 0",
		},
		{
			name:        "acquire single permit",
			capacity:    1,
			permits:     1,
			wantSuccess: true,
			description: "should succeed with single permit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sem := NewChannelSemaphore(tt.capacity)
			success := sem.TryAcquire(tt.permits)
			assert.Equal(t, tt.wantSuccess, success, tt.description)
		})
	}
}

func TestChannelSemaphore_Release(t *testing.T) {
	tests := []struct {
		name           string
		capacity       int64
		acquireFirst   int64
		releasePermits int64
		description    string
	}{
		{
			name:           "release after acquire",
			capacity:       3,
			acquireFirst:   2,
			releasePermits: 2,
			description:    "should release permits back to semaphore",
		},
		{
			name:           "partial release",
			capacity:       5,
			acquireFirst:   3,
			releasePermits: 1,
			description:    "should partially release permits",
		},
		{
			name:           "release more than acquired",
			capacity:       3,
			acquireFirst:   2,
			releasePermits: 5,
			description:    "should handle releasing more than acquired",
		},
		{
			name:           "release without acquire",
			capacity:       3,
			acquireFirst:   0,
			releasePermits: 2,
			description:    "should handle release without prior acquire",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sem := NewChannelSemaphore(tt.capacity)

			if tt.acquireFirst > 0 {
				success := sem.TryAcquire(tt.acquireFirst)
				require.True(t, success, "initial acquire should succeed")
			}

			assert.NotPanics(t, func() {
				sem.Release(tt.releasePermits)
			}, tt.description)
		})
	}
}

func TestChannelSemaphore_Acquire(t *testing.T) {
	tests := []struct {
		name        string
		capacity    int64
		permits     int64
		timeout     time.Duration
		wantErr     bool
		description string
	}{
		{
			name:        "acquire within capacity",
			capacity:    3,
			permits:     2,
			timeout:     100 * time.Millisecond,
			wantErr:     false,
			description: "should succeed immediately when permits available",
		},
		{
			name:        "acquire with timeout",
			capacity:    1,
			permits:     2,
			timeout:     50 * time.Millisecond,
			wantErr:     true,
			description: "should timeout when not enough permits",
		},
		{
			name:        "acquire zero permits",
			capacity:    3,
			permits:     0,
			timeout:     100 * time.Millisecond,
			wantErr:     false,
			description: "should succeed immediately with zero permits",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sem := NewChannelSemaphore(tt.capacity)
			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			defer cancel()

			err := sem.Acquire(ctx, tt.permits)

			if tt.wantErr {
				assert.Error(t, err, tt.description)
				assert.Equal(t, context.DeadlineExceeded, err)
			} else {
				assert.NoError(t, err, tt.description)
			}
		})
	}
}
