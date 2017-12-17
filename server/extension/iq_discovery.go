package extension

import (
	"encoding/xml"

	"github.com/genofire/yaja/messages"
	"github.com/genofire/yaja/server/utils"
)

type ExtensionDiscovery struct {
	IQExtension
	GetSpaces func() []string
}

func (ex *ExtensionDiscovery) Spaces() []string { return []string{} }

func (ex *ExtensionDiscovery) Get(msg *messages.IQ, client *utils.Client) bool {
	log := client.Log.WithField("extension", "roster").WithField("id", msg.ID)

	// query encode
	type query struct {
		XMLName xml.Name `xml:"http://jabber.org/protocol/disco#info query"`
		Body    []byte   `xml:",innerxml"`
	}
	q := &query{}
	err := xml.Unmarshal(msg.Body, q)
	if err != nil {
		return false
	}

	// answer query
	q.Body = []byte{}

	// build answer body
	type feature struct {
		XMLName xml.Name `xml:"feature"`
		Var     string   `xml:"var,attr"`
	}
	for _, namespace := range ex.GetSpaces() {
		if namespace == "" {
			continue
		}
		itemByte, err := xml.Marshal(&feature{
			Var: namespace,
		})
		if err != nil {
			log.Warn(err)
			continue
		}
		q.Body = append(q.Body, itemByte...)
	}

	// decode query
	queryByte, err := xml.Marshal(q)
	if err != nil {
		log.Warn(err)
		return false
	}

	// replay
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
