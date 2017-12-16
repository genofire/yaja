package extension

import (
	"encoding/xml"

	"github.com/genofire/yaja/server/utils"
)

type Message struct {
	Extension
}

func (m *Message) Process(element *xml.StartElement, client *utils.Client) bool {
	return false
}
