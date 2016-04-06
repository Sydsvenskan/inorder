package inorder

import (
	"errors"
	"sync"
	"time"
)

var (
	// ErrTaskTimedOut the task didn't complete within the specified time limit
	ErrTaskTimedOut = errors.New("task timed out")
)

// Task is something that we can wait for
type Task interface {
	Wait() chan bool
}

// Result is the result of a task
type Result struct {
	Error error
	Task  Task
}

// InOrder maintains the order of enqueued jobs
type InOrder struct {
	order []Task
	mutex sync.Mutex

	Timeout time.Duration
	Done    chan *Result
}

// NewInOrder creates a new orderer
func NewInOrder(timeout time.Duration) *InOrder {
	return &InOrder{
		Done:    make(chan *Result),
		Timeout: timeout,
	}
}

// Enqueue a new task
func (in *InOrder) Enqueue(task Task) {
	in.mutex.Lock()
	in.order = append(in.order, task)

	// If the inOrder slice was empty before we want to start up the forwarding
	// routine.
	if len(in.order) == 1 {
		go in.forwardOnDone(task)
	}
	in.mutex.Unlock()
}

func (in *InOrder) forwardOnDone(task Task) {
	select {
	case <-task.Wait():
		in.Done <- &Result{
			Task: task,
		}
	case <-time.After(in.Timeout):
		in.Done <- &Result{
			Task:  task,
			Error: ErrTaskTimedOut,
		}
	}

	in.mutex.Lock()
	in.order = in.order[1:]

	// Keep forwarding until we have emptied the in-order list
	if len(in.order) > 0 {
		go in.forwardOnDone(in.order[0])
	}
	in.mutex.Unlock()
}
