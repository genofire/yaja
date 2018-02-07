package toclient

import (
	"encoding/xml"
	"fmt"

	"dev.sum7.eu/genofire/yaja/database"
	"dev.sum7.eu/genofire/yaja/messages"
	"dev.sum7.eu/genofire/yaja/model"
	"dev.sum7.eu/genofire/yaja/server/state"
	"dev.sum7.eu/genofire/yaja/server/utils"
)

type RegisterFormRequest struct {
	Next                  state.State
	Client                *utils.Client
	domainRegisterAllowed utils.DomainRegisterAllowed
	element               *xml.StartElement
}

// Process message
func (state *RegisterFormRequest) Process() state.State {
	state.Client.Log = state.Client.Log.WithField("state", "register form request")
	state.Client.Log.Debug("running")
	defer state.Client.Log.Debug("leave")

	if !state.domainRegisterAllowed(state.Client.JID) {
		state.Client.Log.Error("unpossible to reach this state, register on this domain is not allowed")
		return nil
	}

	var msg messages.IQClient
	if err := state.Client.In.DecodeElement(&msg, state.element); err != nil {
		state.Client.Log.Warn("is no iq: ", err)
		return state
	}
	if msg.Type != messages.IQTypeGet {
		state.Client.Log.Warn("is no get iq")
		return state
	}
	if msg.Error != nil {
		state.Client.Log.Warn("iq with error: ", msg.Error.Code)
		return state
	}
	type query struct {
		XMLName xml.Name `xml:"query"`
	}
	q := &query{}
	err := xml.Unmarshal(msg.Body, q)

	if q.XMLName.Space != messages.NSIQRegister || err != nil {
		state.Client.Log.Warn("is no iq register: ", err)
		return nil
	}
	state.Client.Out.Encode(&messages.IQClient{
		Type: messages.IQTypeResult,
		To:   state.Client.JID.String(),
		From: state.Client.JID.Domain,
		ID:   msg.ID,
		Body: []byte(fmt.Sprintf(`<query xmlns='%s'><instructions>
					Choose a username and password for use with this service.
				</instructions>
				<username/>
				<password/>
			</query>`, messages.NSIQRegister)),
	})
	return state.Next
}

type RegisterRequest struct {
	Next                  state.State
	Client                *utils.Client
	database              *database.State
	domainRegisterAllowed utils.DomainRegisterAllowed
}

// Process message
func (state *RegisterRequest) Process() state.State {
	state.Client.Log = state.Client.Log.WithField("state", "register request")
	state.Client.Log.Debug("running")
	defer state.Client.Log.Debug("leave")

	if !state.domainRegisterAllowed(state.Client.JID) {
		state.Client.Log.Error("unpossible to reach this state, register on this domain is not allowed")
		return nil
	}

	element, err := state.Client.Read()
	if err != nil {
		state.Client.Log.Warn("unable to read: ", err)
		return nil
	}
	var msg messages.IQClient
	if err = state.Client.In.DecodeElement(&msg, element); err != nil {
		state.Client.Log.Warn("is no iq: ", err)
		return state
	}
	if msg.Type != messages.IQTypeGet {
		state.Client.Log.Warn("is no get iq")
		return state
	}
	if msg.Error != nil {
		state.Client.Log.Warn("iq with error: ", msg.Error.Code)
		return state
	}
	type query struct {
		XMLName  xml.Name `xml:"query"`
		Username string   `xml:"username"`
		Password string   `xml:"password"`
	}
	q := &query{}
	err = xml.Unmarshal(msg.Body, q)
	if err != nil {
		state.Client.Log.Warn("is no iq register: ", err)
		return nil
	}

	state.Client.JID.Local = q.Username
	state.Client.Log = state.Client.Log.WithField("jid", state.Client.JID.Full())
	account := model.NewAccount(state.Client.JID, q.Password)
	err = state.database.AddAccount(account)
	if err != nil {
		state.Client.Out.Encode(&messages.IQClient{
			Type: messages.IQTypeResult,
			To:   state.Client.JID.String(),
			From: state.Client.JID.Domain,
			ID:   msg.ID,
			Body: []byte(fmt.Sprintf(`<query xmlns='%s'>
					<username>%s</username>
					<password>%s</password>
				</query>`, messages.NSIQRegister, q.Username, q.Password)),
			Error: &messages.ErrorClient{
				Code: "409",
				Type: "cancel",
				Any: xml.Name{
					Local: "conflict",
					Space: "urn:ietf:params:xml:ns:xmpp-stanzas",
				},
			},
		})
		state.Client.Log.Warn("database error: ", err)
		return state
	}
	state.Client.Out.Encode(&messages.IQClient{
		Type: messages.IQTypeResult,
		To:   state.Client.JID.String(),
		From: state.Client.JID.Domain,
		ID:   msg.ID,
	})

	state.Client.Log.Infof("registered client %s", state.Client.JID.Bare())
	return state.Next
}
