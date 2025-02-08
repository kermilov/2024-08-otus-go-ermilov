package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	x := min(n, len(tasks))
	taskChannel := make(chan Task, x)
	var wg sync.WaitGroup
	var countErrors atomic.Int64
	for i := 0; i < x; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskChannel {
				if isCountErrExeed(&countErrors, m) {
					return
				}
				err := task()
				if err != nil {
					countErrors.Add(1)
				}
			}
		}()
	}
	taskCount := len(tasks)
	for i := 0; i < taskCount; i++ {
		if isCountErrExeed(&countErrors, m) {
			break
		}
		taskChannel <- tasks[i]
	}
	close(taskChannel)
	wg.Wait()
	if isCountErrExeed(&countErrors, m) {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func isCountErrExeed(countErrors *atomic.Int64, m int) bool {
	return countErrors.Load() >= int64(m)
}
