package eebus

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/andig/evcc/hems/eebus/ship"
	"github.com/gorilla/websocket"
)

func NewServer(addr string, cert tls.Certificate) (*http.Server, error) {
	s := &http.Server{
		Addr:    addr,
		Handler: &Handler{},
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientAuth:   tls.NoClientCert,
			CipherSuites: ship.CipherSuites,
		},
	}

	go func() {
		// if err := s.ListenAndServeTLS("", ""); err != nil {
		// 	log.Fatal(err)
		// }

		if err := serve(s); err != nil {
			log.Println(err)
		}
	}()

	return s, nil
}

func serve(srv *http.Server) error {
	ln, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		return err
	}

	defer ln.Close()

	tlsListener := tls.NewListener(ln, srv.TLSConfig)
	return srv.Serve(tlsListener)
}

type Handler struct{}

func (s *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log := log.New(&writer{os.Stdout, "2006/01/02 15:04:05 "}, "[server] ", 0)
	log.Printf("request: %v", r)

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
		Subprotocols:    []string{ship.SubProtocol},
	}

	// upgrade
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	defer conn.Close()

	// return and close connection
	if conn.Subprotocol() != ship.SubProtocol {
		log.Println("protocol mismatch:", conn.Subprotocol())
		return
	}

	// ship
	sc := &ship.Server{
		Log: log,
	}

	if err := sc.Serve(conn); err != nil {
		log.Println(err)
	}

	log.Println("done")
}
