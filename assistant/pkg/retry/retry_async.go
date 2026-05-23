package retry

import (
	"context"
	"log"
	"sync"
)

type AsyncTask[T any] struct {
	Fn     func() (T, error)
	Config Config
	Done   chan Result[T]
}

type Result[T any] struct {
	Value T
	Err   error
}

type RetryAsyncRunner struct {
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	taskChan chan any
	logger   *log.Logger
}

func NewRetryAsyncRunner(buffer int, logger *log.Logger) *RetryAsyncRunner {
	ctx, cancel := context.WithCancel(context.Background())
	return &RetryAsyncRunner{
		ctx:      ctx,
		cancel:   cancel,
		taskChan: make(chan any, buffer),
		logger:   logger,
	}
}

func (r *RetryAsyncRunner) Submit(task AsyncTask[any]) {
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()

		select {
		case <-r.ctx.Done():
			if task.Config.Logger != nil {
				task.Config.Logger.Println("[retry async] runner canceled")
			}
			return
		case r.taskChan <- task:
		}
	}()
}

func (r *RetryAsyncRunner) Start(workers int) {
	for i := 0; i < workers; i++ {
		go r.worker(i)
	}
}

func (r *RetryAsyncRunner) worker(id int) {
	for {
		select {
		case <-r.ctx.Done():
			return
		case t := <-r.taskChan:
			switch task := t.(type) {
			case AsyncTask[any]:
				result, err := Do(r.ctx, task.Config, func() (any, error) {
					return task.Fn()
				})
				task.Done <- Result[any]{Value: result, Err: err}
			default:
				r.logger.Println("[retry async] unknown task type")
			}
		}
	}
}

func (r *RetryAsyncRunner) Stop() {
	r.cancel()
	r.wg.Wait()
}
