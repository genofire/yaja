package model

import (
	"errors"
	"sync"
)

type Domain struct {
	FQDN     string
	Accounts map[string]*Account
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
	Local  string
	Domain *Domain
}

func (a *Account) GetJID() *JID {
	return &JID{
		Domain: a.Domain.FQDN,
		Local:  a.Local,
	}
}
