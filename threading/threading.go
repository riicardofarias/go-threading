package threading

import (
	"context"
	"sync"
	"sync/atomic"
)

type WorkerPool struct {
	WorkersNumber   chan int
	WorkerGroup     sync.WaitGroup
	NumOfExecutions int32

	ctx    context.Context
	cancel func()
}

func NewWorkerPool(workersNumber int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	w := &WorkerPool{
		WorkersNumber: make(chan int, workersNumber),
		WorkerGroup:   sync.WaitGroup{},

		ctx:    ctx,
		cancel: cancel,
	}

	for i := 0; i < workersNumber; i++ {
		w.WorkersNumber <- 1
	}

	return w
}

func (w *WorkerPool) RunJob(id int, job func(num int) error) {
	num := <-w.WorkersNumber

	w.WorkerGroup.Add(num)

	go func(ctx context.Context) {
		defer w.WorkerGroup.Done()
		defer func() { w.WorkersNumber <- 1 }()

		select {
		case <-w.ctx.Done():
			return
		default:
			if err := job(id); err != nil {
				w.cancel()
			}
		}

		atomic.AddInt32(&w.NumOfExecutions, 1)
	}(w.ctx)
}

func (w *WorkerPool) Wait() {
	defer close(w.WorkersNumber)
	w.WorkerGroup.Wait()
}
