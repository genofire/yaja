package extension

import (
	"encoding/xml"

	"github.com/genofire/yaja/messages"
	"github.com/genofire/yaja/server/utils"
)

type Ping struct {
	IQExtension
}

func (p *Ping) Spaces() []string { return []string{"urn:xmpp:ping"} }

func (p *Ping) Get(msg *messages.IQ, client *utils.Client) bool {
	log := client.Log.WithField("extension", "ping").WithField("id", msg.ID)

	// ping encode
	type ping struct {
		XMLName xml.Name `xml:"urn:xmpp:ping ping"`
	}
	pq := &ping{}
	err := xml.Unmarshal(msg.Body, pq)
	if err != nil {
		return false
	}

	// reply
	client.Out.Encode(&messages.IQ{
		Type: messages.IQTypeResult,
		To:   client.JID.String(),
		From: client.JID.Domain,
		ID:   msg.ID,
	})

	log.Debug("send")

	return true
}