package client

import (
	"fmt"

	"dev.sum7.eu/genofire/yaja/messages"
)

var DefaultChannelSize = 30

func (client *Client) Start() error {
	client.iq = make(chan *messages.IQClient, DefaultChannelSize)
	client.presence = make(chan *messages.PresenceClient, DefaultChannelSize)
	client.mesage = make(chan *messages.MessageClient, DefaultChannelSize)
	client.reply = make(map[string]chan *messages.IQClient)

	for {

		element, err := client.Read()
		if err != nil {
			return err
		}

		errMSG := &messages.StreamError{}
		err = client.Decode(errMSG, element)
		if err == nil {
			return fmt.Errorf("recv stream error: %s: %s -> %s", errMSG.Text, messages.XMLChildrenString(errMSG.StreamErrorGroup), messages.XMLChildrenString(errMSG.Other))
		}

		iq := &messages.IQClient{}
		err = client.Decode(iq, element)
		if err == nil {
			if iq.Ping != nil {
				client.Logging.Debug("answer ping")
				iq.Type = messages.IQTypeResult
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

		pres := &messages.PresenceClient{}
		err = client.Decode(pres, element)
		if err == nil {
			if client.skipError && pres.Error != nil {
				continue
			}
			client.presence <- pres
			continue
		}

		msg := &messages.MessageClient{}
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

func (client *Client) SendRecv(iq *messages.IQClient) *messages.IQClient {
	if iq.ID == "" {
		iq.ID = messages.CreateCookieString()
	}
	ch := make(chan *messages.IQClient, 1)
	client.reply[iq.ID] = ch
	client.Send(iq)
	defer close(ch)
	return <-ch
}

func (client *Client) RecvIQ() *messages.IQClient {
	return <-client.iq
}

func (client *Client) RecvPresence() *messages.PresenceClient {
	return <-client.presence
}

func (client *Client) RecvMessage() *messages.MessageClient {
	return <-client.mesage
}
