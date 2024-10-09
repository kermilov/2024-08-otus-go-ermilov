package hw05parallelexecution

import (
	"errors"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	taskChannel := make(chan Task)
	errorChannel := make(chan error)
	for i := 0; i < n; i++ {
		go worker(taskChannel, errorChannel)
	}
	var countErrors int32
	var result error
	taskCount := len(tasks)
	for start := 0; start < taskCount; start += n {
		end := min(start+n, taskCount)
		for _, task := range tasks[start:end] {
			taskChannel <- task
		}
		for i := start; i < end; i++ {
			if err := <-errorChannel; err != nil {
				countErrors++
			}
		}
		if countErrors >= int32(m) {
			result = ErrErrorsLimitExceeded
			break
		}
	}
	close(taskChannel)
	close(errorChannel)
	return result
}

func worker(taskChannel chan Task, errorChannel chan error) {
	for task := range taskChannel {
		errorChannel <- task()
	}
}
