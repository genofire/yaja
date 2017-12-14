package model

import (
	"errors"
	"sync"
)

type Domain struct {
	FQDN     string              `json:"-"`
	Accounts map[string]*Account `json:"users"`
	sync.Mutex
}

func (d *Domain) GetJID() *JID {
	return &JID{
		Domain: d.FQDN,
	}
}

func (d *Domain) UpdateAccount(a *Account) error {
	if a.Local == "" {
		return errors.New("No localpart exists in account")
	}
	d.Lock()
	d.Accounts[a.Local] = a
	d.Unlock()
	a.Domain = d
	return nil
}

type Account struct {
	Local    string  `json:"-"`
	Domain   *Domain `json:"-"`
	Password string  `json:"password"`
}

func NewAccount(jid *JID, password string) *Account {
	if jid == nil {
		return nil
	}
	return &Account{
		Local: jid.Local,
		Domain: &Domain{
			FQDN: jid.Domain,
		},
		Password: password,
	}
}

func (a *Account) GetJID() *JID {
	return &JID{
		Domain: a.Domain.FQDN,
		Local:  a.Local,
	}
}

func (a *Account) ValidatePassword(password string) bool {
	return a.Password == password
}
