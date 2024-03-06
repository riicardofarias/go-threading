package threading

import (
	"context"
	"sync"
	"sync/atomic"
)

type WorkerPool[T any] struct {
	numOfExecutions int32
	numOfFailures   int32

	jobError     error
	workersCount chan int
	wg           *sync.WaitGroup
	ctx          context.Context
	cancel       func()
}

func NewWorkerPool[T any](workersCount int) *WorkerPool[T] {
	ctx, cancel := context.WithCancel(context.Background())

	// If the workersCount is less than or equal to zero, set it to 1.
	if workersCount <= 0 {
		workersCount = 1
	}

	// Create a new worker pool.
	w := &WorkerPool{
		numOfExecutions: int32(0),
		numOfFailures:   int32(0),
		workersCount:    make(chan int, workersCount),
		wg:              &sync.WaitGroup{},
		ctx:             ctx,
		cancel:          cancel,
	}

	// Fill the workersCount channel with the number of workers.
	for i := 0; i < workersCount; i++ {
		w.workersCount <- 1
	}

	return (*WorkerPool[T])(w)
}

func (w *WorkerPool[T]) NumOfExecutions() int32 {
	return atomic.LoadInt32(&w.numOfExecutions)
}

func (w *WorkerPool[T]) NumOfFailures() int32 {
	return atomic.LoadInt32(&w.numOfFailures)
}

func (w *WorkerPool[T]) Error() error {
	return w.jobError
}

// RunJob runs the given job in a worker.
// The jobFn is a function that takes an integer as an argument and returns an error.
// The integer is the id of the worker.
// If the jobFn returns an error, the worker pool is stopped.
func (w *WorkerPool[T]) RunJob(dataset T, jobFn func(_dataset T) error) {
	w.wg.Add(<-w.workersCount)

	// Run the jobFn in a goroutine.
	go func() {
		defer w.wg.Done()
		defer func() { w.workersCount <- 1 }()

		select {
		case <-w.ctx.Done():
			return
		default:
			if err := jobFn(dataset); err != nil {
				w.cancel()
				atomic.AddInt32(&w.numOfFailures, 1)
				w.jobError = err
			}
		}

		atomic.AddInt32(&w.numOfExecutions, 1)
	}()
}

// Wait waits for all the workers to finish.
// It is a blocking function. It should be called after all the jobs have been added to the worker pool.
func (w *WorkerPool[T]) Wait() {
	defer close(w.workersCount)
	w.wg.Wait()
}
