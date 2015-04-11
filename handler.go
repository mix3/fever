package fever

import (
	"net/http"

	"github.com/codegangsta/negroni"

	"golang.org/x/net/context"
)

type Handler interface {
	ServeHTTP(context.Context, http.ResponseWriter, *http.Request)
}

type HandlerFunc func(c context.Context, w http.ResponseWriter, r *http.Request)

func (h HandlerFunc) ServeHTTP(c context.Context, w http.ResponseWriter, r *http.Request) {
	h(c, w, r)
}

func Wrap(h http.Handler) Handler {
	return HandlerFunc(func(_ context.Context, w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}

type Middleware func(Handler) Handler

func WrapNegroniMiddleware(nh negroni.Handler) Middleware {
	return Middleware(func(fh Handler) Handler {
		return HandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
			if _, ok := w.(negroni.ResponseWriter); !ok {
				w = negroni.NewResponseWriter(w)
			}
			nh.ServeHTTP(w, r, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fh.ServeHTTP(c, w, r)
			}))
		})
	})
}
