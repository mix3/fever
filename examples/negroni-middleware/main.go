package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/meatballhat/negroni-logrus"
	"github.com/mix3/fever"
	"github.com/mix3/fever/mux"
	"golang.org/x/net/context"
)

func main() {
	m := mux.New()
	m.Use(fever.WrapNegroniMiddleware(negronilogrus.NewMiddleware()))
	m.Get("/").ThenFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "success!")
	})
	fever.Run(":19300", 10*time.Second, m)
}
