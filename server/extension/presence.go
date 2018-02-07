package extension

import (
	"encoding/xml"

	"dev.sum7.eu/genofire/yaja/messages"
	"dev.sum7.eu/genofire/yaja/server/utils"
)

type Presence struct {
	Extension
}

//TODO Draft

func (p *Presence) Spaces() []string { return []string{} }

func (p *Presence) Process(element *xml.StartElement, client *utils.Client) bool {
	log := client.Log.WithField("extension", "presence")

	// iq encode
	var msg messages.PresenceClient
	if err := client.In.DecodeElement(&msg, element); err != nil {
		return false
	}
	client.Messages <- &messages.PresenceClient{
		ID: msg.ID,
	}
	log.Debug("send")

	return true
}
