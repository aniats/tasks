package done_channel

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOr(t *testing.T) {
	tests := []struct {
		name            string
		setup           func() []<-chan interface{}
		expectImmediate bool
		maxWaitTime     time.Duration
	}{
		{
			name: "no channels",
			setup: func() []<-chan interface{} {
				return nil
			},
			expectImmediate: true,
		},
		{
			name: "empty slice of channels",
			setup: func() []<-chan interface{} {
				return []<-chan interface{}{}
			},
			expectImmediate: true,
		},
		{
			name: "single channel that closes immediately",
			setup: func() []<-chan interface{} {
				ch := make(chan interface{})
				close(ch)
				return []<-chan interface{}{ch}
			},
			expectImmediate: true,
		},
		{
			name: "single channel that closes after delay",
			setup: func() []<-chan interface{} {
				ch := make(chan interface{})
				go func() {
					time.Sleep(20 * time.Millisecond)
					close(ch)
				}()
				return []<-chan interface{}{ch}
			},
			expectImmediate: false,
			maxWaitTime:     50 * time.Millisecond,
		},
		{
			name: "multiple channels, first closes immediately",
			setup: func() []<-chan interface{} {
				ch1 := make(chan interface{})
				ch2 := make(chan interface{})
				ch3 := make(chan interface{})
				close(ch1)
				return []<-chan interface{}{ch1, ch2, ch3}
			},
			expectImmediate: true,
		},
		{
			name: "multiple channels, second closes first",
			setup: func() []<-chan interface{} {
				ch1 := make(chan interface{})
				ch2 := make(chan interface{})
				ch3 := make(chan interface{})
				go func() {
					time.Sleep(10 * time.Millisecond)
					close(ch2)
				}()
				go func() {
					time.Sleep(30 * time.Millisecond)
					close(ch1)
				}()
				go func() {
					time.Sleep(50 * time.Millisecond)
					close(ch3)
				}()
				return []<-chan interface{}{ch1, ch2, ch3}
			},
			expectImmediate: false,
			maxWaitTime:     25 * time.Millisecond,
		},
		{
			name: "multiple channels, last closes first",
			setup: func() []<-chan interface{} {
				ch1 := make(chan interface{})
				ch2 := make(chan interface{})
				ch3 := make(chan interface{})
				go func() {
					time.Sleep(50 * time.Millisecond)
					close(ch1)
				}()
				go func() {
					time.Sleep(30 * time.Millisecond)
					close(ch2)
				}()
				go func() {
					time.Sleep(10 * time.Millisecond)
					close(ch3)
				}()
				return []<-chan interface{}{ch1, ch2, ch3}
			},
			expectImmediate: false,
			maxWaitTime:     25 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			channels := tt.setup()

			start := time.Now()
			result := Or(channels...)

			if tt.expectImmediate {
				select {
				case <-result:
					elapsed := time.Since(start)
					assert.Less(t, elapsed, 5*time.Millisecond, "Channel should close immediately")
				case <-time.After(10 * time.Millisecond):
					assert.Fail(t, "Or() should have closed immediately")
				}
			} else {
				select {
				case <-result:
					elapsed := time.Since(start)
					assert.Greater(t, elapsed, 5*time.Millisecond, "Channel closed too quickly")
					assert.Less(t, elapsed, tt.maxWaitTime, "Channel took too long to close")
				case <-time.After(tt.maxWaitTime + 10*time.Millisecond):
					assert.Fail(t, "Or() did not close within expected time")
				}
			}

			select {
			case _, ok := <-result:
				assert.False(t, ok, "Channel should be closed")
			default:
				assert.Fail(t, "Channel should be closed and readable")
			}
		})
	}
}
