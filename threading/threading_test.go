package threading

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWorkerPool_RunJob_Success(t *testing.T) {
	t.Run("Should execute all jobs", func(t *testing.T) {
		wp := NewWorkerPool(1)
		numberOfExec := 5

		for i := 0; i < numberOfExec; i++ {
			wp.RunJob(i, func(num int) error {
				return nil
			})
		}

		wp.Wait()
		assert.Equal(t, numberOfExec, int(wp.NumOfExecutions()))
	})
}

func TestWorkerPool_RunJob_Error(t *testing.T) {
	t.Run("Should fail on third run", func(t *testing.T) {
		wp := NewWorkerPool(1)
		numberOfExec := 5

		for i := 0; i < numberOfExec; i++ {
			wp.RunJob(i, func(num int) error {
				if num == 2 {
					return assert.AnError
				}

				return nil
			})
		}

		wp.Wait()
		assert.Equal(t, 3, int(wp.NumOfExecutions()))
	})
}
