package tester

import (
	"crypto/tls"
	"time"

	"dev.sum7.eu/genofire/yaja/client"
	"dev.sum7.eu/genofire/yaja/model"
)

type Status struct {
	client               *client.Client
	account              *Account
	JID                  *model.JID        `json:"jid"`
	Domain               string            `json:"domain"`
	Login                bool              `json:"is_online"`
	MessageForConnection map[string]string `json:"-"`
	Connections          map[string]bool   `json:"-"`
	TLSVersion           string            `json:"tls_version"`
	IPv4                 bool              `json:"ipv4"`
	IPv6                 bool              `json:"ipv6"`
}

func NewStatus(acc *Account) *Status {
	return &Status{
		account:              acc,
		JID:                  acc.JID,
		Domain:               acc.JID.Domain,
		MessageForConnection: make(map[string]string),
		Connections:          make(map[string]bool),
	}
}

func (s *Status) Update(timeout time.Duration) {
	if s.client == nil || !s.Login {
		s.Login = false
		s.TLSVersion = ""
		return
	}

	c := &client.Client{
		JID:      model.NewJID(s.account.JID.Bare()),
		Protocol: "tcp4",
		Logging:  s.client.Logging,
		Timeout:  timeout / 2,
	}

	if err := c.Connect(s.account.Password); err == nil {
		s.IPv4 = true
		c.Close()
	} else {
		s.IPv4 = false
	}

	c.JID = model.NewJID(s.account.JID.Bare())
	c.Protocol = "tcp6"

	if err := c.Connect(s.account.Password); err == nil {
		s.IPv6 = true
		c.Close()
	} else {
		s.IPv6 = false
	}
	if !s.IPv4 && !s.IPv6 {
		s.client.Close()
		s.Login = false
		s.TLSVersion = ""
	}

	if tlsstate := s.client.TLSConnectionState(); tlsstate != nil {
		switch tlsstate.Version {
		case tls.VersionSSL30:
			s.TLSVersion = "SSL 3.0"
		case tls.VersionTLS10:
			s.TLSVersion = "TLS 1.0"
		case tls.VersionTLS11:
			s.TLSVersion = "TLS 1.1"
		case tls.VersionTLS12:
			s.TLSVersion = "TLS 1.2"
		default:
			s.TLSVersion = "unknown " + string(tlsstate.Version)
		}
	} else {
		s.TLSVersion = ""
	}
}
