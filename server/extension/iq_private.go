package extension

import (
	"encoding/xml"

	"dev.sum7.eu/genofire/yaja/messages"
	"dev.sum7.eu/genofire/yaja/server/utils"
)

type IQPrivate struct {
	IQExtension
}

type iqPrivateQuery struct {
	XMLName xml.Name `xml:"jabber:iq:private query"`
	Body    []byte   `xml:",innerxml"`
}

type iqPrivateExtension interface {
	Handle(*messages.IQ, *iqPrivateQuery, *utils.Client) bool
}

func (ex *IQPrivate) Spaces() []string { return []string{"jabber:iq:private"} }

func (ex *IQPrivate) Get(msg *messages.IQ, client *utils.Client) bool {
	log := client.Log.WithField("extension", "private").WithField("id", msg.ID)

	// query encode
	q := &iqPrivateQuery{}
	if err := xml.Unmarshal(msg.Body, q); err != nil {
		return false
	}

	// run every extensions
	count := 0
	for _, e := range []iqPrivateExtension{
		&IQPrivateMetacontact{},
		&IQPrivateRoster{},
		&IQPrivateBookmark{},
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
