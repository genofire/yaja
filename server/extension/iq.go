package extension

import (
	"encoding/xml"

	"dev.sum7.eu/genofire/yaja/messages"
	"dev.sum7.eu/genofire/yaja/server/utils"
)

type IQExtensions []IQExtension

type IQExtension interface {
	Extension
	Get(*messages.IQClient, *utils.Client) bool
	Set(*messages.IQClient, *utils.Client) bool
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
	var msg messages.IQClient
	if err := client.In.DecodeElement(&msg, element); err != nil {
		return false
	}

	log = log.WithField("id", msg.ID)

	// run every extensions
	count := 0
	for _, extension := range iex {
		switch msg.Type {
		case messages.IQTypeGet:
			if extension.Get(&msg, client) {
				count++
			}
		case messages.IQTypeSet:
			if extension.Set(&msg, client) {
				count++
			}
		}
	}

	// not extensions found
	if count != 1 {
		log.Debugf("%s - %s: %v", msg.XMLName.Space, msg.Type, msg.Other)
	}

	return true
}
