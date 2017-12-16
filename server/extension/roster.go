package extension

import (
	"encoding/xml"

	"github.com/genofire/yaja/database"
	"github.com/genofire/yaja/messages"
	"github.com/genofire/yaja/server/utils"
)

type Roster struct {
	Extension
	Database *database.State
}

func (r *Roster) Process(element *xml.StartElement, client *utils.Client) bool {
	var msg messages.IQ
	if err := client.In.DecodeElement(&msg, element); err != nil {
		client.Log.Warn("is no iq: ", err)
		return false
	}
	if msg.Type != messages.IQTypeGet {
		client.Log.Warn("is no get iq")
		return false
	}
	return true
}
