package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrWorkersCount        = errors.New("errors workers count <= 0")
)

type Task func() error

func writeTasksToChannel(tasks []Task, taskCh chan<- Task) {
	defer close(taskCh)
	for _, oneTask := range tasks {
		taskCh <- oneTask
	}
}

func readChanel(taskCh <-chan Task, errorsCount *int32, m int) {
	for {
		taskValue, ok := <-taskCh
		if !ok {
			return
		}

		err := taskValue()
		if err != nil {
			errorsCount := atomic.AddInt32(errorsCount, 1)
			if m > 0 && int(errorsCount) >= m {
				return
			}
		}
	}
}

func Run(tasks []Task, n, m int) error {
	if n <= 0 {
		return ErrWorkersCount
	}

	wg := sync.WaitGroup{}
	taskCh := make(chan Task, len(tasks))
	var errorsCount int32

	go writeTasksToChannel(tasks, taskCh)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			readChanel(taskCh, &errorsCount, m)
		}()
	}
	wg.Wait()

	if errorsCount >= int32(m) && m > 0 {
		return ErrErrorsLimitExceeded
	}
	return nil
}
