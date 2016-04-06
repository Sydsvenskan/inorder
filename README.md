# In Order

Maintains ordering of asynchronous tasks.

## Usage

```golang
type DummyTask struct {
	block        *inorder.Block
	ExpectedName string
	Name         string
}

func (dt *DummyTask) Wait() chan bool {
	return dt.block.Wait()
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
```
