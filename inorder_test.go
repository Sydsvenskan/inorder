package inorder

import (
	"math/rand"
	"testing"
	"time"

	"github.com/xdg/testy"
)

func DummyTask(d time.Duration) *Block {
	b := NewBlock()
	go func() {
		<-time.NewTimer(d).C
		b.Done()
	}()
	return b
}

func TestFixedOrdering(t *testing.T) {
	is := testy.New(t)
	defer func() { t.Logf(is.Done()) }()

	a := DummyTask(20 * time.Millisecond)
	b := DummyTask(10 * time.Millisecond)
	c := DummyTask(15 * time.Millisecond)

	inOrder := NewInOrder()
	inOrder.Enqueue(a)
	inOrder.Enqueue(b)
	inOrder.Enqueue(c)

	is.Label("task A comes first").True(a == <-inOrder.Done)
	is.Label("task B comes second").True(b == <-inOrder.Done)
	is.Label("task C comes last").True(c == <-inOrder.Done)
}

func TestRandomOrdering(t *testing.T) {
	is := testy.New(t).Label("the order is preserved")
	defer func() { t.Logf(is.Done()) }()

	inOrder := NewInOrder()

	var tasks []*Block
	for i := 0; i < 20; i++ {
		ms := rand.Int() % 20
		task := DummyTask(time.Duration(ms) * time.Millisecond)
		tasks = append(tasks, task)
		inOrder.Enqueue(task)
	}

	for _, task := range tasks {
		is.True(task == <-inOrder.Done)
	}
}
