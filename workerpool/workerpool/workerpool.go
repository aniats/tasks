package workerpool

import (
	"errors"
	"sync"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrInvalidWorkerCount  = errors.New("worker count must be greater than zero")
	ErrWorkerPoolClosed    = errors.New("worker pool is closed")
	ErrWorkerPoolFull      = errors.New("worker pool is full")
	ErrTaskInvalid         = errors.New("task invalid")
)

type Task func() error

type WorkerPool struct {
	mutex    sync.Mutex
	doneChan chan struct{}
	taskChan chan Task
	isClosed bool
	//workerCount int
}

func New(workerCount int) (*WorkerPool, error) {
	if workerCount < 1 {
		return nil, ErrInvalidWorkerCount
	}

	workerPool := &WorkerPool{
		//workerCount: workerCount,
		mutex:    sync.Mutex{},
		isClosed: false,
		taskChan: make(chan Task, workerCount),
		doneChan: make(chan struct{}),
	}

	go workerPool.Process()
}

func (workerPool *WorkerPool) Process() {
	wg := &sync.WaitGroup{}
	wg.Add(len(workerPool.taskChan))

	for i := 0; i < len(workerPool.taskChan); i++ {
		go func() {
			defer wg.Done()
			// when taskChan would be closed, we will get out of range
			for task := range workerPool.taskChan {
				task()
			}
		}()
	}

	wg.Wait()
	close(workerPool.doneChan)
}

func (workerPool *WorkerPool) AddTask(task Task) error {
	if task == nil {
		return ErrTaskInvalid
	}

	workerPool.mutex.Lock()
	defer workerPool.mutex.Unlock()
	// what if we have multiple pods accessing here some db connection... we're screwed, right?
	if workerPool.isClosed {
		return ErrWorkerPoolClosed
	}

	select {
	case workerPool.taskChan <- task:
		return nil
	case <-workerPool.doneChan:
		///
	default:
		return ErrWorkerPoolFull
	}
}

func (workerPool *WorkerPool) Close() {}
