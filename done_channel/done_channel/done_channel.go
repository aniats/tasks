package done_channel

import (
	"context"
	"sync"
)

func Or(channels ...<-chan interface{}) <-chan interface{} {
	if len(channels) == 0 {
		c := make(chan interface{})
		close(c)
		return c
	}

	result := make(chan interface{})
	ctx, cancel := context.WithCancel(context.Background())

	var mu sync.Mutex
	var finished bool

	for _, ch := range channels {
		go func(ch <-chan interface{}) {
			select {
			case <-ch:
				mu.Lock()
				if !finished {
					finished = true
					cancel()
					close(result)
				}
				mu.Unlock()
			case <-ctx.Done():
				return
			}
		}(ch)
	}

	return result
}
