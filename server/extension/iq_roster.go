package extension

import (
	"encoding/xml"

	"github.com/genofire/yaja/database"
	"github.com/genofire/yaja/messages"
	"github.com/genofire/yaja/server/utils"
)

type Roster struct {
	IQExtension
	Database *database.State
}

func (r *Roster) Spaces() []string { return []string{"jabber:iq:roster"} }

func (r *Roster) Get(msg *messages.IQ, client *utils.Client) bool {
	log := client.Log.WithField("extension", "roster").WithField("id", msg.ID)

	// query encode
	type query struct {
		XMLName xml.Name `xml:"jabber:iq:roster query"`
		Version string   `xml:"ver,attr"`
		Body    []byte   `xml:",innerxml"`
	}
	q := &query{}
	err := xml.Unmarshal(msg.Body, q)
	if err != nil {
		return false
	}

	// answer query
	q.Body = []byte{}
	q.Version = "1"

	// build answer body
	type item struct {
		XMLName xml.Name `xml:"item"`
		JID     string   `xml:"jid,attr"`
	}
	if acc := r.Database.GetAccount(client.JID); acc != nil {
		for jid, _ := range acc.Roster {
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
	client.Out.Encode(&messages.IQ{
		Type: messages.IQTypeResult,
		To:   client.JID.String(),
		From: client.JID.Domain,
		ID:   msg.ID,
		Body: queryByte,
	})

	log.Debug("send")

	return true
}
