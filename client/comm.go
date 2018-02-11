package client

import (
	"encoding/xml"
	"log"

	"dev.sum7.eu/genofire/yaja/messages"
)

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

func (client *Client) ReadElement(p interface{}) error {
	element, err := client.Read()
	if err != nil {
		return err
	}
	var iq *messages.IQClient
	iq, ok := p.(*messages.IQClient)
	if !ok {
		iq = &messages.IQClient{}
	}
	err = client.In.DecodeElement(iq, element)
	if err == nil && iq.Ping != nil {
		log.Println("answer ping")
		iq.Type = messages.IQTypeResult
		iq.To = iq.From
		iq.From = client.JID
		client.Out.Encode(iq)
		return nil
	}
	if ok {
		return err
	}
	return client.In.DecodeElement(p, element)
}

func (client *Client) Send(p interface{}) error {
	msg, ok := p.(*messages.MessageClient)
	if ok {
		msg.From = client.JID
		return client.Out.Encode(msg)
	}
	iq, ok := p.(*messages.IQClient)
	if ok {
		iq.From = client.JID
		return client.Out.Encode(iq)
	}
	pc, ok := p.(*messages.PresenceClient)
	if ok {
		pc.From = client.JID
		return client.Out.Encode(pc)
	}
	return client.Out.Encode(p)
}
