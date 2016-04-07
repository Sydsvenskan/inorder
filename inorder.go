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
	IsDone() bool
}

// Result is the result of a task
type Result struct {
	mutex *sync.Mutex
	Error error
	Task  Task
}

func (r *Result) SetError(err error) {
	r.mutex.Lock()
	r.Error = err
	r.mutex.Unlock()
}

func (r *Result) IsDone() (bool, error) {
	var err error

	r.mutex.Lock()
	if r.Error != nil {
		err = r.Error
	}
	r.mutex.Unlock()

	if err != nil {
		return false, err
	}
	return r.Task.IsDone(), nil
}

// InOrder maintains the order of enqueued jobs
type InOrder struct {
	order    []*Result
	mutex    sync.Mutex
	taskDone chan bool

	Timeout time.Duration
	Done    chan *Result
}

// NewInOrder creates a new orderer
func NewInOrder(timeout time.Duration) *InOrder {
	in := &InOrder{
		Done:     make(chan *Result),
		taskDone: make(chan bool),
		Timeout:  timeout,
	}

	go func() {
		var doneList []*Result

		for _ = range in.taskDone {
			in.mutex.Lock()

			doneUntil := -1
			for i, res := range in.order {
				done, err := res.IsDone()
				if done || err != nil {
					doneList = append(doneList, res)
				} else {
					break
				}
				doneUntil = i
			}
			if doneUntil >= 0 {
				in.order = in.order[doneUntil+1:]
			}
			in.mutex.Unlock()

			// Send off the the done tasks outside of the mutex lock
			if len(doneList) > 0 {
				for _, res := range doneList {
					in.Done <- res
				}
				doneList = doneList[:0]
			}
		}
	}()

	return in
}

// Enqueue a new task
func (in *InOrder) Enqueue(task Task) {
	in.mutex.Lock()

	result := &Result{
		mutex: &sync.Mutex{},
		Task:  task,
	}
	in.order = append(in.order, result)

	in.mutex.Unlock()

	go func() {
		select {
		case <-result.Task.Wait():
		case <-time.After(in.Timeout):
			result.Error = ErrTaskTimedOut
		}
		in.taskDone <- true
	}()
}
