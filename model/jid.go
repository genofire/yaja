package model

import (
	"errors"
	"regexp"
)

var jidRegex *regexp.Regexp

func init() {
	jidRegex = regexp.MustCompile(`^(?:([^@/<>'\" ]+)@)?([^@/<>'\"]+)(?:/([^<>'\" ][^<>'\"]*))?$`)
}

// JID struct
type JID struct {
	Local    string
	Domain   string
	Resource string
}

// NewJID get JID from string
func NewJID(jidString string) *JID {
	jidSplitTmp := jidRegex.FindAllStringSubmatch(jidString, -1)
	if len(jidSplitTmp) != 1 {
		return nil
	}

	jidSplit := jidSplitTmp[0]

	return &JID{
		Local:    jidSplit[1],
		Domain:   jidSplit[2],
		Resource: jidSplit[3],
	}
}

// Bare get the "bare" jid
func (jid *JID) Bare() string {
	if jid == nil {
		return ""
	}
	if jid.Local != "" {
		return jid.Local + "@" + jid.Domain
	}
	return jid.Domain
}

func (jid *JID) String() string { return jid.Bare() }

// Full get the "full" jid as string
func (jid *JID) Full() string {
	if jid == nil {
		return ""
	}
	if jid.Resource != "" {
		return jid.Bare() + "/" + jid.Resource
	}
	return jid.Bare()
}

//MarshalText to bytearray
func (jid JID) MarshalText() ([]byte, error) {
	return []byte(jid.Full()), nil
}

// UnmarshalText from bytearray
func (jid *JID) UnmarshalText(data []byte) (err error) {
	newJID := NewJID(string(data))
	if newJID == nil {
		return errors.New("not a valid jid")
	}
	jid.Local = newJID.Local
	jid.Domain = newJID.Domain
	jid.Resource = newJID.Resource
	return nil
}
