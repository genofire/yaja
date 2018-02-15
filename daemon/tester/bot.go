package tester

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"dev.sum7.eu/genofire/yaja/xmpp"
	"dev.sum7.eu/genofire/yaja/xmpp/base"
)

func (t *Tester) StartBot(status *Status) {
	logger := log.New()
	logger.SetLevel(t.LoggingBots)
	logCTX := logger.WithFields(log.Fields{
		"type": "bot",
		"jid":  status.client.JID.Full(),
	})
	for {

		element, err := status.client.Read()
		if err != nil {
			status.Disconnect(fmt.Sprintf("could not read any more data from socket: %s", err))
			return
		}

		errMSG := &xmpp.StreamError{}
		err = status.client.Decode(errMSG, element)
		if err == nil {
			status.Disconnect(fmt.Sprintf("recv stream error: %s: %s -> %s", errMSG.Text, xmpp.XMLChildrenString(errMSG.StreamErrorGroup), xmpp.XMLChildrenString(errMSG.Other)))
			return
		}

		iq := &xmpp.IQClient{}
		err = status.client.Decode(iq, element)
		if err == nil {
			if iq.Ping != nil {
				logCTX.Debug("answer ping")
				iq.Type = xmpp.IQTypeResult
				iq.To = iq.From
				iq.From = status.client.JID
				status.client.Send(iq)
			} else {
				logCTX.Warnf("recv iq unsupport: %s", xmpp.XMLChildrenString(iq))
			}
			continue
		}

		pres := &xmpp.PresenceClient{}
		err = status.client.Decode(pres, element)
		if err == nil {
			sender := pres.From
			logPres := logCTX.WithField("from", sender.Full())
			if pres.Type == xmpp.PresenceTypeSubscribe {
				logPres.Debugf("recv presence subscribe")
				pres.Type = xmpp.PresenceTypeSubscribed
				pres.To = sender
				pres.From = nil
				status.client.Send(pres)
				logPres.Debugf("accept new subscribe")

				pres.Type = xmpp.PresenceTypeSubscribe
				pres.ID = ""
				status.client.Send(pres)
				logPres.Info("request also subscribe")
			} else if pres.Type == xmpp.PresenceTypeSubscribed {
				logPres.Info("recv presence accepted subscribe")
			} else if pres.Type == xmpp.PresenceTypeUnsubscribe {
				logPres.Info("recv presence remove subscribe")
			} else if pres.Type == xmpp.PresenceTypeUnsubscribed {
				logPres.Info("recv presence removed subscribe")
			} else if pres.Type == xmpp.PresenceTypeUnavailable {
				logPres.Debug("recv presence unavailable")
			} else {
				logCTX.Warnf("recv presence unsupported: %s -> %s", pres.Type, xmpp.XMLChildrenString(pres))
			}
			continue
		}

		msg := &xmpp.MessageClient{}
		err = status.client.Decode(msg, element)
		if err != nil {
			logCTX.Warnf("unsupport xml recv: %s <-> %v", err, element)
			continue
		}
		logCTX = logCTX.WithField("from", msg.From.Full()).WithField("msg-recv", msg.Body)
		if msg.Error != nil {
			if msg.Error.Type == "auth" {
				status.Disconnect("recv msg with error not auth")
				return
			}
			logCTX.Debugf("recv msg with error %s[%s]: %s -> %s -> %s", msg.Error.Code, msg.Error.Type, msg.Error.Text, xmpp.XMLChildrenString(msg.Error.StanzaErrorGroup), xmpp.XMLChildrenString(msg.Error.Other))
			continue

		}

		msgText := strings.SplitN(msg.Body, " ", 2)
		switch msgText[0] {

		case "ping":
			status.client.Send(xmpp.MessageClient{Type: msg.Type, To: msg.From, Body: "pong"})
		case "admin":
			if len(msgText) == 2 {
				botAdmin(strings.SplitN(msgText[1], " ", 2), logCTX, status, msg.From, botAllowed(t.Admins, status.account.Admins))
			} else {
				status.client.Send(xmpp.MessageClient{Type: msg.Type, To: msg.From, Body: "list, add JID-BARE, del JID-BARE"})
			}
		case "disconnect":
			first := true
			allAdmins := ""
			isAdmin := false
			fromBare := msg.From
			for _, jid := range botAllowed(t.Admins, status.account.Admins) {
				if first {
					first = false
				} else {
					allAdmins += ", "
				}
				allAdmins += jid.Bare().String()
				if jid.Bare().IsEqual(fromBare) {
					isAdmin = true
					status.client.Send(xmpp.MessageClient{Type: msg.Type, To: jid, Body: "last message, disconnect requested by " + fromBare.String()})

				}
			}
			if isAdmin {
				status.Disconnect(fmt.Sprintf("disconnect by admin '%s'", fromBare.String()))
				return
			}
			status.client.Send(xmpp.MessageClient{Type: msg.Type, To: msg.From, Body: "not allowed, ask " + allAdmins})

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
func botAllowed(list []*xmppbase.JID, toConvert map[string]interface{}) []*xmppbase.JID {
	alist := list
	for jid := range toConvert {
		alist = append(alist, xmppbase.NewJID(jid))
	}
	return alist
}

func botAdmin(cmd []string, log *log.Entry, status *Status, from *xmppbase.JID, allowed []*xmppbase.JID) {
	msg := ""
	if len(cmd) == 2 {
		isAdmin := false
		for _, jid := range allowed {
			if jid.Bare() == from.Bare() {
				isAdmin = true
			}
		}
		if status.account.Admins == nil {
			status.account.Admins = make(map[string]interface{})
		}
		if !isAdmin {
			msg = "not allowed"
		} else if cmd[0] == "add" {
			status.account.Admins[cmd[1]] = true
			msg = "ack"
		} else if cmd[0] == "del" {
			delete(status.account.Admins, cmd[1])
			msg = "ack"
		} else {
			msg = "unknown command"
		}
	} else {
		if len(cmd) == 1 && cmd[0] == "list" {
			for jid := range status.account.Admins {
				if msg == "" {
					msg += "admins are: " + jid
				} else {
					msg += ", " + jid
				}
			}
		} else {
			msg = "unknown command"
		}
	}
	status.client.Send(xmpp.MessageClient{Type: xmpp.MessageTypeChat, To: from, Body: msg})
}
