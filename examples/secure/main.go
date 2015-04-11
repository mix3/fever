package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/mix3/fever"
	"github.com/mix3/fever/mux"
	"github.com/unrolled/secure"
	"golang.org/x/net/context"
)

func main() {
	s := secure.New(secure.Options{
		FrameDeny: true,
	})
	f := func(h fever.Handler) fever.Handler {
		return fever.HandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
			err := s.Process(w, r)
			if err != nil {
				return
			}
			h.ServeHTTP(c, w, r)
		})
	}
	m := mux.New()
	m.Use(f)
	m.Get("/").ThenFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})
	fever.Run(":19300", 10*time.Second, m)
}
