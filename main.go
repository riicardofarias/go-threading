package main

import (
	"errors"
	"fmt"
	"go-threading/threading"
	"time"
)

const (
	dataCount = 200
)

func main() {
	workerPool := threading.NewWorkerPool(20)

	for i := 0; i < dataCount; i++ {
		workerPool.RunJob(i, func(num int) error {
			err := Worker(num)
			return err
		})
	}

	workerPool.Wait()

	fmt.Println(fmt.Sprintf("Number of executions: %d", workerPool.NumOfExecutions()))
	fmt.Println(fmt.Sprintf("Number of failures: %d", workerPool.NumOfFailures()))
}

func Worker(id int) error {
	fmt.Println(fmt.Sprintf("Worker %d started", id))

	if id == 100 || id == 101 {
		return errors.New("error")
	}

	time.Sleep(500 * time.Millisecond)

	return nil
}
