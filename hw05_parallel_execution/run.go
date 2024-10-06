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
	var countErrors int32
	taskCount := len(tasks)
	for start := 0; start < taskCount; start += n {
		end := min(start+n, taskCount)
		run(tasks[start:end], &countErrors, m)
		if int(atomic.LoadInt32(&countErrors)) >= m {
			return ErrErrorsLimitExceeded
		}
	}
	return nil
}

func run(tasks []Task, countErrors *int32, m int) error {
	wg := sync.WaitGroup{}
	for _, task := range tasks {
		wg.Add(1)
		go func(task Task) {
			defer wg.Done()
			if int(atomic.LoadInt32(countErrors)) >= m {
				return
			}
			if err := task(); err != nil {
				atomic.AddInt32(countErrors, 1)
			}
		}(task)
	}
	wg.Wait()
	return nil
}
