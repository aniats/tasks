package mutex

import (
	"context"
	"sync"
)

type Semaphore interface {
	Acquire(context.Context, int64) error
	TryAcquire(int64) bool
	Release(int64)
}

type MutexSemaphore struct {
	mu         *sync.Mutex // Защищает доступ к permits
	cond       *sync.Cond  // Для пробуждения ожидающих горутин
	permits    int64       // Сколько разрешений доступно СЕЙЧАС
	maxPermits int64       // Максимальное количество разрешений
}

func NewMutexSemaphore(capacity int64) *MutexSemaphore {
	s := &MutexSemaphore{
		permits:    capacity,
		maxPermits: capacity,
		mu:         &sync.Mutex{},
	}
	s.cond = sync.NewCond(s.mu) // Привязываем condition variable к мьютексу
	return s
}

func (s *MutexSemaphore) TryAcquire(permits int64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.permits >= permits {
		s.permits -= permits
		return true
	}

	return false
}

func (s *MutexSemaphore) Release(permits int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.permits += permits

	if s.permits > s.maxPermits {
		s.permits = s.maxPermits
	}

	s.cond.Broadcast()
}

func (s *MutexSemaphore) Acquire(ctx context.Context, permits int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	done := make(chan struct{})

	go func() {
		select {
		case <-ctx.Done():
			s.cond.Broadcast()
		case <-done:
		}
	}()

	for s.permits < permits {
		if ctx.Err() != nil {
			close(done)
			return ctx.Err()
		}
		s.cond.Wait()
	}

	if ctx.Err() != nil {
		close(done)
		return ctx.Err()
	}

	s.permits -= permits
	close(done)

	return nil
}
