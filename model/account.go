package model

import (
	"errors"
	"sync"

	"dev.sum7.eu/genofire/yaja/xmpp/base"
)

type Domain struct {
	FQDN     string              `json:"-"`
	Accounts map[string]*Account `json:"users"`
	sync.Mutex
}

func (d *Domain) GetJID() *xmppbase.JID {
	return &xmppbase.JID{
		Domain: d.FQDN,
	}
}

func (d *Domain) UpdateAccount(a *Account) error {
	if a.Node == "" {
		return errors.New("No localpart exists in account")
	}
	d.Lock()
	d.Accounts[a.Node] = a
	d.Unlock()
	a.Domain = d
	return nil
}

type Account struct {
	Node      string               `json:"-"`
	Domain    *Domain              `json:"-"`
	Password  string               `json:"password"`
	Roster    map[string]*Buddy    `json:"roster"`
	Bookmarks map[string]*Bookmark `json:"bookmarks"`
}

func NewAccount(jid *xmppbase.JID, password string) *Account {
	if jid == nil {
		return nil
	}
	return &Account{
		Node: jid.Node,
		Domain: &Domain{
			FQDN: jid.Domain,
		},
		Password: password,
	}
}

func (a *Account) GetJID() *xmppbase.JID {
	return &xmppbase.JID{
		Domain: a.Domain.FQDN,
		Node:   a.Node,
	}
}

func (a *Account) ValidatePassword(password string) bool {
	return a.Password == password
}
