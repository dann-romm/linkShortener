package appctl

import (
	"context"
	"fmt"
	"sync"
)

type ParallelRun struct {
	mux sync.Mutex
	wg  sync.WaitGroup
	err arrError
}

func (p *ParallelRun) do(ctx context.Context, f func(context.Context) error) {
	p.wg.Add(1)
	go func() {
		defer func() {
			r := recover()
			if r != nil {
				p.mux.Lock()
				p.err = append(p.err, fmt.Errorf("unhandled error: %v", r))
				p.mux.Unlock()
			}
			p.wg.Done()
		}()
		if err := f(ctx); err != nil {
			p.mux.Lock()
			p.err = append(p.err, fmt.Errorf("%w", err))
			p.mux.Unlock()
		}
	}()
}

func (p *ParallelRun) wait() error {
	p.wg.Wait()
	if len(p.err) > 0 {
		return p.err
	}
	return nil
}
