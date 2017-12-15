package server

import (
	"crypto/tls"
	"net"

	"github.com/genofire/yaja/database"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
)

type Server struct {
	TLSConfig       *tls.Config
	TLSManager      *autocert.Manager
	ClientAddr      []string
	ServerAddr      []string
	Database        *database.State
	LoggingClient   log.Level
	RegisterEnable  bool     `toml:"enable"`
	RegisterDomains []string `toml:"domains"`
}

func (srv *Server) Start() {
	for _, addr := range srv.ServerAddr {
		socket, err := net.Listen("tcp", addr)
		if err != nil {
			log.Warn("create server socket: ", err.Error())
			break
		}
		go srv.listenServer(socket)
	}

	for _, addr := range srv.ClientAddr {
		socket, err := net.Listen("tcp", addr)
		if err != nil {
			log.Warn("create client socket: ", err.Error())
			break
		}
		go srv.listenClient(socket)
	}
}

func (srv *Server) listenServer(s2s net.Listener) {
	for {
		conn, err := s2s.Accept()
		if err != nil {
			log.Warn("accepting server connection: ", err.Error())
			break
		}
		go srv.handleServer(conn)
	}
}

func (srv *Server) listenClient(c2s net.Listener) {
	for {
		conn, err := c2s.Accept()
		if err != nil {
			log.Warn("accepting client connection: ", err.Error())
			break
		}
		go srv.handleClient(conn)
	}
}

func (srv *Server) handleServer(conn net.Conn) {
	log.Info("new server connection:", conn.RemoteAddr())
}

func (srv *Server) handleClient(conn net.Conn) {
	log.Info("new client connection:", conn.RemoteAddr())
	client := NewClient(conn, srv)
	state := ConnectionStartup()

	for {
		state, client = state.Process(client)
		if state == nil {
			client.log.Info("disconnect")
			client.Close()
			//s.DisconnectBus <- Disconnect{Jid: client.jid}
			return
		}
		// run next state
	}
}

func (srv *Server) Close() {

}
