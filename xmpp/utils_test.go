package xmpp

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartElementToString(t *testing.T) {
	assert := assert.New(t)

	str := XMLStartElementToString(nil)
	assert.Equal("<nil>", str)

	str = XMLStartElementToString(&xml.StartElement{
		Name: xml.Name{
			Local: "iq",
			Space: "jabber:client",
		},
		Attr: []xml.Attr{
			xml.Attr{
				Name: xml.Name{
					Local: "foo",
				},
				Value: "bar",
			},
		},
	})

	assert.Equal(`<iq xmlns="jabber:client" foo="bar">`, str)
}

func remarhal(origin StanzaErrorGroup) StanzaErrorGroup {
	el := StanzaErrorGroup{}
	b, _ := xml.Marshal(origin)
	xml.Unmarshal(b, &el)
	return el
}

func TestChildrenString(t *testing.T) {
	assert := assert.New(t)

	el := remarhal(StanzaErrorGroup{
		Conflict:  &xml.Name{},
		Gone:      "a",
		Forbidden: &xml.Name{},
	})
	str := XMLChildrenString(el)
	assert.Equal("conflict, forbidden", str)

	str = XMLChildrenString(&el)
	assert.Equal("conflict, forbidden", str)
}

func TestCreateCookie(t *testing.T) {
	assert := assert.New(t)

	a := CreateCookieString()
	assert.NotEqual("", a)

	b := CreateCookieString()
	assert.NotEqual("", b)

	assert.NotEqual(a, b)
}
