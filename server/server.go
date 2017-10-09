package server

import (
	"crypto/tls"

	"dev.sum7.eu/genofire/yaja/model"
)

type Server struct {
	TLSConfig  *tls.Config
	PortClient int
	PortServer int
	State      *model.State
}

func (srv *Server) Start() {

}

func (srv *Server) Close() {

}
