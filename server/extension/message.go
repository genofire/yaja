package extension

import (
	"encoding/xml"

	"dev.sum7.eu/genofire/yaja/server/utils"
)

type Message struct {
	Extension
}

//TODO Draft

func (m *Message) Spaces() []string { return []string{} }

func (m *Message) Process(element *xml.StartElement, client *utils.Client) bool {
	return false
}
