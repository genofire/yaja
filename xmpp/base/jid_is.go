package xmppbase

// IsDomain checks if jid has only domain but no local and resource
func (jid *JID) IsDomain() bool {
	return jid != nil && jid.Local == "" && jid.Domain != "" && jid.Resource == ""
}

// IsBare checks if jid has local and domain but no resource
func (jid *JID) IsBare() bool {
	return jid != nil && jid.Local != "" && jid.Domain != "" && jid.Resource == ""
}

// IsFull checks if jid has all three parts of a JID
func (jid *JID) IsFull() bool {
	return jid != nil && jid.Local != "" && jid.Domain != "" && jid.Resource != ""
}

// IsEqual to check if two jid has same values
func (a *JID) IsEqual(b *JID) bool {
	if a == nil || b == nil {
		return false
	}
	return a.Local == b.Local && a.Domain == b.Domain && a.Resource == b.Resource
}
