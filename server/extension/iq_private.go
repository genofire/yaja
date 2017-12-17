package extension

import (
	"encoding/xml"

	"github.com/genofire/yaja/messages"
	"github.com/genofire/yaja/server/utils"
)

type Private struct {
	IQExtension
}

type privateQuery struct {
	XMLName xml.Name `xml:"jabber:iq:private query"`
	Body    []byte   `xml:",innerxml"`
}

type ioPrivateExtension interface {
	Handle(*messages.IQ, *privateQuery, *utils.Client) bool
}

func (p *Private) Spaces() []string { return []string{"jabber:iq:private"} }

func (p *Private) Get(msg *messages.IQ, client *utils.Client) bool {
	log := client.Log.WithField("extension", "private").WithField("id", msg.ID)

	// query encode
	q := &privateQuery{}
	err := xml.Unmarshal(msg.Body, q)
	if err != nil {
		return false
	}

	// run every extensions
	count := 0
	for _, e := range []ioPrivateExtension{
		&PrivateMetacontact{},
		&PrivateRoster{},
	} {
		if e.Handle(msg, q, client) {
			count++
		}
	}

	// not extensions found
	if count != 1 {
		log.Debug(msg.XMLName.Space, " - ", msg.Type, ": ", string(q.Body))
	}

	return true
}
