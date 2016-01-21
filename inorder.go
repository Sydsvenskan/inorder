package inorder

import "sync"

// Task is something that we can wait for
type Task interface {
	Wait()
}

// InOrder maintains the order of enqueued jobs
type InOrder struct {
	order []Task
	mutex sync.Mutex

	Done chan Task
}

// NewInOrder creates a new orderer
func NewInOrder() *InOrder {
	return &InOrder{
		Done: make(chan Task),
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
	task.Wait()
	in.Done <- task

	in.mutex.Lock()
	in.order = in.order[1:]

	// Keep forwarding until we have emptied the in-order list
	if len(in.order) > 0 {
		go in.forwardOnDone(in.order[0])
	}
	in.mutex.Unlock()
}
