package tester

import "dev.sum7.eu/genofire/yaja/model"

type Account struct {
	JID      *model.JID             `json:"jid"`
	Password string                 `json:"password"`
	Admins   map[string]interface{} `json:"admins"`
}
