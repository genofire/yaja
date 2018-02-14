package tester

import "dev.sum7.eu/genofire/yaja/xmpp/base"

type Account struct {
	JID      *xmppbase.JID          `json:"jid"`
	Password string                 `json:"password"`
	Admins   map[string]interface{} `json:"admins"`
}
