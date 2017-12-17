package extension

import (
	"encoding/xml"

	"github.com/genofire/yaja/messages"
	"github.com/genofire/yaja/server/utils"
)

type IQPrivateMetacontact struct {
	iqPrivateExtension
}

func (ex *IQPrivateMetacontact) Handle(msg *messages.IQ, q *iqPrivateQuery, client *utils.Client) bool {
	log := client.Log.WithField("extension", "private-metacontact").WithField("id", msg.ID)

	// storage encode
	type storage struct {
		XMLName xml.Name `xml:"storage:metacontacts storage"`
	}
	s := &storage{}
	err := xml.Unmarshal(q.Body, s)
	if err != nil {
		return false
	}
	/*
		TODO full implement XEP-0209
		 https://xmpp.org/extensions/xep-0209.html
	*/

	queryByte, err := xml.Marshal(&iqPrivateQuery{
		Body: q.Body,
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
