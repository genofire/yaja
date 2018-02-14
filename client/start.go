package client

import (
	"fmt"

	"dev.sum7.eu/genofire/yaja/xmpp"
)

var DefaultChannelSize = 30

func (client *Client) Start() error {
	client.iq = make(chan *xmpp.IQClient, DefaultChannelSize)
	client.presence = make(chan *xmpp.PresenceClient, DefaultChannelSize)
	client.mesage = make(chan *xmpp.MessageClient, DefaultChannelSize)
	client.reply = make(map[string]chan *xmpp.IQClient)

	for {

		element, err := client.Read()
		if err != nil {
			return err
		}

		errMSG := &xmpp.StreamError{}
		err = client.Decode(errMSG, element)
		if err == nil {
			return fmt.Errorf("recv stream error: %s: %s -> %s", errMSG.Text, xmpp.XMLChildrenString(errMSG.StreamErrorGroup), xmpp.XMLChildrenString(errMSG.Other))
		}

		iq := &xmpp.IQClient{}
		err = client.Decode(iq, element)
		if err == nil {
			if iq.Ping != nil {
				client.Logging.Debug("answer ping")
				iq.Type = xmpp.IQTypeResult
				iq.To = iq.From
				iq.From = client.JID
				client.Send(iq)
			} else {
				if client.skipError && iq.Error != nil {
					continue
				}
				if ch, ok := client.reply[iq.ID]; ok {
					delete(client.reply, iq.ID)
					ch <- iq
					continue
				}
				client.iq <- iq
			}
			continue
		}

		pres := &xmpp.PresenceClient{}
		err = client.Decode(pres, element)
		if err == nil {
			if client.skipError && pres.Error != nil {
				continue
			}
			client.presence <- pres
			continue
		}

		msg := &xmpp.MessageClient{}
		err = client.Decode(msg, element)
		if err == nil {
			if client.skipError && msg.Error != nil {
				continue
			}
			client.mesage <- msg
			continue
		}
		client.Logging.Warnf("unsupport xml recv: %v", element)
	}
}

func (client *Client) SendRecv(iq *xmpp.IQClient) *xmpp.IQClient {
	if iq.ID == "" {
		iq.ID = xmpp.CreateCookieString()
	}
	ch := make(chan *xmpp.IQClient, 1)
	client.reply[iq.ID] = ch
	client.Send(iq)
	defer close(ch)
	return <-ch
}

func (client *Client) RecvIQ() *xmpp.IQClient {
	return <-client.iq
}

func (client *Client) RecvPresence() *xmpp.PresenceClient {
	return <-client.presence
}

func (client *Client) RecvMessage() *xmpp.MessageClient {
	return <-client.mesage
}
