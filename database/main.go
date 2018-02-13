package database

import (
	"errors"
	"sync"

	"dev.sum7.eu/genofire/yaja/messages"
	"dev.sum7.eu/genofire/yaja/model"
	log "github.com/sirupsen/logrus"
)

type State struct {
	Domains map[string]*model.Domain `json:"domains"`
	sync.Mutex
}

func (s *State) AddAccount(a *model.Account) error {
	if a.Local == "" {
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
		_, ok = domain.Accounts[a.Local]
		if ok {
			return errors.New("exists already")
		}
		domain.Accounts[a.Local] = a
		a.Domain = d
		return nil
	}
	return errors.New("no give domain")
}

func (s *State) Authenticate(jid *messages.JID, password string) (bool, error) {
	logger := log.WithField("database", "auth")

	if domain, ok := s.Domains[jid.Domain]; ok {
		if acc, ok := domain.Accounts[jid.Local]; ok {
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

func (s *State) GetAccount(jid *messages.JID) *model.Account {
	logger := log.WithField("database", "get")

	if domain, ok := s.Domains[jid.Domain]; ok {
		if acc, ok := domain.Accounts[jid.Local]; ok {
			return acc
		} else {
			logger.Debug("account not found")
		}
	} else {
		logger.Debug("domain not found")
	}
	return nil
}
