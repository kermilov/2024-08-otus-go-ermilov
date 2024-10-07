package hw05parallelexecution

import (
	"errors"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	taskChannel := publishTasks(tasks)
	return <-runWorkers(n, m, taskChannel)
}

func publishTasks(tasks []Task) chan Task {
	taskChannel := make(chan Task, len(tasks))
	for _, task := range tasks {
		taskChannel <- task
	}
	close(taskChannel)
	return taskChannel
}

func runWorkers(n, m int, taskChannel chan Task) chan error {
	var countErrors, countAll int32
	finishChannel := make(chan error, 1)
	for i := 0; i < n; i++ {
		go worker(taskChannel, finishChannel, &countErrors, &countAll, m)
	}
	return finishChannel
}

func worker(taskChannel chan Task, finishChannel chan error, countErrors, countAll *int32, m int) {
	for task := range taskChannel {
		cur := int(atomic.LoadInt32(countErrors))
		if cur != 0 && cur >= m {
			return
		}
		if err := task(); err != nil {
			inc := int(atomic.AddInt32(countErrors, 1))
			if inc >= m {
				if inc == m || (inc == 1 && m == 0) {
					finishChannel <- ErrErrorsLimitExceeded
				}
				return
			}
		}
		if int(atomic.AddInt32(countAll, 1)) == cap(taskChannel) {
			finishChannel <- nil
		}
	}
}
