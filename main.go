package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

const (
	dataCount = 1_000_000
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

func main() {
	workerPool := NewWorkerPool(150)

	data := make([]int, dataCount)
	chunkSize := 1500
	i := 0

	for i := 0; i < dataCount; i++ {
		data[i] = i
	}

	for true {
		start := i * chunkSize
		end := (i + 1) * chunkSize

		if start >= dataCount {
			break
		}

		if end > dataCount {
			end = dataCount
		}

		//content := data[start:end]

		workerPool.RunJob(i, func(num int) error {
			err := Worker(num)
			return err
		})

		i++
	}

	workerPool.Wait()

	fmt.Println(fmt.Sprintf("Number of executions: %d", workerPool.NumOfExecutions))
}

func Worker(id int) error {
	fmt.Println(fmt.Sprintf("Worker %d started", id))

	if id == rand.Intn(2) {
		return errors.New("error")
	}

	time.Sleep(3 * time.Second)

	return nil
}
