package xmppbase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test Values for NewJID from RFC7622
// https://tools.ietf.org/html/rfc7622
func TestNewJID(t *testing.T) {
	assert := assert.New(t)

	checkList := map[string]*JID{
		"juliet@example.com": {
			Node:   "juliet",
			Domain: "example.com",
		},
		"juliet@example.com/foo": {
			Node:     "juliet",
			Domain:   "example.com",
			Resource: "foo",
		},
		"juliet@example.com/foo bar": {
			Node:     "juliet",
			Domain:   "example.com",
			Resource: "foo bar",
		},
		"juliet@example.com/foo@bar": {
			Node:     "juliet",
			Domain:   "example.com",
			Resource: "foo@bar",
		},
		"foo\\20bar@example.com": {
			Node:   "foo\\20bar",
			Domain: "example.com",
		},
		"fussball@example.com": {
			Node:   "fussball",
			Domain: "example.com",
		},
		"fu&#xDF;ball@example.com": {
			Node:   "fu&#xDF;ball",
			Domain: "example.com",
		},
		"&#x3C0;@example.com": {
			Node:   "&#x3C0;",
			Domain: "example.com",
		},
		"&#x3A3;@example.com/foo": {
			Node:     "&#x3A3;",
			Domain:   "example.com",
			Resource: "foo",
		},
		"&#x3C3;@example.com/foo": {
			Node:     "&#x3C3;",
			Domain:   "example.com",
			Resource: "foo",
		},
		"&#x3C2;@example.com/foo": {
			Node:     "&#x3C2;",
			Domain:   "example.com",
			Resource: "foo",
		},
		"king@example.com/&#x265A;": {
			Node:     "king",
			Domain:   "example.com",
			Resource: "&#x265A;",
		},
		"example.com": {
			Domain: "example.com",
		},
		"example.com/foobar": {
			Domain:   "example.com",
			Resource: "foobar",
		},
		"a.example.com/b@example.net": {
			Domain:   "a.example.com",
			Resource: "b@example.net",
		},
		"\"juliet\"@example.com":  nil,
		"foo bar@example.com":     nil,
		"juliet@example.com/ foo": nil,
		"@example.com/":           nil,
		// "henry&#x2163;@example.com": nil, -- ignore for easier implementation
		// "&#x265A;@example.com":      nil,
		"juliet@": nil,
		"/foobar": nil,
	}

	for jidString, jidValid := range checkList {
		jid := NewJID(jidString)

		if jidValid != nil {
			assert.NotNil(jid, "this should be a valid JID:"+jidString)
			if jid == nil {
				continue
			}

			assert.Equal(jidValid.Node, jid.Node, "the local part was not right detectet:"+jidString)
			assert.Equal(jidValid.Domain, jid.Domain, "the domain part was not right detectet:"+jidString)
			assert.Equal(jidValid.Resource, jid.Resource, "the resource part was not right detectet:"+jidString)
			assert.Equal(jidValid.Full(), jidString, "the function full of jid did not work")
		} else {
			assert.Nil(jid, "this should not be a valid JID:"+jidString)
		}

	}
}

func TestJIDBare(t *testing.T) {
	assert := assert.New(t)

	checkList := map[string]*JID{
		"aaa@example.com": {
			Node:   "aaa",
			Domain: "example.com",
		},
		"aab@example.com": {
			Node:     "aab",
			Domain:   "example.com",
			Resource: "foo",
		},
		"example.com": {
			Domain:   "example.com",
			Resource: "foo",
		},
	}
	for jidValid, jid := range checkList {
		jidBase := jid.Bare()
		assert.Equal(jidValid, jidBase)

	}
}

func TestMarshal(t *testing.T) {
	assert := assert.New(t)

	jid := &JID{}
	err := jid.UnmarshalText([]byte("juliet@example.com/foo"))

	assert.NoError(err)
	assert.Equal(jid.Node, "juliet")
	assert.Equal(jid.Domain, "example.com")
	assert.Equal(jid.Resource, "foo")

	err = jid.UnmarshalText([]byte("juliet@example.com/ foo"))

	assert.Error(err)

	jid = &JID{
		Node:     "romeo",
		Domain:   "example.com",
		Resource: "bar",
	}
	jidString, err := jid.MarshalText()
	assert.NoError(err)
	assert.Equal("romeo@example.com/bar", string(jidString))
}
