package extension

import (
	"dev.sum7.eu/genofire/yaja/server/utils"
	"dev.sum7.eu/genofire/yaja/xmpp"
	"dev.sum7.eu/genofire/yaja/xmpp/base"
)

type IQPing struct {
	IQExtension
}

func (ex *IQPing) Spaces() []string { return []string{"urn:xmpp:ping"} }

func (ex *IQPing) Get(msg *xmpp.IQClient, client *utils.Client) bool {
	log := client.Log.WithField("extension", "ping").WithField("id", msg.ID)

	if msg.Ping == nil {
		return false
	}

	// reply
	client.Messages <- &xmpp.IQClient{
		Type: xmpp.IQTypeResult,
		To:   client.JID,
		From: xmppbase.NewJID(client.JID.Domain),
		ID:   msg.ID,
	}

	log.Debug("send")

	return true
}
