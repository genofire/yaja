package xmppbase

import (
	"errors"
	"regexp"
)

var jidRegex *regexp.Regexp

func init() {
	jidRegex = regexp.MustCompile(`^(?:([^@/<>'\" ]+)@)?([^@/<>'\"]+)(?:/([^<>'\" ][^<>'\"]*))?$`)
}

// JID implements RFC 6122: XMPP - Address Format
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

// Clone JID struct address/pointer
func (jid *JID) Clone() *JID {
	if jid != nil {
		return &JID{
			Local:    jid.Local,
			Domain:   jid.Domain,
			Resource: jid.Resource,
		}
	}
	return nil
}

// Full get the "full" jid as string
func (jid *JID) Full() *JID {
	return jid.Clone()
}

// Bare get the "bare" jid
func (jid *JID) Bare() *JID {
	if jid != nil {
		return &JID{
			Local:  jid.Local,
			Domain: jid.Domain,
		}
	}
	return nil
}

func (jid *JID) String() string {
	if jid == nil {
		return ""
	}
	str := jid.Domain
	if jid.Local != "" {
		str = jid.Local + "@" + str
	}
	if jid.Resource != "" {
		str = str + "/" + jid.Resource
	}
	return str
}

//MarshalText to bytearray
func (jid JID) MarshalText() ([]byte, error) {
	return []byte(jid.String()), nil
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
