package extension

import (
	"encoding/xml"

	"dev.sum7.eu/genofire/yaja/messages"
	"dev.sum7.eu/genofire/yaja/server/utils"
)

//TODO Draft

type IQLast struct {
	IQExtension
}

func (ex *IQLast) Spaces() []string { return []string{"jabber:iq:last"} }

func (ex *IQLast) Get(msg *messages.IQClient, client *utils.Client) bool {
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
	client.Messages <- &messages.IQClient{
		Type: messages.IQTypeResult,
		To:   client.JID.String(),
		From: client.JID.Domain,
		ID:   msg.ID,
		Body: queryByte,
	}

	log.Debug("send")

	return true
}
