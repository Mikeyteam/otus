package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks.
func Run(tasks []Task, countParallelGoroutine, countError int) error {
	if countError <= 0 {
		return ErrErrorsLimitExceeded
	}
	channelWithTask := make(chan Task)
	emptyChannel := make(chan struct{})
	waitGroup := &sync.WaitGroup{}
	err := int32(0)

	waitGroup.Add(countParallelGoroutine)

	for i := 0; i < countParallelGoroutine; i++ {
		go workerChannel(channelWithTask, emptyChannel, waitGroup, countError, &err)
	}

	for _, task := range tasks {
		select {
		case <-emptyChannel:
			waitGroup.Wait()
			return ErrErrorsLimitExceeded
		case channelWithTask <- task:
		}
	}

	close(channelWithTask)
	waitGroup.Wait()

	select {
	case <-emptyChannel:
		return ErrErrorsLimitExceeded
	default:
	}

	return nil
}

// workerChannel read task from channel.
func workerChannel(
	channelWithTask chan Task,
	emptyChannel chan struct{},
	waitGroup *sync.WaitGroup,
	countError int,
	errors *int32,
) {
	defer waitGroup.Done()

	for {
		select {
		case <-emptyChannel:
			return
		case task, ok := <-channelWithTask:
			if !ok {
				return
			}

			if err := task(); err != nil {
				if atomic.AddInt32(errors, 1) == int32(countError) {
					close(emptyChannel)
				}
			}
		}
	}
}
