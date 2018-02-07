package extension

import (
	"encoding/xml"

	"dev.sum7.eu/genofire/yaja/messages"
	"dev.sum7.eu/genofire/yaja/server/utils"
)

type IQPrivateBookmark struct {
	iqPrivateExtension
}

func (ex *IQPrivateBookmark) Handle(msg *messages.IQ, q *iqPrivateQuery, client *utils.Client) bool {
	log := client.Log.WithField("extension", "private").WithField("id", msg.ID)

	// storage encode
	type storage struct {
		XMLName xml.Name `xml:"storage:bookmarks storage"`
	}
	s := &storage{}
	if err := xml.Unmarshal(q.Body, s); err != nil {
		return false
	}
	/*
		TODO full implement
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
