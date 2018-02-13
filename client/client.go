package client

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"net"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"dev.sum7.eu/genofire/yaja/messages"
)

// Client holds XMPP connection opitons
type Client struct {
	Protocol string // tcp tcp4 tcp6 are supported
	Timeout  time.Duration
	conn     net.Conn // connection to server
	out      *xml.Encoder
	in       *xml.Decoder

	Logging *log.Logger

	JID *messages.JID

	reply map[string]chan *messages.IQClient

	skipError bool
	iq        chan *messages.IQClient
	presence  chan *messages.PresenceClient
	mesage    chan *messages.MessageClient
}

func NewClient(jid *messages.JID, password string) (*Client, error) {
	client := &Client{
		Protocol: "tcp",
		JID:      jid,
		Logging:  log.New(),
	}
	return client, client.Connect(password)

}
func (client *Client) Connect(password string) error {
	_, srvEntries, err := net.LookupSRV("xmpp-client", "tcp", client.JID.Domain)
	addr := client.JID.Domain + ":5222"
	if err == nil && len(srvEntries) > 0 {
		bestSrv := srvEntries[0]
		for _, srv := range srvEntries {
			if srv.Priority <= bestSrv.Priority && srv.Weight >= bestSrv.Weight {
				bestSrv = srv
				addr = fmt.Sprintf("%s:%d", srv.Target, srv.Port)
			}
		}
	}
	a := strings.SplitN(addr, ":", 2)
	if len(a) == 1 {
		addr += ":5222"
	}
	if client.Protocol == "" {
		client.Protocol = "tcp"
	}
	conn, err := net.DialTimeout(client.Protocol, addr, client.Timeout)
	client.setConnection(conn)
	if err != nil {
		return err
	}

	if err = client.connect(password); err != nil {
		client.Close()
		return err
	}
	return nil
}

// Close closes the XMPP connection
func (c *Client) Close() error {
	if c.conn != (*tls.Conn)(nil) {
		return c.conn.Close()
	}
	return nil
}
