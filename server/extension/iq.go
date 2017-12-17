package extension

import (
	"encoding/xml"

	"github.com/genofire/yaja/messages"
	"github.com/genofire/yaja/server/utils"
)

type IQExtensions []IQExtension

type IQExtension interface {
	Extension
	Get(*messages.IQ, *utils.Client) bool
	Set(*messages.IQ, *utils.Client) bool
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
	var msg messages.IQ
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
		log.Debug(msg.XMLName.Space, " - ", msg.Type, ": ", string(msg.Body))
	}

	return true
}
