package tester

import (
	"crypto/tls"
	"fmt"
	"time"

	"dev.sum7.eu/genofire/yaja/client"
	"dev.sum7.eu/genofire/yaja/xmpp"
	"dev.sum7.eu/genofire/yaja/xmpp/base"
	"dev.sum7.eu/genofire/yaja/xmpp/iq"
)

type Status struct {
	backupClient         *client.Client
	client               *client.Client
	account              *Account
	JID                  *xmppbase.JID     `json:"jid"`
	Domain               string            `json:"domain"`
	Login                bool              `json:"is_online"`
	MessageForConnection map[string]string `json:"-"`
	Connections          map[string]bool   `json:"-"`
	TLSVersion           string            `json:"tls_version"`
	IPv4                 bool              `json:"ipv4"`
	IPv6                 bool              `json:"ipv6"`
	Software             string            `json:"software,omitempty"`
	OS                   string            `json:"os,omitempty"`
}

func NewStatus(backupClient *client.Client, acc *Account) *Status {
	return &Status{
		backupClient:         backupClient,
		account:              acc,
		JID:                  acc.JID,
		Domain:               acc.JID.Domain,
		MessageForConnection: make(map[string]string),
		Connections:          make(map[string]bool),
	}
}
func (s *Status) Disconnect(reason string) {
	if s.Login {
		msg := &xmpp.MessageClient{
			Type: xmpp.MessageTypeChat,
			Body: fmt.Sprintf("you receive a notify that '%s' disconnect: %s", s.JID.Full(), reason),
		}
		for jid := range s.account.Admins {
			msg.To = xmppbase.NewJID(jid)
			if err := s.backupClient.Send(msg); err != nil {
				s.client.Send(msg)
			}
		}
	}
	s.client.Logging.Warnf("status-disconnect: %s", reason)
	s.client.Close()
	s.Login = false
	s.TLSVersion = ""
}

func (s *Status) update(timeout time.Duration) {
	if s.client == nil || !s.Login {
		return
	}

	c := &client.Client{
		JID:      s.account.JID.Bare(),
		Protocol: "tcp4",
		Logging:  s.client.Logging.WithField("status", "ipv4"),
		Timeout:  timeout / 2,
	}

	if err := c.Connect(s.account.Password); err == nil {
		s.IPv4 = true
		c.Close()
	} else {
		s.IPv4 = false
	}

	c.Logging = s.client.Logging.WithField("status", "ipv4")
	c.JID = s.account.JID.Bare()
	c.Protocol = "tcp6"

	if err := c.Connect(s.account.Password); err == nil {
		s.IPv6 = true
		c.Close()
	} else {
		s.IPv6 = false
	}
	if !s.IPv4 && !s.IPv6 {
		s.Disconnect("check of ipv4 and ipv6 failed -> client should not be connected anymore")
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
	iq, err := s.client.SendRecv(&xmpp.IQClient{
		To:      &xmppbase.JID{Domain: s.JID.Domain},
		Type:    xmpp.IQTypeGet,
		Version: &xmppiq.Version{},
	})
	if err != nil {
		s.client.Logging.Errorf("status-update: %s", err.Error())
	} else if iq != nil {
		if iq.Error != nil && iq.Error.ServiceUnavailable != nil {
			s.Software = "unknown"
			s.OS = "unknown"
		} else if iq.Version != nil {
			s.Software = iq.Version.Name
			if iq.Version.Version != "" {
				s.Software += "-" + iq.Version.Version
			}
			if s.Software == "" {
				s.Software = "unknown"
			}
			s.OS = iq.Version.OS
			if s.OS == "" {
				s.OS = "unknown"
			}
		} else {
			s.Software = ""
			s.OS = ""
		}
	}
}
