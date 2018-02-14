package extension

import (
	"encoding/xml"

	"dev.sum7.eu/genofire/yaja/server/utils"
	"dev.sum7.eu/genofire/yaja/xmpp"
)

type IQExtensions []IQExtension

type IQExtension interface {
	Extension
	Get(*xmpp.IQClient, *utils.Client) bool
	Set(*xmpp.IQClient, *utils.Client) bool
}

func (iex IQExtensions) Spaces() (result []string) {
	for _, extension := range iex {
		spaces := extension.Spaces()
		result = append(result, spaces...)
	}
	return result
}

func (iex IQExtensions) Process(element *xml.StartElement, client *utils.Client) bool {
	log := client.Log.WithField("extension", "iq")

	// iq encode
	var msg xmpp.IQClient
	if err := client.In.DecodeElement(&msg, element); err != nil {
		return false
	}

	log = log.WithField("id", msg.ID)

	// run every extensions
	count := 0
	for _, extension := range iex {
		switch msg.Type {
		case xmpp.IQTypeGet:
			if extension.Get(&msg, client) {
				count++
			}
		case xmpp.IQTypeSet:
			if extension.Set(&msg, client) {
				count++
			}
		}
	}

	// not extensions found
	if count != 1 {
		log.Debugf("%s - %s: %s", msg.XMLName.Space, msg.Type, xmpp.XMLChildrenString(msg.Other))
	}

	return true
}
