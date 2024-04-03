package job

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
	"visa/global"

	"go.uber.org/zap"
)

type ExecutionQueue struct {
	queue    chan StepFunction
	wg       sync.WaitGroup
	Ctx      context.Context
	cancel   context.CancelFunc
	mu       sync.Mutex
	Deadline time.Time
	abort    bool
}

type StepFunction struct {
	Function         func(ctx context.Context) (err error)
	FunctionWithArgs func(ctx context.Context, args ...interface{}) (err error)
	Args             []interface{}
}

func NewExecutionQueue(ctx context.Context, cancel context.CancelFunc, deadline time.Time) *ExecutionQueue {
	return &ExecutionQueue{
		queue:    make(chan StepFunction),
		Ctx:      ctx,
		cancel:   cancel,
		Deadline: deadline,
	}
}

func (eq *ExecutionQueue) Run() {
	for fn := range eq.queue {
		fmt.Println("fn: ", fn)
		fmt.Println("eq.queue: ", eq.queue)
		eq.mu.Lock()
		if eq.abort {
			eq.mu.Unlock()
			break
		}
		eq.mu.Unlock()
		// Execute the function
		if len(fn.Args) == 0 {
			if err := fn.Function(eq.Ctx); err != nil {
				eq.Abort(err)
			}
		} else {
			if err := fn.FunctionWithArgs(eq.Ctx, fn.Args...); err != nil {
				eq.Abort(err)
			}
		}

		// time.Sleep(5 * time.Second)
		eq.wg.Done()
	}
}

func (eq *ExecutionQueue) Close() {
	close(eq.queue)
	eq.cancel()
}

func (eq *ExecutionQueue) Wait() {
	eq.wg.Wait()
}

func (eq *ExecutionQueue) Enqueue(fn StepFunction) {
	if time.Now().Before(eq.Deadline) {
		eq.wg.Add(1)
		eq.queue <- fn
	} else {
		fmt.Println("enquque deadline comes")
		eq.Abort(errors.New("deadline comes"))
	}
}

func (eq *ExecutionQueue) Abort(err error) {
	global.LOGGER.Error("preview article validation error", zap.Error(err))
	eq.mu.Lock()
	defer eq.mu.Unlock()

	// Signal an abort
	eq.abort = true

	// Close the channel and cancel context
	eq.Close()
}
