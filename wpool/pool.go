package wpool

import (
	"context"
	"sync"
	"sync/atomic"
)

var DefaultPool *defaultPool

const (
	defaultCap            = 10000
	defaultScaleThreshold = 1
)

func init() {
	DefaultPool = &defaultPool{
		cap:            defaultCap,
		scaleThreshold: defaultScaleThreshold,
	}
}

type WPool interface {
	SetCap(int32)
	Run(func())
	RunWithCtx(context.Context, func())
	SetPanicHandler(func(context.Context, any))
}

type defaultPool struct {
	cap int32

	taskHead *task
	taskTail *task
	mux      sync.Mutex

	// config
	scaleThreshold int32

	taskCount   atomic.Int32
	workerCount atomic.Int32

	panicHandler func(context.Context, any)
}

func (p *defaultPool) SetCap(i int32) {
	p.cap = i
}

func (p *defaultPool) Run(f func()) {
	p.RunWithCtx(context.Background(), f)
}

func (p *defaultPool) RunWithCtx(ctx context.Context, f func()) {
	t := taskPool.Get().(*task)
	t.f = f

	p.mux.Lock()
	if p.taskHead == nil {
		p.taskHead = t
		p.taskTail = t
	} else {
		p.taskTail.next = t
		p.taskTail = t
	}
	p.taskCount.Add(1)
	p.mux.Unlock()

	if (p.taskCount.Load() > p.scaleThreshold && p.workerCount.Load() < p.cap) || p.workerCount.Load() == 0 {
		w := workerPool.Get().(*worker)
		w.p = p
		w.Run()
		p.workerCount.Add(1)
	}
}

func (p *defaultPool) SetPanicHandler(handler func(context.Context, any)) {
	p.panicHandler = handler
}

func Run(f func()) {
	DefaultPool.Run(f)
}
