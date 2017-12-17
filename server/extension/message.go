package extension

import (
	"encoding/xml"

	"github.com/genofire/yaja/server/utils"
)

type Message struct {
	Extension
}

//TODO Draft

func (m *Message) Spaces() []string { return []string{} }

func (m *Message) Process(element *xml.StartElement, client *utils.Client) bool {
	return false
}
