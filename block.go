package inorder

// Block provides a wait mechanic based on whether a task is done
type Block struct {
	done chan bool
}

// NewBlock creates a new block
func NewBlock() *Block {
	return &Block{
		done: make(chan bool),
	}
}

// IsDone checks if the task is done
func (b *Block) IsDone() (done bool) {
	select {
	case <-b.done:
		return true
	default:
		return false
	}
}

// Wait returns a channel that is closed when the task is done
func (b *Block) Wait() chan bool {
	return b.done
}

// Done sets the task as done
func (b *Block) Done() {
	close(b.done)
}
