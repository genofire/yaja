package xmppbase

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
	Node     string
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
		Node:     jidSplit[1],
		Domain:   jidSplit[2],
		Resource: jidSplit[3],
	}
}

// Bare get the "bare" jid
func (jid *JID) Bare() string {
	if jid == nil {
		return ""
	}
	if jid.Node != "" {
		return jid.Node + "@" + jid.Domain
	}
	return jid.Domain
}

// IsBare checks if jid has node and domain but no resource
func (jid *JID) IsBare() bool {
	return jid != nil && jid.Node != "" && jid.Domain != "" && jid.Resource == ""
}

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

// IsFull checks if jid has all three parts of a JID
func (jid *JID) IsFull() bool {
	return jid != nil && jid.Node != "" && jid.Domain != "" && jid.Resource != ""
}

func (jid *JID) String() string { return jid.Bare() }

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
	jid.Node = newJID.Node
	jid.Domain = newJID.Domain
	jid.Resource = newJID.Resource
	return nil
}
