package eebus

import (
	"crypto/tls"
	"log"
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
		if err := s.ListenAndServeTLS("", ""); err != nil {
			log.Fatal(err)
		}
	}()

	return s, nil
}

type Handler struct{}

func (s *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
		Subprotocols:    []string{ship.SubProtocol},
	}

	// upgrade
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	// return and close connection
	if conn.Subprotocol() != ship.SubProtocol {
		return
	}

	// ship
	sc := ship.New(conn)
	sc.Log = log.New(&writer{os.Stdout, "2006/01/02 15:04:05 "}, "[server] ", 0)

	if err := sc.Serve(); err != nil {
		sc.Log.Println("connect:", err)
		return
	}
}
