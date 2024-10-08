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
	taskChannel := publishTasks(tasks)
	return runWorkers(n, m, taskChannel)
}

func publishTasks(tasks []Task) chan Task {
	taskChannel := make(chan Task, len(tasks))
	for _, task := range tasks {
		taskChannel <- task
	}
	close(taskChannel)
	return taskChannel
}

func runWorkers(n, m int, taskChannel chan Task) error {
	var countErrors int32
	wg := &sync.WaitGroup{}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(taskChannel, &countErrors, m)
		}()
	}
	wg.Wait()
	if int(countErrors) >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func worker(taskChannel chan Task, countErrors *int32, m int) {
	for task := range taskChannel {
		if err := task(); err != nil {
			if int(atomic.AddInt32(countErrors, 1)) >= m {
				return
			}
		}
	}
}
