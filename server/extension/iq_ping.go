package extension

import (
	"encoding/xml"

	"dev.sum7.eu/genofire/yaja/messages"
	"dev.sum7.eu/genofire/yaja/server/utils"
)

type IQPing struct {
	IQExtension
}

func (ex *IQPing) Spaces() []string { return []string{"urn:xmpp:ping"} }

func (ex *IQPing) Get(msg *messages.IQ, client *utils.Client) bool {
	log := client.Log.WithField("extension", "ping").WithField("id", msg.ID)

	// ping encode
	type ping struct {
		XMLName xml.Name `xml:"urn:xmpp:ping ping"`
	}
	pq := &ping{}
	if err := xml.Unmarshal(msg.Body, pq); err != nil {
		return false
	}

	// reply
	client.Messages <- &messages.IQ{
		Type: messages.IQTypeResult,
		To:   client.JID.String(),
		From: client.JID.Domain,
		ID:   msg.ID,
	}

	log.Debug("send")

	return true
}
