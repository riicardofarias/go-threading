package main

import (
	"errors"
	"fmt"
	"go-threading/threading"
	"time"
)

const (
	dataCount = 1_000_000
)

func main() {
	workerPool := threading.NewWorkerPool(1)

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

	if id == 10 {
		return errors.New("error")
	}

	time.Sleep(100 * time.Millisecond)

	return nil
}
