package tester

import (
	"strings"

	log "github.com/sirupsen/logrus"

	"dev.sum7.eu/genofire/yaja/messages"
	"dev.sum7.eu/genofire/yaja/server/utils"
)

func (t *Tester) StartBot(status *Status) {
	for {
		logCTX := log.WithField("jid", status.client.JID.Full())

		element, err := status.client.Read()
		if err != nil {
			logCTX.Errorf("read client %s", err)
			status.client.Close()
			status.Login = false
			return
		}

		errMSG := &messages.StreamError{}
		err = status.client.In.DecodeElement(errMSG, element)
		if err == nil {
			logCTX.Errorf("recv stream error: %s: %v", errMSG.Text, errMSG.Any)
			status.client.Close()
			status.Login = false
			return
		}

		iq := &messages.IQClient{}
		err = status.client.In.DecodeElement(iq, element)
		if err == nil {
			if iq.Ping != nil {
				logCTX.Debug("answer ping")
				iq.Type = messages.IQTypeResult
				iq.To = iq.From
				iq.From = status.client.JID
				status.client.Out.Encode(iq)
			} else {
				logCTX.Warnf("unsupport iq recv: %v", iq)
			}
			continue
		}

		pres := &messages.PresenceClient{}
		err = status.client.In.DecodeElement(pres, element)
		if err == nil {
			sender := pres.From
			logPres := logCTX.WithField("from", sender.Full())
			if pres.Type == messages.PresenceTypeSubscribe {
				logPres.Debugf("recv subscribe")
				pres.Type = messages.PresenceTypeSubscribed
				pres.To = sender
				pres.From = nil
				status.client.Out.Encode(pres)
				logPres.Debugf("accept new subscribe")

				pres.Type = messages.PresenceTypeSubscribe
				pres.ID = utils.CreateCookieString()
				status.client.Out.Encode(pres)
				logPres.Info("request also subscribe")
			} else if pres.Type == messages.PresenceTypeSubscribed {
				logPres.Info("recv accepted subscribe")
			} else if pres.Type == messages.PresenceTypeUnsubscribe {
				logPres.Info("recv remove subscribe")
			} else if pres.Type == messages.PresenceTypeUnsubscribed {
				logPres.Info("recv removed subscribe")
			} else {
				logCTX.Warnf("unsupported presence recv: %v", pres)
			}
			continue
		}

		msg := &messages.MessageClient{}
		err = status.client.In.DecodeElement(msg, element)
		if err != nil {
			logCTX.Warnf("unsupport xml recv: %s <-> %v", err, element)
			continue
		}
		logCTX = logCTX.WithField("from", msg.From.Full()).WithField("msg-recv", msg.Body)
		if msg.Error != nil {
			logCTX.Debugf("recv msg with error %s[%s]: %s -> %v -> %v", msg.Error.Code, msg.Error.Type, msg.Error.Text, msg.Error.StanzaErrorGroup, msg.Error.Other)
			continue

		}

		msgText := strings.SplitN(msg.Body, " ", 2)
		switch msgText[0] {

		case "ping":
			status.client.Send(messages.MessageClient{Type: msg.Type, To: msg.From, Body: "pong"})

		case "checkmsg":
			if len(msgText) == 2 {
				t.UpdateConnectionStatus(msg.From, status.client.JID, msgText[1])
			} else {
				logCTX.Debug("undetect")
			}

		default:
			logCTX.Debug("undetect")
		}
	}
}
