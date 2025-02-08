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
	taskCount := len(tasks)
	workerCount := min(n, taskCount)
	taskChannel := make(chan Task, workerCount)
	var countErrors atomic.Int64
	// публикация заданий
	go func() {
		for i := 0; i < taskCount; i++ {
			if isCountErrExeed(&countErrors, m) {
				break
			}
			taskChannel <- tasks[i]
		}
		close(taskChannel)
	}()
	// запуск обработчиков заданий
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
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
	wg.Wait()
	if isCountErrExeed(&countErrors, m) {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func isCountErrExeed(countErrors *atomic.Int64, m int) bool {
	return countErrors.Load() >= int64(m)
}
