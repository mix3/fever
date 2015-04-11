package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/mix3/fever"
	"github.com/mix3/fever/mux"
	"golang.org/x/net/context"
)

func main() {
	m := mux.New()
	m.Get("/").ThenFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})
	m.Get("/:name").ThenFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
		p := mux.Params(c)
		fmt.Fprintf(w, "Hello %s", p.ByName("name"))
	})
	srv := fever.NewServer(&http.Server{Addr: ":19300", Handler: m}, 10*time.Second)
	srv.ListenAndServeTLS("examples/tls/ca.crt", "examples/tls/ca.key")
}
