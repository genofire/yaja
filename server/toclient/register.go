package toclient

import (
	"encoding/xml"

	"dev.sum7.eu/genofire/yaja/database"
	"dev.sum7.eu/genofire/yaja/model"
	"dev.sum7.eu/genofire/yaja/server/state"
	"dev.sum7.eu/genofire/yaja/server/utils"
	"dev.sum7.eu/genofire/yaja/xmpp"
	"dev.sum7.eu/genofire/yaja/xmpp/base"
	"dev.sum7.eu/genofire/yaja/xmpp/iq"
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

	var msg xmpp.IQClient
	if err := state.Client.In.DecodeElement(&msg, state.element); err != nil {
		state.Client.Log.Warn("is no iq: ", err)
		return state
	}
	if msg.Type != xmpp.IQTypeGet {
		state.Client.Log.Warn("is no get iq")
		return state
	}
	if msg.Error != nil {
		state.Client.Log.Warn("iq with error: ", msg.Error.Code)
		return state
	}

	if msg.PrivateRegister == nil {
		state.Client.Log.Warn("is no iq register")
		return nil
	}
	state.Client.Out.Encode(&xmpp.IQClient{
		Type: xmpp.IQTypeResult,
		To:   state.Client.JID,
		From: xmppbase.NewJID(state.Client.JID.Domain),
		ID:   msg.ID,
		PrivateRegister: &xmppiq.IQPrivateRegister{
			Instructions: "Choose a username and password for use with this service.",
			Username:     "",
			Password:     "",
		},
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
	var msg xmpp.IQClient
	if err = state.Client.In.DecodeElement(&msg, element); err != nil {
		state.Client.Log.Warn("is no iq: ", err)
		return state
	}
	if msg.Type != xmpp.IQTypeGet {
		state.Client.Log.Warn("is no get iq")
		return state
	}
	if msg.Error != nil {
		state.Client.Log.Warn("iq with error: ", msg.Error.Code)
		return state
	}
	if msg.PrivateRegister == nil {
		state.Client.Log.Warn("is no iq register: ", err)
		return nil
	}

	state.Client.JID.Node = msg.PrivateRegister.Username
	state.Client.Log = state.Client.Log.WithField("jid", state.Client.JID.Full())
	account := model.NewAccount(state.Client.JID, msg.PrivateRegister.Password)
	err = state.database.AddAccount(account)
	if err != nil {
		state.Client.Out.Encode(&xmpp.IQClient{
			Type:            xmpp.IQTypeResult,
			To:              state.Client.JID,
			From:            xmppbase.NewJID(state.Client.JID.Domain),
			ID:              msg.ID,
			PrivateRegister: msg.PrivateRegister,
			Error: &xmpp.ErrorClient{
				Code: "409",
				Type: "cancel",
				StanzaErrorGroup: xmpp.StanzaErrorGroup{
					Conflict: &xml.Name{},
				},
			},
		})
		state.Client.Log.Warn("database error: ", err)
		return state
	}
	state.Client.Out.Encode(&xmpp.IQClient{
		Type: xmpp.IQTypeResult,
		To:   state.Client.JID,
		From: xmppbase.NewJID(state.Client.JID.Domain),
		ID:   msg.ID,
	})

	state.Client.Log.Infof("registered client %s", state.Client.JID.Bare())
	return state.Next
}
