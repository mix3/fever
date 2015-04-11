package chain_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mix3/fever"
	"github.com/mix3/fever/chain"
	"github.com/stretchr/testify/assert"

	"golang.org/x/net/context"
)

func TestChain(t *testing.T) {
	m1 := func(h fever.Handler) fever.Handler {
		return fever.HandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "m1s")
			h.ServeHTTP(c, w, r)
			fmt.Fprint(w, "m1e")
		})
	}
	m2 := func(h fever.Handler) fever.Handler {
		return fever.HandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "m2s")
			h.ServeHTTP(c, w, r)
			fmt.Fprint(w, "m2e")
		})
	}
	m3 := func(h fever.Handler) fever.Handler {
		return fever.HandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "m3s")
			h.ServeHTTP(c, w, r)
			fmt.Fprint(w, "m3e")
		})
	}
	m4 := func(h fever.Handler) fever.Handler {
		return fever.HandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "m4s")
			h.ServeHTTP(c, w, r)
			fmt.Fprint(w, "m4e")
		})
	}

	c1 := chain.New(m1, m2)
	c2 := c1.Append(m3, m4)
	c3 := c1.Append(m4, m3)
	{
		r, _ := http.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		h := c1.ThenFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "c1")
		})
		h.ServeHTTP(context.Background(), w, r)
		assert.Equal(t, "m1sm2sc1m2em1e", w.Body.String())
	}
	{
		r, _ := http.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		h := c2.ThenFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "c2")
		})
		h.ServeHTTP(context.Background(), w, r)
		assert.Equal(t, "m1sm2sm3sm4sc2m4em3em2em1e", w.Body.String())
	}
	{
		r, _ := http.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		h := c3.ThenFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "c3")
		})
		h.ServeHTTP(context.Background(), w, r)
		assert.Equal(t, "m1sm2sm4sm3sc3m3em4em2em1e", w.Body.String())
	}
}
