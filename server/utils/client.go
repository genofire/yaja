package utils

import (
	"encoding/xml"
	"net"

	log "github.com/sirupsen/logrus"

	"dev.sum7.eu/genofire/yaja/model"
	"dev.sum7.eu/genofire/yaja/xmpp/base"
)

type Client struct {
	Log *log.Entry

	Conn net.Conn
	Out  *xml.Encoder
	In   *xml.Decoder

	JID     *xmppbase.JID
	account *model.Account

	Messages chan interface{}
	close    chan interface{}
}

func NewClient(conn net.Conn, level log.Level) *Client {
	logger := log.New()
	logger.SetLevel(level)
	client := &Client{
		Conn:     conn,
		Log:      log.NewEntry(logger),
		In:       xml.NewDecoder(conn),
		Out:      xml.NewEncoder(conn),
		Messages: make(chan interface{}),
		close:    make(chan interface{}),
	}
	return client
}

func (client *Client) SetConnecting(conn net.Conn) {
	client.Conn = conn
	client.In = xml.NewDecoder(conn)
	client.Out = xml.NewEncoder(conn)
}

func (client *Client) Read() (*xml.StartElement, error) {
	for {
		nextToken, err := client.In.Token()
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

func (client *Client) OnClose() <-chan interface{} {
	return client.close
}

func (client *Client) Close() {
	client.close <- true
	client.Conn.Close()
}
