package toclient

import (
	"encoding/xml"
	"fmt"

	"github.com/genofire/yaja/database"
	"github.com/genofire/yaja/messages"
	"github.com/genofire/yaja/model"
	"github.com/genofire/yaja/server/state"
	"github.com/genofire/yaja/server/utils"
)

type RegisterFormRequest struct {
	Next                  state.State
	domainRegisterAllowed utils.DomainRegisterAllowed
	element               *xml.StartElement
}

// Process message
func (state *RegisterFormRequest) Process(client *utils.Client) (state.State, *utils.Client) {
	client.Log = client.Log.WithField("state", "register form request")
	client.Log.Debug("running")
	defer client.Log.Debug("leave")

	if !state.domainRegisterAllowed(client.JID) {
		client.Log.Error("unpossible to reach this state, register on this domain is not allowed")
		return nil, client
	}

	var msg messages.IQ
	if err := client.In.DecodeElement(&msg, state.element); err != nil {
		client.Log.Warn("is no iq: ", err)
		return state, client
	}
	if msg.Type != messages.IQTypeGet {
		client.Log.Warn("is no get iq")
		return state, client
	}
	if msg.Error != nil {
		client.Log.Warn("iq with error: ", msg.Error.Code)
		return state, client
	}
	type query struct {
		XMLName xml.Name `xml:"query"`
	}
	q := &query{}
	err := xml.Unmarshal(msg.Body, q)

	if q.XMLName.Space != messages.NSIQRegister || err != nil {
		client.Log.Warn("is no iq register: ", err)
		return nil, client
	}
	client.Out.Encode(&messages.IQ{
		Type: messages.IQTypeResult,
		To:   client.JID.String(),
		From: client.JID.Domain,
		ID:   msg.ID,
		Body: []byte(fmt.Sprintf(`<query xmlns='%s'><instructions>
					Choose a username and password for use with this service.
				</instructions>
				<username/>
				<password/>
			</query>`, messages.NSIQRegister)),
	})
	return state.Next, client
}

type RegisterRequest struct {
	Next                  state.State
	database              *database.State
	domainRegisterAllowed utils.DomainRegisterAllowed
}

// Process message
func (state *RegisterRequest) Process(client *utils.Client) (state.State, *utils.Client) {
	client.Log = client.Log.WithField("state", "register request")
	client.Log.Debug("running")
	defer client.Log.Debug("leave")

	if !state.domainRegisterAllowed(client.JID) {
		client.Log.Error("unpossible to reach this state, register on this domain is not allowed")
		return nil, client
	}

	element, err := client.Read()
	if err != nil {
		client.Log.Warn("unable to read: ", err)
		return nil, client
	}
	var msg messages.IQ
	if err = client.In.DecodeElement(&msg, element); err != nil {
		client.Log.Warn("is no iq: ", err)
		return state, client
	}
	if msg.Type != messages.IQTypeGet {
		client.Log.Warn("is no get iq")
		return state, client
	}
	if msg.Error != nil {
		client.Log.Warn("iq with error: ", msg.Error.Code)
		return state, client
	}
	type query struct {
		XMLName  xml.Name `xml:"query"`
		Username string   `xml:"username"`
		Password string   `xml:"password"`
	}
	q := &query{}
	err = xml.Unmarshal(msg.Body, q)
	if err != nil {
		client.Log.Warn("is no iq register: ", err)
		return nil, client
	}

	client.JID.Local = q.Username
	client.Log = client.Log.WithField("jid", client.JID.Full())
	account := model.NewAccount(client.JID, q.Password)
	err = state.database.AddAccount(account)
	if err != nil {
		client.Out.Encode(&messages.IQ{
			Type: messages.IQTypeResult,
			To:   client.JID.String(),
			From: client.JID.Domain,
			ID:   msg.ID,
			Body: []byte(fmt.Sprintf(`<query xmlns='%s'>
					<username>%s</username>
					<password>%s</password>
				</query>`, messages.NSIQRegister, q.Username, q.Password)),
			Error: &messages.Error{
				Code: "409",
				Type: "cancel",
				Any: xml.Name{
					Local: "conflict",
					Space: "urn:ietf:params:xml:ns:xmpp-stanzas",
				},
			},
		})
		client.Log.Warn("database error: ", err)
		return state, client
	}
	client.Out.Encode(&messages.IQ{
		Type: messages.IQTypeResult,
		To:   client.JID.String(),
		From: client.JID.Domain,
		ID:   msg.ID,
	})

	client.Log.Infof("registered client %s", client.JID.Bare())
	return state.Next, client
}
