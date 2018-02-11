package messages

import (
	"encoding/xml"

	"dev.sum7.eu/genofire/yaja/model"
)

type MessageType string

const (
	MessageTypeChat      MessageType = "chat"
	MessageTypeGroupchat MessageType = "groupchat"
	MessageTypeError     MessageType = "error"
	MessageTypeHeadline  MessageType = "headline"
	MessageTypeNormal    MessageType = "normal"
)

// MessageClient element
type MessageClient struct {
	XMLName xml.Name    `xml:"jabber:client message"`
	From    *model.JID  `xml:"from,attr,omitempty"`
	ID      string      `xml:"id,attr,omitempty"`
	To      *model.JID  `xml:"to,attr,omitempty"`
	Type    MessageType `xml:"type,attr,omitempty"`
	Lang    string      `xml:"lang,attr,omitempty"`
	Subject string      `xml:"subject"`
	Body    string      `xml:"body"`
	Thread  string      `xml:"thread"`
	// Any hasn't matched element
	Other []XMLElement `xml:",any"`

	Delay *Delay `xml:"delay"`
	Error *ErrorClient
}
