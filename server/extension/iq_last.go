package extension

import (
	"encoding/xml"

	"github.com/genofire/yaja/messages"
	"github.com/genofire/yaja/server/utils"
)

//TODO Draft

type IQLast struct {
	IQExtension
}

func (ex *IQLast) Spaces() []string { return []string{"jabber:iq:last"} }

func (ex *IQLast) Get(msg *messages.IQ, client *utils.Client) bool {
	log := client.Log.WithField("extension", "last").WithField("id", msg.ID)

	// query encode
	type query struct {
		XMLName xml.Name `xml:"jabber:iq:last query"`
		Seconds uint     `xml:"seconds,attr,omitempty"`
		Body    []byte   `xml:",innerxml"`
	}
	q := &query{}
	if err := xml.Unmarshal(msg.Body, q); err != nil {
		return false
	}

	// answer query
	q.Body = []byte{}

	// build answer body
	type item struct {
		XMLName xml.Name `xml:"item"`
		JID     string   `xml:"jid,attr"`
	}
	// decode query
	queryByte, err := xml.Marshal(q)
	if err != nil {
		log.Warn(err)
		return false
	}

	// reply
	client.Messages <- &messages.IQ{
		Type: messages.IQTypeResult,
		To:   client.JID.String(),
		From: client.JID.Domain,
		ID:   msg.ID,
		Body: queryByte,
	}

	log.Debug("send")

	return true
}
