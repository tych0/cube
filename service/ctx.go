package service

import (
	"context"
)

// Context provides a wrapper interface for go context.
//
// Ctx() returns the underlying go context.
//
// Shutdown() cancels the go context.
type Context interface {
	Ctx() context.Context
	Shutdown()
}

// NewContext creates a new service context.
func NewContext() Context {
	ctx, cancelFunc := context.WithCancel(context.Background())
	return &srvCtx{
		ctx:        ctx,
		cancelFunc: cancelFunc,
	}
}

func (sc *srvCtx) Ctx() context.Context {
	return sc.ctx
}

func (sc *srvCtx) Shutdown() {
	sc.cancelFunc()
}

type srvCtx struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
}
