package extension

import (
	"dev.sum7.eu/genofire/yaja/messages"
	"dev.sum7.eu/genofire/yaja/model"
	"dev.sum7.eu/genofire/yaja/server/utils"
)

type IQPing struct {
	IQExtension
}

func (ex *IQPing) Spaces() []string { return []string{"urn:xmpp:ping"} }

func (ex *IQPing) Get(msg *messages.IQClient, client *utils.Client) bool {
	log := client.Log.WithField("extension", "ping").WithField("id", msg.ID)

	if msg.Ping == nil {
		return false
	}

	// reply
	client.Messages <- &messages.IQClient{
		Type: messages.IQTypeResult,
		To:   client.JID,
		From: model.NewJID(client.JID.Domain),
		ID:   msg.ID,
	}

	log.Debug("send")

	return true
}
