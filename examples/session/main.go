package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/mix3/fever"
	"github.com/mix3/fever-sessions"
	"github.com/mix3/fever/mux"
	"github.com/soh335/go-test-redisserver"
	"golang.org/x/net/context"
)

func init() {
	gob.Register(struct{}{})
}

func main() {
	rt, err := redistest.NewServer(true, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer rt.Stop()
	store, err := sessions.NewRedisStore("unix", rt.Config["unixsocket"], "")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = store.Close()
	}()

	ss := sessions.New(store)
	m := mux.New()
	m.Use(ss.Middleware)
	m.Get("/").ThenFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
		s := sessions.Session(c)
		if s.Exists("username") {
			fmt.Fprintf(w, "hello %s", s.Get("username").(string))
		} else {
			fmt.Fprintf(w, "OK")
		}
	})
	m.Get("/counter").ThenFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
		s := sessions.Session(c)
		v := 1
		if s.Exists("counter") {
			v = s.Get("counter").(int) + 1
		}
		s.Set("counter", v)
		fmt.Fprintf(w, "counter => %d", v)
	})
	m.Get("/login").ThenFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
		s := sessions.Session(c)
		s.Set("username", "foo")
		fmt.Fprintf(w, "login")
	})
	m.Get("/logout").ThenFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
		s := sessions.Session(c)
		if s.Exists("username") {
			s.Expire(true)
		}
		fmt.Fprintf(w, "logout")
	})
	m.Get("/switch").ThenFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
		s := sessions.Session(c)
		flashes := s.Flashes()
		if 0 < len(flashes) {
			fmt.Fprintf(w, "ON")
		} else {
			s.AddFlash(struct{}{})
			fmt.Fprintf(w, "OFF")
		}
	})
	fever.Run(":19300", 10*time.Second, m)
}
