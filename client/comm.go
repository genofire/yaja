package client

import (
	"encoding/xml"

	"dev.sum7.eu/genofire/yaja/xmpp"
)

func (client *Client) Read() (*xml.StartElement, error) {
	for {
		nextToken, err := client.in.Token()
		if err != nil {
			return nil, err
		}
		switch nextToken.(type) {
		case xml.StartElement:
			element := nextToken.(xml.StartElement)
			client.Logging.Debug("recv xml: ", xmpp.XMLStartElementToString(&element))
			return &element, nil
		}
	}
}
func (client *Client) Decode(p interface{}, element *xml.StartElement) error {
	err := client.in.DecodeElement(p, element)
	if err != nil {
		client.Logging.Debugf("decode failed xml: %s to: %v", xmpp.XMLStartElementToString(element), p)
	} else {
		client.Logging.Debugf("decode xml: %s to: %v with children %s", xmpp.XMLStartElementToString(element), p, xmpp.XMLChildrenString(p))
	}
	return err
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
	if err == nil && iq.Ping != nil {
		client.Logging.Info("ReadElement: auto answer ping")
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
func (client *Client) encode(p interface{}) error {
	err := client.out.Encode(p)
	if err != nil {
		client.Logging.Debugf("encode failed %v", p)
	} else {
		client.Logging.Debugf("encode %v with children %s", p, xmpp.XMLChildrenString(p))
	}
	return err
}

func (client *Client) Send(p interface{}) error {
	msg, ok := p.(*xmpp.MessageClient)
	if ok {
		msg.From = client.JID
		return client.encode(msg)
	}
	iq, ok := p.(*xmpp.IQClient)
	if ok {
		iq.From = client.JID
		return client.encode(iq)
	}
	pc, ok := p.(*xmpp.PresenceClient)
	if ok {
		pc.From = client.JID
		return client.encode(pc)
	}
	return client.encode(p)
}
