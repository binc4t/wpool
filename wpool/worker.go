package wpool

import (
	"fmt"
	"runtime/debug"
	"sync"
)

var workerPool = &sync.Pool{}

func init() {
	workerPool.New = NewWorker
}

type worker struct {
	p *defaultPool
}

func NewWorker() any {
	return &worker{}
}

func (w *worker) Run() {
	go func() {
		for {
			var t *task
			w.p.mux.Lock()
			if w.p.taskHead != nil {
				t = w.p.taskHead
				w.p.taskHead = w.p.taskHead.next
			}
			w.p.mux.Unlock()

			if t == nil {
				w.p.workerCount.Add(-1)
				w.Recycle()
				return
			}

			w.p.taskCount.Add(-1)

			func() {
				defer func() {
					if r := recover(); r != nil {
						if w.p.panicHandler != nil {
							w.p.panicHandler(t.ctx, r)
						} else {
							fmt.Printf("panic in worker, ctx: %v, err: %v, stack %v\n", t.ctx.Value("name"), r, debug.Stack())
						}
					}
				}()
				t.f()
			}()
			t.Recycle()
		}
	}()
}

func (w *worker) Recycle() {
	w.p = nil
	workerPool.Put(w)
}
