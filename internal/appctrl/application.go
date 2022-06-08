package appctrl

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type Resources interface {
	Init(context.Context) error
	Watch(context.Context) error
	Stop()
	Release() error
}

type Application struct {
	MainFunc func(ctx context.Context, halt <-chan struct{}) error

	Resources Resources

	TerminationTimeout    time.Duration
	InitializationTimeout time.Duration

	appState int32
	halt     chan struct{}
	done     chan struct{}
	errMux   sync.Mutex
	err      error
}

const (
	appStateInit int32 = iota
	appStateRunning
	appStateHalt
	appStateShutdown
)

const (
	defaultTerminationTimeout    = time.Second
	defaultInitializationTimeout = time.Second * 15
)

func (a *Application) Run() error {
	if a.MainFunc == nil {
		return ErrMainOmitted
	}
	if !a.checkState(appStateInit, appStateRunning) {
		return ErrWrongState
	}
	if err := a.init(); err != nil {
		a.err = err
		a.appState = appStateShutdown
		return err
	}

	services := make(chan struct{})
	if a.Resources != nil {
		go a.watchResources(services)
	}

	osSig := make(chan os.Signal, 1)
	signal.Notify(osSig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	a.setError(a.run(osSig))

	if a.Resources != nil {
		a.Resources.Stop()
		select {
		case <-services:
		case <-time.After(a.TerminationTimeout):
		}
		a.setError(a.Resources.Release())
	}
	return a.getError()
}

func (a *Application) init() error {
	log.Println("[appctrl] Initializing application")
	if a.TerminationTimeout == 0 {
		a.TerminationTimeout = defaultTerminationTimeout
	}
	if a.InitializationTimeout == 0 {
		a.InitializationTimeout = defaultInitializationTimeout
	}
	a.halt = make(chan struct{})
	a.done = make(chan struct{})
	if a.Resources != nil {
		ctx, cancel := context.WithTimeout(a, a.InitializationTimeout)
		defer cancel()
		return a.Resources.Init(ctx)
	}
	return nil
}

func (a *Application) watchResources(services chan<- struct{}) {
	log.Println("[appctrl] Starting resource watcher")
	defer close(services)
	defer a.Shutdown()
	a.setError(a.Resources.Watch(a))
}

func (a *Application) run(osSig <-chan os.Signal) error {
	log.Println("[appctrl] Running application")
	defer a.Shutdown()
	errRun := make(chan error, 1)
	errHalt := make(chan error, 1)

	go func() {
		defer close(errRun)
		if err := a.MainFunc(a, a.halt); err != nil {
			errRun <- err
		}
	}()

	go func() {
		defer close(errHalt)
		select {
		case <-osSig:
			a.Halt()
			select {
			case <-time.After(a.TerminationTimeout):
				errHalt <- ErrTermTimeout
			case <-a.done:
			}
		case <-a.done:
		}
	}()

	select {
	case err, ok := <-errRun:
		if ok && err != nil {
			return err
		}
	case err, ok := <-errHalt:
		if ok && err != nil {
			return err
		}
	case <-a.done:
	}
	return nil
}

func (a *Application) Halt() {
	if a.checkState(appStateRunning, appStateHalt) {
		log.Println("[appctrl] Halting application")
		close(a.halt)
	}
}

func (a *Application) Shutdown() {
	a.Halt()
	if a.checkState(appStateHalt, appStateShutdown) {
		log.Println("[appctrl] Shutting down application")
		close(a.done)
	}
}

func (a *Application) checkState(old, new int32) bool {
	return atomic.CompareAndSwapInt32(&a.appState, old, new)
}

func (a *Application) setError(err error) {
	if err == nil {
		return
	}
	a.errMux.Lock()
	if a.err == nil {
		a.err = err
	}
	a.errMux.Unlock()
}

func (a *Application) getError() error {
	a.errMux.Lock()
	err := a.err
	a.errMux.Unlock()
	return err
}

func (a *Application) Deadline() (deadline time.Time, ok bool) {
	return time.Time{}, false
}

func (a *Application) Done() <-chan struct{} {
	return a.done
}

func (a *Application) Err() error {
	if err := a.getError(); err != nil {
		return err
	}
	if atomic.LoadInt32(&a.appState) == appStateShutdown {
		return ErrShutdown
	}
	return nil
}

func (a *Application) Value(any) any {
	return a
}
