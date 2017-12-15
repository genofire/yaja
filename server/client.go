package server

import (
	"encoding/xml"
	"net"

	"github.com/genofire/yaja/model"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	log *log.Entry

	Conn net.Conn
	out  *xml.Encoder
	in   *xml.Decoder

	Server  *Server
	jid     *model.JID
	account *model.Account

	messages chan interface{}
	close    chan interface{}
}

func NewClient(conn net.Conn, srv *Server) *Client {
	logger := log.New()
	logger.SetLevel(srv.LoggingClient)
	client := &Client{
		Conn:   conn,
		Server: srv,
		log:    log.NewEntry(logger),
		in:     xml.NewDecoder(conn),
		out:    xml.NewEncoder(conn),
	}
	return client
}

func (client *Client) NewConnecting(conn net.Conn) {
	client.Conn = conn
	client.in = xml.NewDecoder(conn)
	client.out = xml.NewEncoder(conn)
}

func (client *Client) Read() (*xml.StartElement, error) {
	for {
		nextToken, err := client.in.Token()
		if err != nil {
			return nil, err
		}
		switch nextToken.(type) {
		case xml.StartElement:
			element := nextToken.(xml.StartElement)
			return &element, nil
		}
	}
}

func (client *Client) DomainRegisterAllowed() bool {
	if client.jid.Domain == "" {
		return false
	}

	for _, domain := range client.Server.RegisterDomains {
		if domain == client.jid.Domain {

			return !client.Server.RegisterEnable
		}
	}
	return client.Server.RegisterEnable
}

func (client *Client) Close() {
	client.close <- true
	client.Conn.Close()
}
