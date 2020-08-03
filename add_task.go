package workqueue

import "fmt"

type AddTask struct {
	Number1 int
	Number2 int
}

func (t AddTask) Run() int {
	return t.Number1 + t.Number2
}
func (t AddTask) String() string {
	return fmt.Sprintf("AddTask(%d, %d)", t.Number1, t.Number2)
}
