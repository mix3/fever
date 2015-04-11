package mux_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mix3/fever"
	"github.com/mix3/fever/mux"
	"github.com/stretchr/testify/assert"

	"golang.org/x/net/context"
)

func TestMux(t *testing.T) {
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

	m := mux.New()
	m.Use(m1, m2)
	m.Get("/hoge").Use(m3, m4).ThenFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hoge")
	})
	m.Get("/fuga").Use(m4, m3).ThenFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "fuga")
	})
	{
		r, _ := http.NewRequest("GET", "/hoge", nil)
		w := httptest.NewRecorder()
		m.ServeHTTP(w, r)
		assert.Equal(t, "m1sm2sm3sm4shogem4em3em2em1e", w.Body.String())
	}
	{
		r, _ := http.NewRequest("GET", "/fuga", nil)
		w := httptest.NewRecorder()
		m.ServeHTTP(w, r)
		assert.Equal(t, "m1sm2sm4sm3sfugam3em4em2em1e", w.Body.String())
	}
}

func TestParams(t *testing.T) {
	m := mux.New()
	m.Get("/show/:id").ThenFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
		p := mux.Params(c)
		fmt.Fprint(w, p.ByName("id"))
	})
	r, _ := http.NewRequest("GET", "/show/hoge", nil)
	w := httptest.NewRecorder()
	m.ServeHTTP(w, r)
	assert.Equal(t, "hoge", w.Body.String())
}

func TestNotFound(t *testing.T) {
	{ // Default Handler
		m := mux.New()
		m.Use(func(h fever.Handler) fever.Handler {
			return fever.HandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Hoge", "Fuga")
				h.ServeHTTP(c, w, r)
			})
		})
		r, _ := http.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		m.ServeHTTP(w, r)
		assert.Equal(t, 404, w.Code)
		assert.Equal(t, "404 page not found\n", w.Body.String())
		assert.Equal(t, "Fuga", w.HeaderMap["Hoge"][0])
	}
	{ // Custom Handler
		m := mux.New()
		m.Use(func(h fever.Handler) fever.Handler {
			return fever.HandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Hoge", "Fuga")
				h.ServeHTTP(c, w, r)
			})
		})
		m.NotFound = fever.HandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
			http.Error(w, "NotFound", http.StatusNotFound)
		})
		r, _ := http.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		m.ServeHTTP(w, r)
		assert.Equal(t, 404, w.Code)
		assert.Equal(t, "NotFound\n", w.Body.String())
		assert.Equal(t, "Fuga", w.HeaderMap["Hoge"][0])
	}
}
