// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
package mux

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mix3/fever"
	"github.com/mix3/fever/chain"
	"golang.org/x/net/context"
)

type Mux struct {
	Router   *httprouter.Router
	Chain    chain.Chain
	NotFound fever.Handler
}

func New() *Mux {
	return &Mux{Router: httprouter.New()}
}

func (m *Mux) Use(mws ...fever.Middleware) {
	m.Chain = m.Chain.Append(mws...)
}

func (m *Mux) Handle(method, path string) *router {
	return &router{mux: m, chain: m.Chain, method: method, path: path}
}

func (m *Mux) Get(path string) *router {
	return m.Handle("GET", path)
}

func (m *Mux) Head(path string) *router {
	return m.Handle("HEAD", path)
}

func (m *Mux) Options(path string) *router {
	return m.Handle("OPTIONS", path)
}

func (m *Mux) Post(path string) *router {
	return m.Handle("POST", path)
}

func (m *Mux) Put(path string) *router {
	return m.Handle("PUT", path)
}

func (m *Mux) Patch(path string) *router {
	return m.Handle("PATCH", path)
}

func (m *Mux) Delete(path string) *router {
	return m.Handle("DELETE", path)
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if m.NotFound == nil {
		m.NotFound = fever.Wrap(http.HandlerFunc(http.NotFound))
	}
	m.Router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.Chain.Then(m.NotFound).ServeHTTP(context.Background(), w, r)
	})
	m.Router.ServeHTTP(w, r)
}

type router struct {
	mux    *Mux
	chain  chain.Chain
	method string
	path   string
}

func (r *router) Use(mws ...fever.Middleware) *router {
	r.chain = r.chain.Append(mws...)
	return r
}

func (r *router) Then(h fever.Handler) {
	r.mux.Router.Handle(r.method, r.path, wrap(r.chain.Then(h)))
}

func (r *router) ThenFunc(f fever.HandlerFunc) {
	r.mux.Router.Handle(r.method, r.path, wrap(r.chain.ThenFunc(f)))
}

var paramsKey = struct{}{}

func wrap(h fever.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		c := context.Background()
		c = context.WithValue(c, paramsKey, p)
		h.ServeHTTP(c, w, r)
	}
}

func Params(c context.Context) httprouter.Params {
	return c.Value(paramsKey).(httprouter.Params)
}
