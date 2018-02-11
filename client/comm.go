package client

import (
	"encoding/xml"

	"dev.sum7.eu/genofire/yaja/messages"
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
			client.Logging.Debug("recv xml: ", messages.XMLStartElementToString(&element))
			return &element, nil
		}
	}
}
func (client *Client) Decode(p interface{}, element *xml.StartElement) error {
	err := client.in.DecodeElement(p, element)
	if err != nil {
		client.Logging.Debugf("decode failed xml: %s to: %v", messages.XMLStartElementToString(element), p)
	} else {
		client.Logging.Debugf("decode xml: %s to: %v with children %s", messages.XMLStartElementToString(element), p, messages.XMLChildrenString(p))
	}
	return err
}
func (client *Client) ReadDecode(p interface{}) error {
	element, err := client.Read()
	if err != nil {
		return err
	}
	var iq *messages.IQClient
	iq, ok := p.(*messages.IQClient)
	if !ok {
		iq = &messages.IQClient{}
	}
	err = client.Decode(iq, element)
	if err == nil && iq.Ping != nil {
		client.Logging.Info("ReadElement: auto answer ping")
		iq.Type = messages.IQTypeResult
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
		client.Logging.Debugf("encode %v with children %s", p, messages.XMLChildrenString(p))
	}
	return err
}

func (client *Client) Send(p interface{}) error {
	msg, ok := p.(*messages.MessageClient)
	if ok {
		msg.From = client.JID
		return client.encode(msg)
	}
	iq, ok := p.(*messages.IQClient)
	if ok {
		iq.From = client.JID
		return client.encode(iq)
	}
	pc, ok := p.(*messages.PresenceClient)
	if ok {
		pc.From = client.JID
		return client.encode(pc)
	}
	return client.encode(p)
}
