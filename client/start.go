package client

import (
	"errors"
	"fmt"

	"dev.sum7.eu/genofire/yaja/xmpp"
)

var DefaultChannelSize = 30

func (client *Client) Start() error {
	if client.msg == nil {
		client.msg = make(chan interface{}, DefaultChannelSize)
	}
	if client.reply == nil {
		client.reply = make(map[string]chan *xmpp.IQClient)
	}

	defer func() {
		for id, ch := range client.reply {
			delete(client.reply, id)
			close(ch)
		}
		client.reply = nil
		close(client.msg)
		client.Logging.Info("client.Start: close")
	}()

	client.Logging.Info("client.Start: start")

	for {

		element, err := client.Read()
		if err != nil {
			return err
		}
		client.Logging.Debugf("client.Start: recv msg %v", xmpp.XMLStartElementToString(element))

		errMSG := &xmpp.StreamError{}
		err = client.Decode(errMSG, element)
		if err == nil {
			return fmt.Errorf("recv stream error: %s: %s -> %s", errMSG.Text, xmpp.XMLChildrenString(errMSG.StreamErrorGroup), xmpp.XMLChildrenString(errMSG.Other))
		}

		iq := &xmpp.IQClient{}
		err = client.Decode(iq, element)
		if err == nil {
			if iq.Ping != nil && iq.Type == xmpp.IQTypeGet {
				client.Logging.Info("client.Start: answer ping")
				iq.Type = xmpp.IQTypeResult
				iq.To = iq.From
				iq.From = client.JID
				client.Send(iq)
			} else {
				if ch, ok := client.reply[iq.ID]; ok {
					delete(client.reply, iq.ID)
					//TODO is this usefull?
					go func() { ch <- iq }()
					continue
				}
				if client.SkipError && iq.Error != nil {
					errStr, err := errorString(iq.Error)
					if err != nil {
						return err
					}
					if errStr != "" {
						client.Logging.WithField("to", iq.To.String()).Error(errStr)
					}
					continue
				}
				client.msg <- iq
			}
			continue
		}

		pres := &xmpp.PresenceClient{}
		err = client.Decode(pres, element)
		if err == nil {
			if client.SkipError && pres.Error != nil {
				errStr, err := errorString(pres.Error)
				if err != nil {
					return err
				}
				if errStr != "" {
					client.Logging.WithField("to", pres.To.String()).Error(errStr)
				}
				continue
			}
			client.msg <- pres
			continue
		}

		msg := &xmpp.MessageClient{}
		err = client.Decode(msg, element)
		if err == nil {
			if client.SkipError && msg.Error != nil {
				errStr, err := errorString(msg.Error)
				if err != nil {
					return err
				}
				if errStr != "" {
					client.Logging.WithField("to", msg.To.String()).Error(errStr)
				}
				continue
			}
			client.msg <- msg
			continue
		}
		client.Logging.Warnf("unsupport xml recv: %v", element)
	}
}
func errorString(e *xmpp.ErrorClient) (string, error) {
	str := fmt.Sprintf("[%s] %s", e.Type, xmpp.XMLChildrenString(e))
	if e.Text != nil {
		str = fmt.Sprintf("[%s] %s -> %s", e.Type, e.Text.Body, xmpp.XMLChildrenString(e))
	}
	if e.Type == xmpp.ErrorTypeAuth {
		return "", errors.New(str)
	}
	if e.RemoteServerNotFound != nil {
		return "", nil
	}
	return str, nil
}

func (client *Client) SendRecv(sendIQ *xmpp.IQClient) (*xmpp.IQClient, error) {
	if sendIQ.ID == "" {
		sendIQ.ID = xmpp.CreateCookieString()
	}
	ch := make(chan *xmpp.IQClient)
	if client.reply == nil {
		return nil, errors.New("client.SendRecv: not init (run client.Start)")
	}
	if client.reply == nil {
		client.reply = make(map[string]chan *xmpp.IQClient)
	}
	client.reply[sendIQ.ID] = ch
	client.Send(sendIQ)
	iq := <-ch
	close(ch)
	if iq.Error != nil {
		_, err := errorString(iq.Error)
		if err != nil {
			return nil, err
		}
	}
	return iq, nil
}

func (client *Client) Recv() (msg interface{}, more bool) {
	if client == nil {
		return nil, false
	}
	if client.msg == nil {
		client.msg = make(chan interface{}, DefaultChannelSize)
	}
	msg, more = <-client.msg
	return
}
