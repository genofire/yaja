package client

import (
	"encoding/xml"
	"fmt"
	"net"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"dev.sum7.eu/genofire/yaja/xmpp"
	"dev.sum7.eu/genofire/yaja/xmpp/base"
)

// Client holds XMPP connection opitons
type Client struct {
	Protocol string // tcp tcp4 tcp6 are supported
	Timeout  time.Duration
	conn     net.Conn // connection to server
	out      *xml.Encoder
	in       *xml.Decoder

	Logging *log.Entry

	JID *xmppbase.JID

	SkipError bool
	msg       chan interface{}
	reply     map[string]chan *xmpp.IQClient
}

func NewClient(jid *xmppbase.JID, password string) (*Client, error) {
	client := &Client{
		JID:     jid,
		Logging: log.New().WithField("jid", jid.String()),
	}
	return client, client.Connect(password)

}
func (client *Client) Connect(password string) error {
	_, srvEntries, err := net.LookupSRV("xmpp-client", "tcp", client.JID.Domain)
	addr := client.JID.Domain
	if err == nil && len(srvEntries) > 0 {
		bestSrv := srvEntries[0]
		for _, srv := range srvEntries {
			if srv.Priority <= bestSrv.Priority && srv.Weight >= bestSrv.Weight {
				bestSrv = srv
				addr = fmt.Sprintf("%s:%d", srv.Target, srv.Port)
			}
		}
	}
	if strings.LastIndex(addr, ":") <= strings.LastIndex(addr, "]") {
		addr += ":5222"
	}
	if client.Protocol == "" {
		client.Protocol = "tcp"
	}
	client.Logging.Debug("try tcp connect")
	conn, err := net.DialTimeout(client.Protocol, addr, client.Timeout)
	if err != nil {
		return err
	}
	client.Logging.Debug("tcp connected")
	client.setConnection(conn)

	if err = client.connect(password); err != nil {
		client.Close()
		return err
	}
	return nil
}

// Close closes the XMPP connection
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
