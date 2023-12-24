package wpool

import (
	"context"
	"sync"
)

var taskPool = &sync.Pool{}

func init() {
	taskPool.New = NewTask
}

type task struct {
	ctx context.Context
	f   func()

	next *task
}

func NewTask() any {
	return &task{}
}

func (t *task) Recycle() {
	t.f = nil
	t.next = nil
	taskPool.Put(t)
}
