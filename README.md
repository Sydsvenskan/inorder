# In Order

Maintains ordering of asynchronous tasks.

## Usage

```golang
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
		Name:  name,
	}
	go func() {
		<-time.NewTimer(d).C
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
```
