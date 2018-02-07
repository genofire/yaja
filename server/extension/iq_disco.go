package extension

import (
	"encoding/xml"

	"dev.sum7.eu/genofire/yaja/database"
	"dev.sum7.eu/genofire/yaja/messages"
	"dev.sum7.eu/genofire/yaja/server/utils"
)

type IQDisco struct {
	IQExtension
	Database *database.State
}

func (ex *IQDisco) Spaces() []string { return []string{"http://jabber.org/protocol/disco#items"} }

func (ex *IQDisco) Get(msg *messages.IQ, client *utils.Client) bool {
	log := client.Log.WithField("extension", "disco-item").WithField("id", msg.ID)

	// query encode
	type query struct {
		XMLName xml.Name `xml:"http://jabber.org/protocol/disco#items query"`
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
	if acc := ex.Database.GetAccount(client.JID); acc != nil {
		for jid, _ := range acc.Bookmarks {
			itemByte, err := xml.Marshal(&item{
				JID: jid,
			})
			if err != nil {
				log.Warn(err)
				continue
			}
			q.Body = append(q.Body, itemByte...)
		}
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
