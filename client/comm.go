package client

import (
	"encoding/xml"

	"dev.sum7.eu/genofire/yaja/xmpp"
)

func read(decoder *xml.Decoder) (*xml.StartElement, error) {
	for {
		nextToken, err := decoder.Token()
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
func (client *Client) Read() (*xml.StartElement, error) {
	return read(client.in)
}
func (client *Client) Decode(p interface{}, element *xml.StartElement) error {
	err := client.in.DecodeElement(p, element)
	if err != nil {
		return err
	} else {
		if b, err := xml.Marshal(p); err == nil {
			client.Logging.Debugf("decode %v", string(b))
		} else {
			client.Logging.Debugf("decode %v", p)
		}
	}
	return nil
}
func (client *Client) ReadDecode(p interface{}) error {
	element, err := client.Read()
	if err != nil {
		return err
	}
	var iq *xmpp.IQClient
	iq, ok := p.(*xmpp.IQClient)
	if !ok {
		iq = &xmpp.IQClient{}
	}
	err = client.Decode(iq, element)
	if err == nil && iq.Ping != nil && iq.Type == xmpp.IQTypeGet {
		client.Logging.Info("client.ReadElement: auto answer ping")
		iq.Type = xmpp.IQTypeResult
		iq.To = iq.From
		iq.From = client.JID
		client.Send(iq)
		return nil
	}
	if ok {
		return err
	}
	return client.Decode(p, element)
}
func (client *Client) send(p interface{}) error {
	b, err := xml.Marshal(p)
	if err != nil {
		client.Logging.Warnf("error send %v", p)
		return err
	}
	client.Logging.Debugf("send %v", string(b))
	_, err = client.conn.Write(b)
	return err
}

func (client *Client) Send(p interface{}) error {
	msg, ok := p.(*xmpp.MessageClient)
	if ok {
		msg.From = client.JID
		return client.send(msg)
	}
	iq, ok := p.(*xmpp.IQClient)
	if ok {
		iq.From = client.JID
		return client.send(iq)
	}
	pc, ok := p.(*xmpp.PresenceClient)
	if ok {
		pc.From = client.JID
		return client.send(pc)
	}
	return client.send(p)
}
