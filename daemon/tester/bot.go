package tester

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"dev.sum7.eu/genofire/yaja/xmpp"
	"dev.sum7.eu/genofire/yaja/xmpp/base"
)

func (t *Tester) startBot(status *Status) {
	logger := log.New()
	logger.SetLevel(t.LoggingBots)
	logCTX := logger.WithFields(log.Fields{
		"log": "bot",
		"jid": status.client.JID.Full(),
	})
	go func(status *Status) {
		if err := status.client.Start(); err != nil {
			status.Disconnect(err.Error())
		} else {
			status.Disconnect("safe closed")
		}
	}(status)
	logCTX.Info("start bot")
	defer logCTX.Info("quit bot")
	for {
		element, more := status.client.Recv()
		if !more {
			logCTX.Info("could not recv msg, closed")
			return
		}
		logCTX.Debugf("recv msg %v", element)

		switch element.(type) {
		case *xmpp.PresenceClient:
			pres := element.(*xmpp.PresenceClient)
			sender := pres.From
			logPres := logCTX.WithField("from", sender.Full())
			switch pres.Type {
			case xmpp.PresenceTypeSubscribe:
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
			case xmpp.PresenceTypeSubscribed:
				logPres.Info("recv presence accepted subscribe")
			case xmpp.PresenceTypeUnsubscribe:
				logPres.Info("recv presence remove subscribe")
			case xmpp.PresenceTypeUnsubscribed:
				logPres.Info("recv presence removed subscribe")
			case xmpp.PresenceTypeUnavailable:
				logPres.Debug("recv presence unavailable")
			default:
				logCTX.Warnf("recv presence unsupported: %s -> %s", pres.Type, xmpp.XMLChildrenString(pres))
			}
		case *xmpp.MessageClient:
			msg := element.(*xmpp.MessageClient)
			logMSG := logCTX.WithField("from", msg.From.Full()).WithField("msg-recv", msg.Body)
			msgText := strings.SplitN(msg.Body, " ", 2)
			switch msgText[0] {

			case "ping":
				status.client.Send(xmpp.MessageClient{Type: msg.Type, To: msg.From, Body: "pong"})
				logMSG.Info("answer ping")

			case "admin":
				if len(msgText) == 2 {
					botAdmin(strings.SplitN(msgText[1], " ", 2), logMSG, status, msg.From, botAllowed(t.Admins, status.account.Admins))
				} else {
					status.client.Send(xmpp.MessageClient{Type: msg.Type, To: msg.From, Body: "list, add JID-BARE, del JID-BARE"})
					logMSG.Info("answer admin help")
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
				logMSG.Info("answer disconnect not allowed")

			case "checkmsg":
				if len(msgText) == 2 {
					t.updateConnectionStatus(msg.From, status.client.JID, msgText[1])
				} else {
					logMSG.Debug("undetect")
				}

			default:
				logMSG.Debug("undetect")
			}
		default:
			logCTX.Debug("unhandle")
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
	log.Infof("admin[%s]: %s", from.String(), msg)
}
