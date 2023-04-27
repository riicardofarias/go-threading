package threading

import (
	"context"
	"sync"
	"sync/atomic"
)

// WorkerPool is a pool of workers.
// It is used to run jobs in parallel.
type WorkerPool struct {
	numOfExecutions int32

	workersCount chan int
	wg           *sync.WaitGroup
	ctx          context.Context
	cancel       func()
}

// NewWorkerPool creates a new worker pool with the given number of workers.
// The workersCount argument is the number of workers that can run in parallel.
func NewWorkerPool(workersCount int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	// Create a new worker pool.
	w := &WorkerPool{
		numOfExecutions: int32(0),
		workersCount:    make(chan int, workersCount),
		wg:              &sync.WaitGroup{},
		ctx:             ctx,
		cancel:          cancel,
	}

	// Fill the workersCount channel with the number of workers.
	for i := 0; i < workersCount; i++ {
		w.workersCount <- 1
	}

	return w
}

// NumOfExecutions returns the number of jobs that have been executed.
// It is a thread-safe function.
func (w *WorkerPool) NumOfExecutions() int32 {
	return atomic.LoadInt32(&w.numOfExecutions)
}

// RunJob runs the given job in a worker.
// The jobFn is a function that takes an integer as an argument and returns an error.
// The integer is the id of the worker.
// If the jobFn returns an error, the worker pool is stopped.
func (w *WorkerPool) RunJob(id int, jobFn func(num int) error) {
	w.wg.Add(<-w.workersCount)

	// Run the jobFn in a goroutine.
	go func() {
		defer w.wg.Done()
		defer func() { w.workersCount <- 1 }()

		select {
		case <-w.ctx.Done():
			return
		default:
			if err := jobFn(id); err != nil {
				w.cancel()
			}
		}

		atomic.AddInt32(&w.numOfExecutions, 1)
	}()
}

// Wait waits for all the workers to finish.
// It is a blocking function and it should be called after all the jobs have been added to the worker pool.
func (w *WorkerPool) Wait() {
	defer close(w.workersCount)
	w.wg.Wait()
}
