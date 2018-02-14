package database

import (
	"errors"
	"sync"

	log "github.com/sirupsen/logrus"

	"dev.sum7.eu/genofire/yaja/model"
	"dev.sum7.eu/genofire/yaja/xmpp/base"
)

type State struct {
	Domains map[string]*model.Domain `json:"domains"`
	sync.Mutex
}

func (s *State) AddAccount(a *model.Account) error {
	if a.Node == "" {
		return errors.New("No localpart exists in account")
	}
	if d := a.Domain; d != nil {
		if d.FQDN == "" {
			return errors.New("No fqdn exists in domain")
		}
		s.Lock()
		domain, ok := s.Domains[d.FQDN]
		if !ok {
			if s.Domains == nil {
				s.Domains = make(map[string]*model.Domain)
			}
			s.Domains[d.FQDN] = d
			domain = d
		}
		s.Unlock()

		domain.Lock()
		defer domain.Unlock()
		if domain.Accounts == nil {
			domain.Accounts = make(map[string]*model.Account)
		}
		_, ok = domain.Accounts[a.Node]
		if ok {
			return errors.New("exists already")
		}
		domain.Accounts[a.Node] = a
		a.Domain = d
		return nil
	}
	return errors.New("no give domain")
}

func (s *State) Authenticate(jid *xmppbase.JID, password string) (bool, error) {
	logger := log.WithField("database", "auth")

	if domain, ok := s.Domains[jid.Domain]; ok {
		if acc, ok := domain.Accounts[jid.Node]; ok {
			if acc.ValidatePassword(password) {
				return true, nil
			} else {
				logger.Debug("password not valid")
			}
		} else {
			logger.Debug("account not found")
		}
	} else {
		logger.Debug("domain not found")
	}
	return false, nil
}

func (s *State) GetAccount(jid *xmppbase.JID) *model.Account {
	logger := log.WithField("database", "get")

	if domain, ok := s.Domains[jid.Domain]; ok {
		if acc, ok := domain.Accounts[jid.Node]; ok {
			return acc
		} else {
			logger.Debug("account not found")
		}
	} else {
		logger.Debug("domain not found")
	}
	return nil
}
