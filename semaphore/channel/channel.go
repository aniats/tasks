package channel

import "context"

type Semaphore interface {
	Acquire(context.Context, int64) error
	TryAcquire(int64) bool
	Release(int64)
}

type ChannelSemaphore struct {
	ch chan struct{}
}

func NewChannelSemaphore(capacity int64) *ChannelSemaphore {
	return &ChannelSemaphore{
		ch: make(chan struct{}, capacity),
	}
}

func (s *ChannelSemaphore) TryAcquire(permits int64) bool {
	for i := int64(0); i < permits; i++ {
		select {
		case s.ch <- struct{}{}:
		default:
			for j := int64(0); j < i; j++ {
				<-s.ch
			}
			return false
		}
	}
	return true
}

func (s *ChannelSemaphore) Release(permits int64) {
	for i := int64(0); i < permits; i++ {
		for i := int64(0); i < permits; i++ {
			select {
			case <-s.ch:
			default:
				return
			}
		}
	}
}

func (s *ChannelSemaphore) Acquire(ctx context.Context, permits int64) error {
	for i := int64(0); i < permits; i++ {
		select {
		case s.ch <- struct{}{}:
		case <-ctx.Done():
			for j := int64(0); j < i; j++ {
				<-s.ch
			}
			return ctx.Err()
		}
	}
	return nil
}
