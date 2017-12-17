package server

import (
	"crypto/tls"
	"net"

	"github.com/genofire/yaja/database"
	"github.com/genofire/yaja/model"
	"github.com/genofire/yaja/server/extension"
	"github.com/genofire/yaja/server/toclient"
	"github.com/genofire/yaja/server/toserver"
	"github.com/genofire/yaja/server/utils"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"
)

type Server struct {
	TLSConfig        *tls.Config
	TLSManager       *autocert.Manager
	ClientAddr       []string
	ServerAddr       []string
	Database         *database.State
	LoggingClient    log.Level
	LoggingServer    log.Level
	RegisterEnable   bool
	RegisterDomains  []string
	ExtensionsClient extension.Extensions
	ExtensionsServer extension.Extensions
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

	client := utils.NewClient(conn, srv.LoggingClient)
	client.Log = client.Log.WithField("c", "s2s")

	state := toserver.ConnectionStartup(srv.Database, srv.TLSConfig, srv.TLSManager, srv.ExtensionsServer, client)

	for {
		state = state.Process()
		if state == nil {
			client.Log.Info("disconnect")
			client.Close()
			return
		}
		// run next state
	}
}

func (srv *Server) handleClient(conn net.Conn) {
	log.Info("new client connection:", conn.RemoteAddr())

	client := utils.NewClient(conn, srv.LoggingServer)
	client.Log = client.Log.WithField("c", "c2s")

	state := toclient.ConnectionStartup(srv.Database, srv.TLSConfig, srv.TLSManager, srv.DomainRegisterAllowed, srv.ExtensionsClient, client)

	for {
		state = state.Process()
		if state == nil {
			client.Log.Info("disconnect")
			client.Close()
			//s.DisconnectBus <- Disconnect{Jid: client.jid}
			return
		}
		// run next state
	}
}

func (srv *Server) DomainRegisterAllowed(jid *model.JID) bool {
	if jid.Domain == "" {
		return false
	}

	for _, domain := range srv.RegisterDomains {
		if domain == jid.Domain {
			return !srv.RegisterEnable
		}
	}
	return srv.RegisterEnable
}

func (srv *Server) Close() {

}
