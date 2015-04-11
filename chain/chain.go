package chain

import (
	"net/http"

	"github.com/mix3/fever"
)

type Chain struct {
	mws []fever.Middleware
}

func New(mws ...fever.Middleware) Chain {
	c := Chain{}
	c.mws = append(c.mws, mws...)
	return c
}

func (c Chain) Then(h fever.Handler) fever.Handler {
	var final fever.Handler
	if h != nil {
		final = h
	} else {
		final = fever.Wrap(http.DefaultServeMux)
	}
	for i := len(c.mws) - 1; i >= 0; i-- {
		fn := c.mws[i]
		final = fn(final)
	}
	return final
}

func (c Chain) ThenFunc(fn fever.HandlerFunc) fever.Handler {
	if fn == nil {
		return c.Then(nil)
	}
	return c.Then(fever.HandlerFunc(fn))
}

func (c Chain) Append(mws ...fever.Middleware) Chain {
	newMws := make([]fever.Middleware, len(c.mws)+len(mws))
	copy(newMws, c.mws)
	copy(newMws[len(c.mws):], mws)
	newChain := New(newMws...)
	return newChain
}
