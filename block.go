package inorder

import "sync"

// Block provides a wait mechanic based on whether a task is done
type Block struct {
	done  bool
	cond  *sync.Cond
	mutex sync.RWMutex
}

// NewBlock creates a new block
func NewBlock() *Block {
	return &Block{
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

// IsDone checks if the task is done
func (b *Block) IsDone() (done bool) {
	b.mutex.RLock()
	done = b.done
	b.mutex.RUnlock()
	return
}

// Wait blocks until the task is done
func (b *Block) Wait() {
	b.cond.L.Lock()
	for !b.IsDone() {
		b.cond.Wait()
	}
	b.cond.L.Unlock()
}

// Done sets the task as done
func (b *Block) Done() {
	b.mutex.Lock()
	b.done = true
	b.mutex.Unlock()
	b.cond.Broadcast()
}
