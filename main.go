package main

import (
	"fmt"
	"go-threading/threading"
	"time"
)

const (
	dataCount = 667
)

func main() {
	workerPool := threading.NewWorkerPool(150)

	for i := 0; i < dataCount; i++ {
		workerPool.RunJob(i, func(num int) error {
			err := Worker(num)
			return err
		})
	}

	workerPool.Wait()

	fmt.Println(fmt.Sprintf("Number of executions: %d", workerPool.NumOfExecutions()))
}

func Worker(id int) error {
	fmt.Println(fmt.Sprintf("Worker %d started", id))

	time.Sleep(5 * time.Second)

	return nil
}
