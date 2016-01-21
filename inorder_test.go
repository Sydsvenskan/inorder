package inorder_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/Sydsvenskan/inorder"
	"github.com/xdg/testy"
)

type DummyTask struct {
	block *inorder.Block
	Name  string
}

func (dt *DummyTask) Wait() {
	dt.block.Wait()
}

func NewDummyTask(d time.Duration, name string) *DummyTask {
	t := &DummyTask{
		block: inorder.NewBlock(),
	}
	go func() {
		<-time.NewTimer(d).C
		t.Name = name
		t.block.Done()
	}()
	return t
}

func ExampleUse() {
	a := NewDummyTask(20*time.Millisecond, "a")
	b := NewDummyTask(10*time.Millisecond, "b")
	c := NewDummyTask(15*time.Millisecond, "c")

	inOrder := inorder.NewInOrder()
	inOrder.Enqueue(a)
	inOrder.Enqueue(b)
	inOrder.Enqueue(c)

	fmt.Println("first:", (<-inOrder.Done).(*DummyTask).Name)
	fmt.Println("second:", (<-inOrder.Done).(*DummyTask).Name)
	fmt.Println("last:", (<-inOrder.Done).(*DummyTask).Name)
	// Output:
	// first: a
	// second: b
	// last: c
}

func TestRandomOrdering(t *testing.T) {
	is := testy.New(t).Label("the order is preserved")
	defer func() { t.Logf(is.Done()) }()

	inOrder := inorder.NewInOrder()

	var tasks []*DummyTask
	for i := 0; i < 20; i++ {
		ms := rand.Int() % 20
		task := NewDummyTask(time.Duration(ms)*time.Millisecond, "")
		tasks = append(tasks, task)
		inOrder.Enqueue(task)
	}

	for _, task := range tasks {
		is.True(task == <-inOrder.Done)
	}
}
