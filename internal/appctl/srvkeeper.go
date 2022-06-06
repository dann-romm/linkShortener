package appctl

import (
	"context"
	"sync/atomic"
	"time"
)

type Service interface {
	Init(ctx context.Context) error
	Ping(ctx context.Context) error
	Close() error
}

type ServiceKeeper struct {
	Services []Service

	PingPeriod      time.Duration
	PingTimeout     time.Duration
	ShutdownTimeout time.Duration

	state int32
	stop  chan struct{}
}

const (
	srvStateInit int32 = iota
	srvStateReady
	srvStateRunning
	srvStateShutdown
	srvStateOff
)

const (
	defaultPingPeriod      = time.Second * 15
	defaultPingTimeout     = time.Millisecond * 1500
	defaultShutdownTimeout = time.Millisecond * 15000
)

func (s *ServiceKeeper) Init(ctx context.Context) error {
	if !s.checkState(srvStateInit, srvStateReady) {
		return ErrWrongState
	}
	if err := s.initAllServices(ctx); err != nil {
		return err
	}
	s.stop = make(chan struct{})
	if s.PingPeriod == 0 {
		s.PingPeriod = defaultPingPeriod
	}
	if s.PingTimeout == 0 {
		s.PingTimeout = defaultPingTimeout
	}
	if s.ShutdownTimeout == 0 {
		s.ShutdownTimeout = defaultShutdownTimeout
	}
	return nil
}

func (s *ServiceKeeper) Watch(ctx context.Context) error {
	if !s.checkState(srvStateReady, srvStateRunning) {
		return ErrWrongState
	}
	if err := s.repeatPingServices(ctx); err != nil && err != ErrShutdown {
		return err
	}
	return nil
}

func (s *ServiceKeeper) Stop() {
	if s.checkState(srvStateRunning, srvStateShutdown) {
		close(s.stop)
	}
}

func (s *ServiceKeeper) Release() error {
	if !s.checkState(srvStateShutdown, srvStateOff) {
		return ErrWrongState
	}
	return s.release()
}

func (s *ServiceKeeper) initAllServices(ctx context.Context) error {
	initCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	p := ParallelRun{}
	for _, service := range s.Services {
		p.do(initCtx, service.Init)
	}
	return p.wait()
}

func (s *ServiceKeeper) pingServices(ctx context.Context) error {
	pingCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	p := ParallelRun{}
	for _, service := range s.Services {
		p.do(pingCtx, service.Ping)
	}
	return p.wait()
}

func (s *ServiceKeeper) repeatPingServices(ctx context.Context) error {
	for {
		select {
		case <-s.stop:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(s.PingPeriod):
			if err := s.pingServices(ctx); err != nil {
				return err
			}
		}
	}
}

func (s *ServiceKeeper) release() error {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.ShutdownTimeout)
	defer cancel()

	p := ParallelRun{}
	for _, service := range s.Services {
		p.do(shutdownCtx, func(context.Context) error {
			return service.Close()
		})
	}

	errWait := make(chan error)
	go func() {
		defer close(errWait)
		if err := p.wait(); err != nil {
			errWait <- err
		}
	}()

	for {
		select {
		case err, ok := <-errWait:
			if ok {
				return err
			}
			return nil
		case <-shutdownCtx.Done():
			return shutdownCtx.Err()
		}
	}
}

func (s *ServiceKeeper) checkState(old, new int32) bool {
	return atomic.CompareAndSwapInt32(&s.state, old, new)
}
