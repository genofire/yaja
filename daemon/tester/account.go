package tester

import "dev.sum7.eu/genofire/yaja/messages"

type Account struct {
	JID      *messages.JID          `json:"jid"`
	Password string                 `json:"password"`
	Admins   map[string]interface{} `json:"admins"`
}
