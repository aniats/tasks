package workerpool

import (
	"errors"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRun_EmptyTasks(t *testing.T) {
	err := Run([]Task{}, 2, 1)
	require.NoError(t, err)
}

func TestRun_AllTasksSuccess(t *testing.T) {
	var counter int64
	tasks := make([]Task, 10)
	for i := range tasks {
		tasks[i] = func() error {
			atomic.AddInt64(&counter, 1)
			return nil
		}
	}

	err := Run(tasks, 3, 5)
	require.NoError(t, err)
	require.Equal(t, int64(10), counter)
}

func TestRun_ErrorsLimitExceeded(t *testing.T) {
	var taskCount int64
	tasks := make([]Task, 10)
	for i := range tasks {
		tasks[i] = func() error {
			atomic.AddInt64(&taskCount, 1)
			return errors.New("test error")
		}
	}

	err := Run(tasks, 2, 3)
	require.ErrorIs(t, err, ErrErrorsLimitExceeded)

	require.Greater(t, taskCount, int64(0))
	require.LessOrEqual(t, taskCount, int64(5))
}

func TestRun_ZeroErrorsAllowed(t *testing.T) {
	tasks := []Task{
		func() error { return errors.New("error") },
	}

	err := Run(tasks, 1, 0)
	require.NoError(t, err)
}

func TestRun_NegativeErrorsAllowed(t *testing.T) {
	var taskCount int64
	tasks := make([]Task, 5)
	for i := range tasks {
		tasks[i] = func() error {
			atomic.AddInt64(&taskCount, 1)
			return errors.New("error")
		}
	}

	err := Run(tasks, 2, -1)
	require.NoError(t, err)
	require.Equal(t, int64(5), taskCount)
}

func TestRun_ZeroWorkers(t *testing.T) {
	var taskCount int64
	tasks := []Task{
		func() error {
			atomic.AddInt64(&taskCount, 1)
			return nil
		},
	}

	err := Run(tasks, 0, 1)
	require.NoError(t, err)
	require.Equal(t, int64(1), taskCount)
}
