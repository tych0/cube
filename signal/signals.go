package cube

import (
	"os"
	"os/signal"

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

type sh struct {
	signalCh   chan os.Signal
	signals    map[os.Signal]SignalHandler
	ignSignals map[os.Signal]struct{}
	ctx        context.Context
	cancelFunc context.CancelFunc
	running    bool
}

// NewSignalRouter returns a signal router.
func NewSignalRouter() SignalRouter {
	ctx, cancelFunc := context.WithCancel(context.Background())
	return &sh{
		signalCh:   make(chan os.Signal),
		signals:    make(map[os.Signal]SignalHandler),
		ignSignals: make(map[os.Signal]struct{}),
		ctx:        ctx,
		cancelFunc: cancelFunc,
		running:    false,
	}
}

func (s *sh) Handle(sig os.Signal, h SignalHandler) {
	s.signals[sig] = h
	signal.Notify(s.signalCh, sig)
	delete(s.ignSignals, sig)
}

func (s *sh) Reset(sig os.Signal) {
	delete(s.signals, sig)
	signal.Reset(sig)
}

func (s *sh) Ignore(sig os.Signal) {
	delete(s.signals, sig)
	signal.Ignore(sig)
	s.ignSignals[sig] = struct{}{}
}

func (s *sh) IsHandled(sig os.Signal) bool {
	_, ok := s.signals[sig]
	return ok
}

func (s *sh) IsIgnored(sig os.Signal) bool {
	_, ok := s.ignSignals[sig]
	return ok
}

// Implement service lifecycle receivers

func (s *sh) OnStart() error {
	go func() {
		defer func() {
			s.running = false
		}()
		// This go routine dies with the server
		for {
			select {
			case <-s.ctx.Done():
				// We are done exit.
				return
			case sig := <-s.signalCh:
				if h, ok := s.signals[sig]; ok {
					h(sig)
				}
			}
		}
	}()
	s.running = true
	return nil
}

func (s *sh) OnStop() error {
	s.cancelFunc()
	return nil
}

func (s *sh) OnConfigure(cfg interface{}) error {
	return nil
}

func (s *sh) IsHealthy() bool {
	return s.running
}
