package extension

import (
	"encoding/xml"

	"dev.sum7.eu/genofire/yaja/server/utils"
	"dev.sum7.eu/genofire/yaja/xmpp"
)

type Presence struct {
	Extension
}

//TODO Draft

func (p *Presence) Spaces() []string { return []string{} }

func (p *Presence) Process(element *xml.StartElement, client *utils.Client) bool {
	log := client.Log.WithField("extension", "presence")

	// iq encode
	var msg xmpp.PresenceClient
	if err := client.In.DecodeElement(&msg, element); err != nil {
		return false
	}
	client.Messages <- &xmpp.PresenceClient{
		ID: msg.ID,
	}
	log.Debug("send")

	return true
}
