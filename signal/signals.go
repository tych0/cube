package cube

import (
	"os"
	"os/signal"
	"sync"

	"context"
)

// SignalHandler is a function that handles the signal.
type SignalHandler func(os.Signal)

// SignalRouter routes signals to registered handler.
type SignalRouter interface {
	// Handle registers a signal handler.
	Handle(sig os.Signal, h SignalHandler)

	// Reset resets a signal handler.
	Reset(sig os.Signal)

	// Ignore ignores a signal.
	Ignore(sig os.Signal)

	// IsHandled checks if a signal being routed to a handler.
	IsHandled(sig os.Signal) bool

	// IsIgnored checks if a signal is being ignored.
	IsIgnored(sig os.Signal) bool
}

type router struct {
	signalCh   chan os.Signal
	signals    map[os.Signal]SignalHandler
	ignSignals map[os.Signal]struct{}
	ctx        context.Context
	cancelFunc context.CancelFunc
	running    bool
	lock       *sync.RWMutex
}

// NewSignalRouter returns a signal router.
func NewSignalRouter() SignalRouter {
	ctx, cancelFunc := context.WithCancel(context.Background())
	return &router{
		signalCh:   make(chan os.Signal),
		signals:    make(map[os.Signal]SignalHandler),
		ignSignals: make(map[os.Signal]struct{}),
		ctx:        ctx,
		cancelFunc: cancelFunc,
		running:    false,
		lock:       &sync.RWMutex{},
	}
}

func (s *router) Handle(sig os.Signal, h SignalHandler) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.signals[sig] = h
	signal.Notify(s.signalCh, sig)
	delete(s.ignSignals, sig)
}

func (s *router) Reset(sig os.Signal) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.signals, sig)
	signal.Reset(sig)
}

func (s *router) Ignore(sig os.Signal) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.signals, sig)
	signal.Ignore(sig)
	s.ignSignals[sig] = struct{}{}
}

func (s *router) IsHandled(sig os.Signal) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	_, ok := s.signals[sig]
	return ok
}

func (s *router) IsIgnored(sig os.Signal) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	_, ok := s.ignSignals[sig]
	return ok
}

// Implement service lifecycle receivers

func (s *router) OnStart() error {
	go func() {
		defer func() {
			s.lock.Lock()
			defer s.lock.Unlock()
			s.running = false
		}()
		// This go routine dies with the server
		for {
			select {
			case <-s.ctx.Done():
				// We are done exit.
				return
			case sig := <-s.signalCh:
				func() {
					s.lock.RLock()
					defer s.lock.RUnlock()
					if h, ok := s.signals[sig]; ok {
						h(sig)
					}
				}()
			}
		}
	}()
	s.lock.Lock()
	s.running = true
	s.lock.Unlock()
	return nil
}

func (s *router) OnStop() error {
	s.cancelFunc()
	return nil
}

func (s *router) OnConfigure(cfg interface{}) error {
	return nil
}

func (s *router) IsHealthy() bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.running
}
