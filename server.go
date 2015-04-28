package fever

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/netutil"

	"github.com/lestrrat/go-server-starter-listener"
	"gopkg.in/tylerb/graceful.v1"
)

type server struct {
	*graceful.Server
}

func NewServer(srv *http.Server, timeout time.Duration) *server {
	return &server{&graceful.Server{
		Timeout: timeout,
		Server:  srv,
	}}
}

func Run(addr string, timeout time.Duration, n http.Handler) {
	srv := NewServer(&http.Server{Addr: addr, Handler: n}, timeout)
	if err := srv.ListenAndServe(); err != nil {
		if opErr, ok := err.(*net.OpError); !ok || (ok && opErr.Op != "accept") {
			logger := log.New(os.Stdout, "[graceful] ", 0)
			logger.Fatal(err)
		}
	}
}

func (srv *server) newListener() (net.Listener, error) {
	l, _ := ss.NewListener()
	if l == nil {
		addr := srv.Addr
		if addr == "" {
			addr = ":http"
		}
		var err error
		l, err = net.Listen("tcp", addr)
		if err != nil {
			return nil, err
		}
	}
	return l, nil
}

func (srv *server) ListenAndServe() error {
	l, err := srv.newListener()
	if err != nil {
		return err
	}
	if srv.ListenLimit != 0 {
		l = netutil.LimitListener(l, srv.ListenLimit)
	}
	return srv.Serve(l)
}

func (srv *server) ListenAndServeTLS(certFile, keyFile string) error {
	config := &tls.Config{}
	if srv.TLSConfig != nil {
		*config = *srv.TLSConfig
	}
	if config.NextProtos == nil {
		config.NextProtos = []string{"http/1.1"}
	}
	var err error
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}
	l, err := srv.newListener()
	if err != nil {
		return err
	}
	tlsListener := tls.NewListener(l, config)
	return srv.Serve(tlsListener)
}
