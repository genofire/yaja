package server

import (
	"encoding/xml"
	"fmt"

	"github.com/genofire/yaja/messages"
	"github.com/genofire/yaja/model"
)

type RegisterFormRequest struct {
	Next    State
	element *xml.StartElement
}

// Process message
func (state *RegisterFormRequest) Process(client *Client) (State, *Client) {
	client.log = client.log.WithField("state", "register form request")
	client.log.Debug("running")
	defer client.log.Debug("leave")

	var msg messages.IQ
	if err := client.in.DecodeElement(&msg, state.element); err != nil {
		client.log.Warn("is no iq: ", err)
		return state, client
	}
	if msg.Type != messages.IQTypeGet {
		client.log.Warn("is no get iq")
		return state, client
	}
	if msg.Error != nil {
		client.log.Warn("iq with error: ", msg.Error.Code)
		return state, client
	}
	type query struct {
		XMLName xml.Name `xml:"query"`
	}
	q := &query{}
	err := xml.Unmarshal(msg.Body, q)

	if q.XMLName.Space != messages.NSIQRegister || err != nil {
		client.log.Warn("is no iq register: ", err)
		return nil, client
	}
	client.out.Encode(&messages.IQ{
		Type: messages.IQTypeResult,
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
	Next State
}

// Process message
func (state *RegisterRequest) Process(client *Client) (State, *Client) {
	client.log = client.log.WithField("state", "register request")
	client.log.Debug("running")
	defer client.log.Debug("leave")

	element, err := client.Read()
	if err != nil {
		client.log.Warn("unable to read: ", err)
		return nil, client
	}
	var msg messages.IQ
	if err = client.in.DecodeElement(&msg, element); err != nil {
		client.log.Warn("is no iq: ", err)
		return state, client
	}
	if msg.Type != messages.IQTypeGet {
		client.log.Warn("is no get iq")
		return state, client
	}
	if msg.Error != nil {
		client.log.Warn("iq with error: ", msg.Error.Code)
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
		client.log.Warn("is no iq register: ", err)
		return nil, client
	}

	client.jid.Local = q.Username
	client.log = client.log.WithField("jid", client.jid.Full())
	account := model.NewAccount(client.jid, q.Password)
	err = client.Server.Database.AddAccount(account)
	if err != nil {
		client.out.Encode(&messages.IQ{
			Type: messages.IQTypeResult,
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
		client.log.Warn("database error: ", err)
		return state, client
	}
	client.account = account
	client.out.Encode(&messages.IQ{
		Type: messages.IQTypeResult,
		ID:   msg.ID,
	})

	client.log.Infof("registered client %s", client.jid.Bare())
	return state.Next, client
}
