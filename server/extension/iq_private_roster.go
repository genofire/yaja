package extension

import (
	"encoding/xml"

	"github.com/genofire/yaja/messages"
	"github.com/genofire/yaja/server/utils"
)

type IQPrivateRoster struct {
	iqPrivateExtension
}

func (ex *IQPrivateRoster) Handle(msg *messages.IQ, q *iqPrivateQuery, client *utils.Client) bool {
	log := client.Log.WithField("extension", "private").WithField("id", msg.ID)

	// roster encode
	type roster struct {
		XMLName xml.Name `xml:"roster:delimiter roster"`
		Body    []byte   `xml:",innerxml"`
	}
	r := &roster{}
	if err := xml.Unmarshal(q.Body, r); err != nil {
		return false
	}

	rosterByte, err := xml.Marshal(&roster{
		Body: []byte("::"),
	})
	if err != nil {
		log.Warn(err)
		return true
	}
	queryByte, err := xml.Marshal(&iqPrivateQuery{
		Body: rosterByte,
	})
	if err != nil {
		log.Warn(err)
		return true
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
