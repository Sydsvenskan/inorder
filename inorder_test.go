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
	block        *inorder.Block
	ExpectedName string
	Name         string
}

func (dt *DummyTask) Wait() chan bool {
	return dt.block.Wait()
}

func (dt *DummyTask) IsDone() bool {
	return dt.block.IsDone()
}

func NewDummyTask(d time.Duration, name string) *DummyTask {
	t := &DummyTask{
		ExpectedName: name,
		block:        inorder.NewBlock(),
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

	inOrder := inorder.NewInOrder(50 * time.Millisecond)
	inOrder.Enqueue(a)
	inOrder.Enqueue(b)
	inOrder.Enqueue(c)

	fmt.Println("first:", (<-inOrder.Done).Task.(*DummyTask).Name)
	fmt.Println("second:", (<-inOrder.Done).Task.(*DummyTask).Name)
	fmt.Println("last:", (<-inOrder.Done).Task.(*DummyTask).Name)
	// Output:
	// first: a
	// second: b
	// last: c
}

func TestRandomOrdering(t *testing.T) {
	is := testy.New(t).Label("the order is preserved")
	defer func() { t.Logf(is.Done()) }()

	inOrder := inorder.NewInOrder(100 * time.Millisecond)

	var tasks []*DummyTask
	for i := 0; i < 5000; i++ {
		name := fmt.Sprintf("normal-%d", i)
		ms := rand.Int() % 50
		if i%500 == 0 {
			ms = 20000
			name = fmt.Sprintf("too-long-%d", i)
		}
		task := NewDummyTask(time.Duration(ms)*time.Millisecond, name)
		tasks = append(tasks, task)
		inOrder.Enqueue(task)
		if i%100 == 0 {
			time.Sleep(20 * time.Millisecond)
		}
	}

	for i := range tasks {
		result := <-inOrder.Done
		if result.Error == inorder.ErrTaskTimedOut {
			is.Equal(result.Task.(*DummyTask).ExpectedName, fmt.Sprintf("too-long-%d", i))
		} else if result.Error == nil {
			is.Equal(result.Task.(*DummyTask).Name, fmt.Sprintf("normal-%d", i))
		} else {
			is.Errorf("Task failed: %s", result.Error.Error())
		}
	}
}
