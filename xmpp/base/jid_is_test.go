package xmppbase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJIDIsDomain(t *testing.T) {
	assert := assert.New(t)
	var jid *JID
	assert.False(jid.IsDomain())

	jid = &JID{}
	assert.False(jid.IsDomain())

	jid = &JID{Local: "a"}
	assert.False(jid.IsDomain())

	jid = &JID{Domain: "a"}
	assert.True(jid.IsDomain())

	jid = &JID{Resource: "a"}
	assert.False(jid.IsDomain())

	jid = &JID{Local: "a", Domain: "b"}
	assert.False(jid.IsDomain())

	jid = &JID{Local: "a", Resource: "b"}
	assert.False(jid.IsDomain())

	jid = &JID{Domain: "a", Resource: "b"}
	assert.False(jid.IsDomain())

	jid = &JID{Local: "a", Domain: "b", Resource: "a"}
	assert.False(jid.IsDomain())
}

func TestJIDIsBare(t *testing.T) {
	assert := assert.New(t)
	var jid *JID
	assert.False(jid.IsBare())

	jid = &JID{}
	assert.False(jid.IsBare())

	jid = &JID{Local: "a"}
	assert.False(jid.IsBare())

	jid = &JID{Domain: "a"}
	assert.False(jid.IsBare())

	jid = &JID{Resource: "a"}
	assert.False(jid.IsBare())

	jid = &JID{Local: "a", Domain: "b"}
	assert.True(jid.IsBare())

	jid = &JID{Local: "a", Resource: "b"}
	assert.False(jid.IsBare())

	jid = &JID{Domain: "a", Resource: "b"}
	assert.False(jid.IsBare())

	jid = &JID{Local: "a", Domain: "b", Resource: "a"}
	assert.False(jid.IsBare())
}

func TestJIDIsFull(t *testing.T) {
	assert := assert.New(t)
	var jid *JID
	assert.False(jid.IsFull())

	jid = &JID{}
	assert.False(jid.IsFull())

	jid = &JID{Local: "a"}
	assert.False(jid.IsFull())

	jid = &JID{Domain: "a"}
	assert.False(jid.IsFull())

	jid = &JID{Resource: "a"}
	assert.False(jid.IsFull())

	jid = &JID{Local: "a", Domain: "b"}
	assert.False(jid.IsFull())

	jid = &JID{Local: "a", Resource: "b"}
	assert.False(jid.IsFull())

	jid = &JID{Domain: "a", Resource: "b"}
	assert.False(jid.IsFull())

	jid = &JID{Local: "a", Domain: "b", Resource: "a"}
	assert.True(jid.IsFull())
}

func TestJIDIsEqual(t *testing.T) {
	assert := assert.New(t)

	// just one null
	var a *JID
	b := &JID{}
	assert.False(a.IsEqual(b))

	a = &JID{}
	// two empty JID
	assert.True(a.IsEqual(b))

	a.Local = "bot"
	b.Local = "bot"
	a.Domain = "example.org"
	b.Domain = "example.org"
	a.Resource = "notebook"
	b.Resource = "notebook"

	assert.True(a.IsEqual(b))

	b.Resource = "mobile"
	assert.False(a.IsEqual(b))

}
